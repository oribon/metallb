// SPDX-License-Identifier:Apache-2.0

package config

import (
	"context"
	"time"

	operatorv1beta1 "github.com/metallb/metallb-operator/api/v1beta1"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Updater interface {
	Update(f File) error
	Clean() error
}

type configmapUpdater struct {
	clientset.Interface
	name      string
	namespace string
}

func UpdaterForConfigMap(c clientset.Interface, ns string) *configmapUpdater {
	return &configmapUpdater{
		Interface: c,
		name:      "config",
		namespace: ns,
	}
}

func (c configmapUpdater) Update(f File) error {
	resData, err := yaml.Marshal(f)
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

func (c configmapUpdater) Clean() error {
	return c.Update(File{})
}

type operatorUpdater struct {
	client.Client
	namespace string
}

func UpdaterForOperator(r *rest.Config, ns string) (*operatorUpdater, error) {
	myScheme := runtime.NewScheme()

	if err := operatorv1beta1.AddToScheme(myScheme); err != nil {
		return nil, err
	}

	cl, err := client.New(r, client.Options{
		Scheme: myScheme,
	})

	if err != nil {
		return nil, err
	}

	return &operatorUpdater{
		Client:    cl,
		namespace: ns,
	}, nil
}

func (o operatorUpdater) Update(f File) error {
	for _, a := range f.Pools {
		err := o.createPool(a)
		if err != nil {
			return err
		}
	}

	for _, bp := range f.BFDProfiles {
		err := o.createBFDProfile(bp)
		if err != nil {
			return err
		}
	}

	/*
		Since a peer doesn't have a name we need to clean them before creating.
		Without this the webhook fails the request because we are trying to create
		a BGPPeer with the same configuration under a different (generated) name.
		TODO: is there a better way to handle this?
	*/
	err := o.DeleteAllOf(context.Background(), &operatorv1beta1.BGPPeer{}, client.InNamespace(o.namespace))
	if err != nil {
		return err
	}

	for _, p := range f.Peers {
		err := o.createPeer(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o operatorUpdater) Clean() error {
	err := o.DeleteAllOf(context.Background(), &operatorv1beta1.AddressPool{}, client.InNamespace(o.namespace))
	if err != nil {
		return err
	}

	err = o.DeleteAllOf(context.Background(), &operatorv1beta1.BFDProfile{}, client.InNamespace(o.namespace))
	if err != nil {
		return err
	}

	err = o.DeleteAllOf(context.Background(), &operatorv1beta1.BGPPeer{}, client.InNamespace(o.namespace))
	if err != nil {
		return err
	}

	return nil
}

func (o operatorUpdater) createPool(a AddressPool) error {
	pool := o.poolToOperator(a)
	existing := &operatorv1beta1.AddressPool{}
	err := o.Client.Get(context.TODO(), types.NamespacedName{Name: pool.Name, Namespace: pool.Namespace}, existing)

	if err != nil {
		if apierr.IsNotFound(err) {
			return o.Client.Create(context.TODO(), pool)
		}
		return err
	}

	pool.ResourceVersion = existing.ResourceVersion
	return o.Client.Update(context.TODO(), pool)
}

func (o operatorUpdater) poolToOperator(a AddressPool) *operatorv1beta1.AddressPool {
	addrs := make([]string, len(a.Addresses))
	copy(addrs, a.Addresses)

	bgppadvs := make([]operatorv1beta1.BgpAdvertisement, len(a.BGPAdvertisements))
	for _, b := range a.BGPAdvertisements {
		adv := operatorv1beta1.BgpAdvertisement{}
		if b.AggregationLength != nil {
			al := int32(*b.AggregationLength)
			adv.AggregationLength = &al
		}
		if b.LocalPref != nil {
			adv.LocalPref = *b.LocalPref
		}
		copy(adv.Communities, b.Communities)
		bgppadvs = append(bgppadvs, adv)
	}

	return &operatorv1beta1.AddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: o.namespace,
		},
		Spec: operatorv1beta1.AddressPoolSpec{
			Protocol:  string(a.Protocol),
			Addresses: addrs,
			// AvoidBuggyIPs missing from operator
			AutoAssign:        a.AutoAssign,
			BGPAdvertisements: bgppadvs,
		},
	}
}

func (o operatorUpdater) createPeer(p Peer) error {
	peer, err := o.peerToOperator(p)
	if err != nil {
		return err
	}

	// peer's name is autogenerated, will not hit a name conflict here.
	return o.Client.Create(context.TODO(), peer)
}

func (o operatorUpdater) peerToOperator(p Peer) (*operatorv1beta1.BGPPeer, error) {
	var holdtime time.Duration
	var err error
	if p.HoldTime != "" {
		holdtime, err = time.ParseDuration(p.HoldTime)
		if err != nil {
			return nil, err
		}
	}

	nodeselectors := make([]operatorv1beta1.NodeSelector, len(p.NodeSelectors))
	for _, ns := range p.NodeSelectors {
		n := operatorv1beta1.NodeSelector{
			MatchLabels: ns.MatchLabels,
		}
		for _, e := range ns.MatchExpressions {
			vals := make([]string, len(e.Values))
			copy(vals, e.Values)
			expr := operatorv1beta1.MatchExpression{
				Key:      e.Key,
				Operator: e.Operator,
				Values:   vals,
			}
			n.MatchExpressions = append(n.MatchExpressions, expr)
		}
		nodeselectors = append(nodeselectors, n)
	}

	return &operatorv1beta1.BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "testpeer-",
			Namespace:    o.namespace,
		},
		Spec: operatorv1beta1.BGPPeerSpec{
			MyASN:         p.MyASN,
			ASN:           p.ASN,
			Address:       p.Addr,
			SrcAddress:    p.SrcAddr,
			Port:          p.Port,
			HoldTime:      holdtime,
			RouterID:      p.RouterID,
			NodeSelectors: nodeselectors,
			Password:      p.Password,
			BFDProfile:    p.BFDProfile,
		},
	}, nil

}

func (o operatorUpdater) createBFDProfile(bp BfdProfile) error {
	bfdprofile := o.bfdProfileToOperator(bp)
	existing := &operatorv1beta1.BFDProfile{}
	err := o.Client.Get(context.TODO(), types.NamespacedName{Name: bfdprofile.Name, Namespace: bfdprofile.Namespace}, existing)

	if err != nil {
		if apierr.IsNotFound(err) {
			return o.Client.Create(context.TODO(), bfdprofile)
		}
		return err
	}

	bfdprofile.ResourceVersion = existing.ResourceVersion
	return o.Client.Update(context.TODO(), bfdprofile)
}

func (o operatorUpdater) bfdProfileToOperator(bp BfdProfile) *operatorv1beta1.BFDProfile {
	return &operatorv1beta1.BFDProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bp.Name,
			Namespace: o.namespace,
		},
		Spec: operatorv1beta1.BFDProfileSpec{
			ReceiveInterval:  bp.ReceiveInterval,
			TransmitInterval: bp.TransmitInterval,
			DetectMultiplier: bp.DetectMultiplier,
			EchoInterval:     bp.EchoInterval,
			EchoMode:         bp.EchoMode,
			PassiveMode:      bp.PassiveMode,
			MinimumTTL:       bp.MinimumTTL,
		},
	}
}
