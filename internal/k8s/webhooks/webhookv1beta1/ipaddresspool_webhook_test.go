// SPDX-License-Identifier:Apache-2.0

package webhookv1beta1

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateIPAddressPool(t *testing.T) {
	MetalLBNamespace = MetalLBTestNameSpace
	ipAddressPool := metallbv1.IPAddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ippool",
			Namespace: MetalLBTestNameSpace,
		},
	}
	Logger = log.NewNopLogger()

	toRestoreIPAddressPools := getExistingIPAddressPools
	getExistingIPAddressPools = func() (*metallbv1.IPAddressPoolList, error) {
		return &metallbv1.IPAddressPoolList{
			Items: []metallbv1.IPAddressPool{
				ipAddressPool,
			},
		}, nil
	}

	defer func() {
		getExistingIPAddressPools = toRestoreIPAddressPools
	}()

	tests := []struct {
		desc          string
		ipAddressPool *metallbv1.IPAddressPool
		isNew         bool
		failValidate  bool
		expected      *metallbv1.IPAddressPoolList
	}{
		{
			desc: "Second IPAddressPool",
			ipAddressPool: &metallbv1.IPAddressPool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ippool1",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: true,
			expected: &metallbv1.IPAddressPoolList{
				Items: []metallbv1.IPAddressPool{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-ippool",
							Namespace: MetalLBTestNameSpace,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-ippool1",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same IPAddressPool, update",
			ipAddressPool: &metallbv1.IPAddressPool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ippool",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: false,
			expected: &metallbv1.IPAddressPoolList{
				Items: []metallbv1.IPAddressPool{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-ippool",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Validation fails",
			ipAddressPool: &metallbv1.IPAddressPool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ippool",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNew: false,
			expected: &metallbv1.IPAddressPoolList{
				Items: []metallbv1.IPAddressPool{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-ippool",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
			failValidate: true,
		},
		{
			desc: "Validation must fail if created in different namespace",
			ipAddressPool: &metallbv1.IPAddressPool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ippool2",
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
			err = validateIPAddressPoolCreate(test.ipAddressPool)
		} else {
			err = validateIPAddressPoolUpdate(test.ipAddressPool, nil)
		}
		if test.failValidate && err == nil {
			t.Fatalf("test %s failed, expecting error", test.desc)
		}
		if !cmp.Equal(test.expected, mock.ipAddressPools) {
			t.Fatalf("test %s failed, %s", test.desc, cmp.Diff(test.expected, mock.ipAddressPools))
		}
	}
}
