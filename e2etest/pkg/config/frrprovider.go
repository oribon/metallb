// SPDX-License-Identifier:Apache-2.0

package config

import (
	"fmt"

	"go.universe.tf/e2etest/pkg/executor"
	"go.universe.tf/e2etest/pkg/metallb"
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type FRRProvider interface {
	FRRForSpeaker(namespace, name string) (executor.Executor, error)
	BGPMetricsPodForSpeaker(namespace, name string) (pod *corev1.Pod, metricsPrefix string, err error)
}

type frrModeProvider struct {
	speakers map[string]*corev1.Pod
}

func NewFRRModeProvider(r *rest.Config) (FRRProvider, error) {
	cl, err := clientset.NewForConfig(r)
	if err != nil {
		return nil, err
	}

	speakerPods, err := metallb.SpeakerPods(cl)
	if err != nil {
		return nil, err
	}

	speakers := map[string]*corev1.Pod{}
	for _, p := range speakerPods {
		speakers[p.Name] = p
	}

	return frrModeProvider{speakers: speakers}, nil
}

func (f frrModeProvider) FRRForSpeaker(ns, name string) (executor.Executor, error) {
	_, ok := f.speakers[name]
	if !ok {
		return nil, fmt.Errorf("speakers %s/%s not found in known speakers %v", ns, name, f.speakers)
	}

	return executor.ForPod(ns, name, "frr"), nil
}

func (f frrModeProvider) BGPMetricsPodForSpeaker(ns, name string) (*corev1.Pod, string, error) {
	p, ok := f.speakers[name]
	if !ok {
		return nil, "", fmt.Errorf("speakers %s/%s not found in in known speakers %v", ns, name, f.speakers)
	}

	return p, "metallb", nil
}

type frrk8sModeProvider struct {
	frrk8sPodForSpeakerPod map[string]*corev1.Pod
}

func NewFRRK8SModeProvider(r *rest.Config) (FRRProvider, error) { // TODO: make it work regardless of namespaces
	cl, err := clientset.NewForConfig(r)
	if err != nil {
		return nil, err
	}

	speakerPods, err := metallb.SpeakerPods(cl)
	if err != nil {
		return nil, err
	}

	frrk8sPods, err := metallb.FRRK8SPods(cl)
	if err != nil {
		return nil, err
	}

	frrK8SForSpeaker := map[string]*corev1.Pod{}
	for _, s := range speakerPods {
		found := false
		for _, f := range frrk8sPods {
			if s.Spec.NodeName == f.Spec.NodeName {
				frrK8SForSpeaker[s.Name] = f
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("speaker %s/%s on node %s does not have a matching frr-k8s", s.Namespace, s.Name, s.Spec.NodeName)
		}
	}

	res := frrk8sModeProvider{
		frrk8sPodForSpeakerPod: frrK8SForSpeaker,
	}

	return res, nil
}

func (f frrk8sModeProvider) FRRForSpeaker(ns, name string) (executor.Executor, error) {
	frrk8s, ok := f.frrk8sPodForSpeakerPod[name]
	if !ok {
		return nil, fmt.Errorf("speaker %s/%s does not have a matching frr-k8s", ns, name)
	}

	return executor.ForPod(ns, frrk8s.Name, "frr"), nil
}

func (f frrk8sModeProvider) BGPMetricsPodForSpeaker(ns, name string) (*corev1.Pod, string, error) {
	p, ok := f.frrk8sPodForSpeakerPod[name]
	if !ok {
		return nil, "", fmt.Errorf("speakers %s/%s not found in map %v", ns, name, f.frrk8sPodForSpeakerPod)
	}

	return p, "frrk8s", nil
}
