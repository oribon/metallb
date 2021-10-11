package e2e

import (
	"context"

	operatorv1alpha1 "github.com/metallb/metallb-operator/api/v1alpha1"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Updater interface {
	Update(cf configFile) error
}

type Cleaner interface {
	Clean() error
}

type UpdateCleaner interface {
	Updater
	Cleaner
}

type cmUpdateCleaner struct {
	clientset.Interface
	name      string
	namespace string
}

func newConfigMap(c clientset.Interface, ns string) *cmUpdateCleaner {
	return &cmUpdateCleaner{
		Interface: c,
		name:      "config",
		namespace: ns,
	}
}

func (c cmUpdateCleaner) Update(cf configFile) error {
	resData, err := yaml.Marshal(cf)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal MetalLB ConfigMap data")
	}

	_, err = c.CoreV1().ConfigMaps(c.namespace).Update(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.name,
			Namespace: c.namespace,
		},
		Data: map[string]string{"config": string(resData)},
	}, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrapf(err, "Failed to update MetalLB ConfigMap")
	}

	return nil
}

func (c cmUpdateCleaner) Clean() error {
	return c.Update(configFile{})
}

type operatorUpdateCleaner struct {
	client.Client
	namespace string
}

func newOperator(r *rest.Config, ns string) (*operatorUpdateCleaner, error) {
	myScheme := runtime.NewScheme()

	if err := operatorv1alpha1.AddToScheme(myScheme); err != nil {
		return nil, err
	}

	cl, err := client.New(r, client.Options{
		Scheme: myScheme,
	})

	if err != nil {
		return nil, err
	}

	return &operatorUpdateCleaner{
		Client:    cl,
		namespace: ns,
	}, nil
}

func (o operatorUpdateCleaner) Update(cf configFile) error {
	for _, p := range cf.Pools {
		err := o.createPool(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o operatorUpdateCleaner) Clean() error {
	err := o.DeleteAllOf(context.Background(), &operatorv1alpha1.AddressPool{}, client.InNamespace(o.namespace))
	if err != nil {
		return err
	}

	// BGPPeer unexported
	//err = cs.DeleteAllOf(context.Background(), &operatorv1alpha1.BGPPeer, client.InNamespace(o.namespace))
	//if err != nil {
	//return err
	//}

	return nil
}
func (o operatorUpdateCleaner) createPool(p addressPool) error {
	pool := o.poolToOperator(p)
	return o.Client.Create(context.TODO(), pool)
}

func (o operatorUpdateCleaner) poolToOperator(p addressPool) *operatorv1alpha1.AddressPool {
	res := &operatorv1alpha1.AddressPool{}
	res.Name = p.Name
	res.Namespace = o.namespace
	res.Spec.Protocol = string(p.Protocol)
	res.Spec.Addresses = p.Addresses
	res.Spec.AutoAssign = p.AutoAssign
	// AvoidBuggyIPs missing from operator
	// BGPAdvertisements missing from v0.10.2
	return res
}

// BGPPeer missing from v0.10.2
//func peerToOperator(p peer) {}
