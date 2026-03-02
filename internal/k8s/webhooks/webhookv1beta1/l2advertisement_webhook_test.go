// SPDX-License-Identifier:Apache-2.0

package webhookv1beta1

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateL2Advertisement(t *testing.T) {
	MetalLBNamespace = MetalLBTestNameSpace
	l2Adv := metallbv1.L2Advertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-l2adv",
			Namespace: MetalLBTestNameSpace,
		},
	}

	Logger = log.NewNopLogger()

	toRestore := getExistingL2Advs
	getExistingL2Advs = func() (*metallbv1.L2AdvertisementList, error) {
		return &metallbv1.L2AdvertisementList{
			Items: []metallbv1.L2Advertisement{
				l2Adv,
			},
		}, nil
	}
	toRestoreIPAddressPools := getExistingIPAddressPools
	getExistingIPAddressPools = func() (*metallbv1.IPAddressPoolList, error) {
		return &metallbv1.IPAddressPoolList{}, nil
	}

	defer func() {
		getExistingL2Advs = toRestore
		getExistingIPAddressPools = toRestoreIPAddressPools
	}()

	tests := []struct {
		desc         string
		l2Adv        *metallbv1.L2Advertisement
		isNew        bool
		failValidate bool
		expected     *metallbv1.L2AdvertisementList
	}{
		{
			desc: "Second Adv",
			l2Adv: &metallbv1.L2Advertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: true,
			expected: &metallbv1.L2AdvertisementList{
				Items: []metallbv1.L2Advertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-l2adv",
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
			l2Adv: &metallbv1.L2Advertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-l2adv",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: false,
			expected: &metallbv1.L2AdvertisementList{
				Items: []metallbv1.L2Advertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-l2adv",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same, new",
			l2Adv: &metallbv1.L2Advertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-l2adv",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: true,
			expected: &metallbv1.L2AdvertisementList{
				Items: []metallbv1.L2Advertisement{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-l2adv",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
			failValidate: true,
		},
		{
			desc: "Validation must fail if created in different namespace",
			l2Adv: &metallbv1.L2Advertisement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-l2adv1",
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
			err = validateL2AdvCreate(test.l2Adv)
		} else {
			err = validateL2AdvUpdate(test.l2Adv, nil)
		}
		if test.failValidate && err == nil {
			t.Fatalf("test %s failed, expecting error", test.desc)
		}
		if !cmp.Equal(test.expected, mock.l2Advs) {
			t.Fatalf("test %s failed, %s", test.desc, cmp.Diff(test.expected, mock.l2Advs))
		}
	}
}
