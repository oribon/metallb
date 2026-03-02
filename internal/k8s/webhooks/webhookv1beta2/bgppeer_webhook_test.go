// SPDX-License-Identifier:Apache-2.0

package webhookv1beta2

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const testNamespace = "namespace"

func TestValidateBGPPeer(t *testing.T) {
	MetalLBNamespace = testNamespace
	bgpPeer := metallbv1.BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-peer",
			Namespace: testNamespace,
		},
	}

	Logger = log.NewNopLogger()

	toRestore := GetExistingBGPPeers
	GetExistingBGPPeers = func() (*metallbv1.BGPPeerList, error) {
		return &metallbv1.BGPPeerList{
			Items: []metallbv1.BGPPeer{
				bgpPeer,
			},
		}, nil
	}

	defer func() {
		GetExistingBGPPeers = toRestore
	}()

	tests := []struct {
		desc         string
		bgpPeer      *metallbv1.BGPPeer
		isNew        bool
		failValidate bool
		expected     *metallbv1.BGPPeerList
	}{
		{
			desc: "Second Peer",
			bgpPeer: &metallbv1.BGPPeer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: testNamespace,
				},
			},
			isNew: true,
			expected: &metallbv1.BGPPeerList{
				Items: []metallbv1.BGPPeer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-peer",
							Namespace: testNamespace,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test",
							Namespace: testNamespace,
						},
					},
				},
			},
		},
		{
			desc: "Same, update",
			bgpPeer: &metallbv1.BGPPeer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-peer",
					Namespace: testNamespace,
				},
			},
			isNew: false,
			expected: &metallbv1.BGPPeerList{
				Items: []metallbv1.BGPPeer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-peer",
							Namespace: testNamespace,
						},
					},
				},
			},
		},
		{
			desc: "Validation failed",
			bgpPeer: &metallbv1.BGPPeer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-peer",
					Namespace: testNamespace,
				},
			},
			isNew: false,
			expected: &metallbv1.BGPPeerList{
				Items: []metallbv1.BGPPeer{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-peer",
							Namespace: testNamespace,
						},
					},
				},
			},
			failValidate: true,
		},
		{
			desc: "Validation must fail if created in different namespace",
			bgpPeer: &metallbv1.BGPPeer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-peer1",
					Namespace: "default",
				},
			},
			isNew:        true,
			expected:     nil,
			failValidate: true,
		},
	}
	for _, test := range tests {
		var err error
		mock := &mockValidator{}
		Validator = mock
		mock.forceError = test.failValidate
		if test.isNew {
			_, err = validatePeerCreate(test.bgpPeer)
		} else {
			err = validatePeerUpdate(test.bgpPeer, nil)
		}
		if test.failValidate && err == nil {
			t.Fatalf("test %s failed, expecting error", test.desc)
		}
		if !cmp.Equal(test.expected, mock.bgpPeers) {
			t.Fatalf("test %s failed, %s", test.desc, cmp.Diff(test.expected, mock.bgpPeers))
		}
	}
}
