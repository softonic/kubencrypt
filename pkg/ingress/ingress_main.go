package ingress

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Kubencrypt:
	"github.com/softonic/kubencrypt/pkg/cli"

	// Kubernetes:
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	types "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	//_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// Community:
	log "github.com/Sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Typedefs:
//-----------------------------------------------------------------------------

// Data holds ingress data.
type Data struct {

	// Flags:
	Namespace   string
	IngressName string
	ServiceName string
	ServicePort int

	// Handlers:
	clientset *kubernetes.Clientset
	client    types.IngressInterface
	paths     *[]extensionsv1beta1.HTTPIngressPath

	// Data:
	ingress *extensionsv1beta1.Ingress
	backup  *extensionsv1beta1.Ingress
}

//-----------------------------------------------------------------------------
// func: connect
//-----------------------------------------------------------------------------

func (d *Data) connect() (err error) {

	if d.clientset == nil {

		// Log:
		log.Info("Connecting to kubernetes...")

		// Connect to the cluster:
		d.clientset, err = cli.K8sConnect()
		if err != nil {
			return err
		}

		// Ingress client handler:
		d.client = d.clientset.ExtensionsV1beta1().Ingresses(d.Namespace)
	}

	return nil
}

//-----------------------------------------------------------------------------
// func: Backup
//-----------------------------------------------------------------------------

// Backup retrieves the current ingress object and makes a copy.
func (d *Data) Backup() {

	// Connect:
	err := d.connect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Log:
	log.Info("Backing up the current ingress...")

	// Get my ingress:
	d.ingress, err = d.client.Get(d.IngressName, metav1.GetOptions{})
	if err != nil {
		log.Panic(err.Error())
	}

	// Backup the ingress:
	d.backup = d.ingress.DeepCopy()
}

//-----------------------------------------------------------------------------
// func: Update
//-----------------------------------------------------------------------------

// Update adds a path into the first ingress rule.
func (d *Data) Update() {

	// Connect:
	err := d.connect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Log:
	log.Info("Adding the letsencrypt path...")

	// Forge the new data:
	d.paths = &d.ingress.Spec.Rules[0].HTTP.Paths
	*d.paths = append(*d.paths, extensionsv1beta1.HTTPIngressPath{
		Path: "/.well-known/*",
		Backend: extensionsv1beta1.IngressBackend{
			ServiceName: d.ServiceName,
			ServicePort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(d.ServicePort),
			},
		},
	})

	// Insert the new data:
	if d.ingress, err = d.client.Update(d.ingress); err != nil {
		log.Panic(err)
	}
}

//-----------------------------------------------------------------------------
// func: Restore
//-----------------------------------------------------------------------------

// Restore restores the original ingress object.
func (d *Data) Restore() {

	// Connect:
	err := d.connect()
	if err != nil {
		log.Panic(err.Error())
	}

	// Log:
	log.Info("Restoring the original ingress...")

	// Modify loop:
	for {

		d.ingress.Spec.Rules[0].HTTP.Paths = d.backup.Spec.Rules[0].HTTP.Paths

		// Retry Update() until you no longer get a conflict error:
		if _, err = d.client.Update(d.ingress); errors.IsConflict(err) {
			log.Warn("Encountered conflict, retrying")
			d.ingress, err = d.client.Get(d.IngressName, metav1.GetOptions{})
			if err != nil {
				log.Panic(err.Error())
			}
		} else if err != nil {
			log.Panic(err.Error())
		} else {
			break
		}
	}
}
