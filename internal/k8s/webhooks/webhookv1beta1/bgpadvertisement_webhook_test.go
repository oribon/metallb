// SPDX-License-Identifier:Apache-2.0

package webhookv1beta1

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	metallbv1 "go.universe.tf/metallb/api/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateBGPAdvertisement(t *testing.T) {
	MetalLBNamespace = MetalLBTestNameSpace
	bgpAdv := metallbv1.BGPAdvertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-bgpadv",
			Namespace: MetalLBTestNameSpace,
		},
	}

	Logger = log.NewNopLogger()

	toRestore := getExistingBGPAdvs
	getExistingBGPAdvs = func() (*metallbv1.BGPAdvertisementList, error) {
		return &metallbv1.BGPAdvertisementList{
			Items: []metallbv1.BGPAdvertisement{
				bgpAdv,
			},
		}, nil
	}
	toRestoreIPAddressPools := getExistingIPAddressPools
	getExistingIPAddressPools = func() (*metallbv1.IPAddressPoolList, error) {
		return &metallbv1.IPAddressPoolList{}, nil
	}
	toRestoreNodes := getExistingNodes
	getExistingNodes = func() (*v1core.NodeList, error) {
		return &v1core.NodeList{}, nil
	}

	defer func() {
		getExistingBGPAdvs = toRestore
		getExistingIPAddressPools = toRestoreIPAddressPools
		getExistingNodes = toRestoreNodes
	}()

	tests := []struct {
		desc         string
		bgpAdv       *metallbv1.BGPAdvertisement
		isNew        bool
		failValidate bool
		expected     *metallbv1.BGPAdvertisementList
	}{
		{
			desc: "Second Adv",
			bgpAdv: &metallbv1.BGPAdvertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: true,
			expected: &metallbv1.BGPAdvertisementList{
				Items: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-bgpadv",
							Namespace: MetalLBTestNameSpace,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same, update",
			bgpAdv: &metallbv1.BGPAdvertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-bgpadv",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: false,
			expected: &metallbv1.BGPAdvertisementList{
				Items: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-bgpadv",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same, new",
			bgpAdv: &metallbv1.BGPAdvertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-bgpadv",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: true,
			expected: &metallbv1.BGPAdvertisementList{
				Items: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-bgpadv",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
			failValidate: true,
		},
		{
			desc: "Validation must fail if created in different namespace",
			bgpAdv: &metallbv1.BGPAdvertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-bgpadv1",
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
			err = validateBGPAdvCreate(test.bgpAdv)
		} else {
			err = validateBGPAdvUpdate(test.bgpAdv, nil)
		}
		if test.failValidate && err == nil {
			t.Fatalf("test %s failed, expecting error", test.desc)
		}
		if !cmp.Equal(test.expected, mock.bgpAdvs) {
			t.Fatalf("test %s failed, %s", test.desc, cmp.Diff(test.expected, mock.bgpAdvs))
		}
	}
}
