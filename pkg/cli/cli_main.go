package cli

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"os"

	// Kubernetes:
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// App contains flags, arguments and commands for an application:
	App = kingpin.New("kubencrypt", "Letsencrypt on Kubernetes.")

	// Flags:
	flgKubeconfig = App.Flag("kubeconfig",
		"Absolute path to the kubeconfig file.").
		Default(kubeconfigPath()).ExistingFileOrDir()

	// FlgNamespace contains the namespace name:
	FlgNamespace = App.Flag("namespace",
		"If present, the namespace scope for this request.").
		Default(apiv1.NamespaceAll).HintAction(listNamespaces).String()

	// FlgIngress contains the ingress name:
	FlgIngress = App.Flag("ingress",
		"Name of the ingress object to be altered.").
		Required().HintAction(listIngresses).String()

	// FlgSecret contains the secret name:
	FlgSecret = App.Flag("secret",
		"Name of the secret object to be altered.").
		Required().HintAction(listSecrets).String()

	// FlgServiceName contains the service name:
	FlgServiceName = App.Flag("service-name",
		"Name of the k8s letsencrypt service.").
		Required().String()

	// FlgServicePort contains the service port:
	FlgServicePort = App.Flag("service-port",
		"Port of the k8s letsencrypt service.").
		Required().Int()

	// FlgDomain contains the domain name:
	FlgDomain = App.Flag("domain",
		"Prove control of this domain.").
		Required().String()
)

//-----------------------------------------------------------------------------
// func init() is called after all the variable declarations in the package
// have evaluated their initializers, and those are evaluated only after all
// the imported packages have been initialized:
//-----------------------------------------------------------------------------

func init() {

	// Customize kingpin:
	App.Version("v0.1.0").Author("Marc Villacorta Morera")
	App.HelpFlag.Short('h')
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

// K8sConnect establishes a k8s connection.
func K8sConnect() (*kubernetes.Clientset, error) {

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
	clientset, err := K8sConnect()
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
	clientset, err := K8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get the list of ingresses objects:
	l, err := clientset.ExtensionsV1beta1().Ingresses(*FlgNamespace).List(metav1.ListOptions{})
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
	clientset, err := K8sConnect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Get the list of secrets objects:
	l, err := clientset.CoreV1().Secrets(*FlgNamespace).List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Extract the name of each secret:
	for _, v := range l.Items {
		list = append(list, v.Name)
	}

	return
}
