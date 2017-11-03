package main

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"os"
	"time"

	// Kubernetes:
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Community:
	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//-----------------------------------------------------------------------------
// Command, flags and arguments:
//-----------------------------------------------------------------------------

var (

	// Root level command:
	app = kingpin.New("kubencrypt", "Letsencrypt on Kubernetes.")

	// Flags:
	flgKubeconfig = app.Flag("kubeconfig",
		"Absolute path to the kubeconfig file.").
		Default(kubeconfigPath()).ExistingFileOrDir()

	flgNamespace = app.Flag("namespace",
		"If present, the namespace scope for this request.").
		Default(apiv1.NamespaceAll).HintAction(listNamespaces).String()

	flgIngress = app.Flag("ingress",
		"Name of the ingress object to be altered.").
		Required().HintAction(listIngresses).String()

	flgSecret = app.Flag("secret",
		"Name of the secret object to be altered.").
		Required().HintAction(listSecrets).String()

	flgServiceName = app.Flag("service-name",
		"Name of the k8s letsencrypt service.").
		Required().String()

	flgServicePort = app.Flag("service-port",
		"Port of the k8s letsencrypt service.").
		Required().Int()
)

//-----------------------------------------------------------------------------
// func init() is called after all the variable declarations in the package
// have evaluated their initializers, and those are evaluated only after all
// the imported packages have been initialized:
//-----------------------------------------------------------------------------

func init() {

	// Customize kingpin:
	app.Version("v0.1.0").Author("Marc Villacorta Morera")
	app.HelpFlag.Short('h')

	// Customize the default logger:
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

//-----------------------------------------------------------------------------
// Entry point:
//-----------------------------------------------------------------------------

func main() {

	// Parse the command-line flags:
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Connect to the cluster:
	clientset, err := k8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get my ingress:
	ingressClient := clientset.ExtensionsV1beta1().Ingresses(*flgNamespace)
	myIngress, err := ingressClient.Get(*flgIngress, metav1.GetOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Backup the first rule paths:
	paths := &myIngress.Spec.Rules[0].HTTP.Paths
	pathsBackup := *paths

	// Add a path to the first rule:
	*paths = append(*paths, extensionsv1beta1.HTTPIngressPath{
		Path: "/.well-known/*",
		Backend: extensionsv1beta1.IngressBackend{
			ServiceName: *flgServiceName,
			ServicePort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(*flgServicePort),
			},
		},
	})

	if _, err = ingressClient.Update(myIngress); err != nil {
		log.Panic(err)
	}

	// Restore the backup:
	time.Sleep(time.Second * 10)
	myIngress, err = ingressClient.Get(*flgIngress, metav1.GetOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	paths = &myIngress.Spec.Rules[0].HTTP.Paths
	*paths = pathsBackup

	if _, err := ingressClient.Update(myIngress); err != nil {
		log.Panic(err.Error())
	}

	// Get my secret:
	mySecret, err := clientset.CoreV1().Secrets(*flgNamespace).Get(*flgSecret, metav1.GetOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Log ingresses and secrets count:
	log.WithField("name", myIngress.GetObjectMeta().GetName()).Info("My ingress")
	log.WithField("name", mySecret.GetObjectMeta().GetName()).Info("My secret")
}

//-----------------------------------------------------------------------------
// kubeconfigPath:
//-----------------------------------------------------------------------------

func kubeconfigPath() (path string) {

	// Return ~/.kube/config if exists...
	if _, err := os.Stat(os.Getenv("HOME") + "/.kube/config"); err == nil {
		return os.Getenv("HOME") + "/.kube/config"
	}

	// ...otherwise return '.':
	return "."
}

//-----------------------------------------------------------------------------
// k8sConnect:
//-----------------------------------------------------------------------------

func k8sConnect() (*kubernetes.Clientset, error) {

	// Build the config:
	config, err := buildConfig(*flgKubeconfig)
	if err != nil {
		return nil, err
	}

	// Create the clientset:
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Return the clientset:
	return clientset, nil
}

//-----------------------------------------------------------------------------
// buildConfig:
//-----------------------------------------------------------------------------

func buildConfig(kubeconfig string) (*rest.Config, error) {

	// Use kubeconfig if given...
	if kubeconfig != "" && kubeconfig != "." {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	// ...otherwise assume in-cluster:
	log.Info("Running in-cluster using environment variables")
	return rest.InClusterConfig()
}

//-----------------------------------------------------------------------------
// listNamespaces:
//-----------------------------------------------------------------------------

func listNamespaces() (list []string) {

	// Connect to the cluster:
	clientset, err := k8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get the list of namespace objects:
	l, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Extract the name of each namespace:
	for _, v := range l.Items {
		list = append(list, v.Name)
	}

	return
}

//-----------------------------------------------------------------------------
// listIngresses:
//-----------------------------------------------------------------------------

func listIngresses() (list []string) {

	// Connect to the cluster:
	clientset, err := k8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get the list of ingresses objects:
	l, err := clientset.ExtensionsV1beta1().Ingresses(*flgNamespace).List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Extract the name of each ingress:
	for _, v := range l.Items {
		list = append(list, v.Name)
	}

	return
}

//-----------------------------------------------------------------------------
// listSecrets:
//-----------------------------------------------------------------------------

func listSecrets() (list []string) {

	// Connect to the cluster:
	clientset, err := k8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get the list of secrets objects:
	l, err := clientset.CoreV1().Secrets(*flgNamespace).List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Extract the name of each secret:
	for _, v := range l.Items {
		list = append(list, v.Name)
	}

	return
}
