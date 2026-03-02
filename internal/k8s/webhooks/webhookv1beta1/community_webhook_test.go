// SPDX-License-Identifier:Apache-2.0

package webhookv1beta1

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateCommunity(t *testing.T) {
	MetalLBNamespace = MetalLBTestNameSpace
	Logger = log.NewNopLogger()

	toRestoreCommunities := getExistingCommunities
	getExistingCommunities = func() (*metallbv1.CommunityList, error) {
		return &metallbv1.CommunityList{
			Items: []metallbv1.Community{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-commuinty1",
						Namespace: MetalLBTestNameSpace,
					},
				},
			},
		}, nil
	}
	defer func() {
		getExistingCommunities = toRestoreCommunities
	}()

	tests := []struct {
		desc           string
		commuinty      *metallbv1.Community
		isNewCommunity bool
		failValidate   bool
		expected       *metallbv1.CommunityList
	}{
		{
			desc: "Second Community",
			commuinty: &metallbv1.Community{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-community2",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNewCommunity: true,
			expected: &metallbv1.CommunityList{
				Items: []metallbv1.Community{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-commuinty1",
							Namespace: MetalLBTestNameSpace,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-community2",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same Community, update",
			commuinty: &metallbv1.Community{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-commuinty1",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNewCommunity: false,
			expected: &metallbv1.CommunityList{
				Items: []metallbv1.Community{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-commuinty1",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
		},
		{
			desc: "Same community, new",
			commuinty: &metallbv1.Community{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-commuinty1",
					Namespace: MetalLBTestNameSpace,
				},
			},
			isNewCommunity: true,
			expected: &metallbv1.CommunityList{
				Items: []metallbv1.Community{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-commuinty1",
							Namespace: MetalLBTestNameSpace,
						},
					},
				},
			},
			failValidate: true,
		},
		{
			desc: "Validation must fail if created in different namespace",
			commuinty: &metallbv1.Community{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-commuinty2",
					Namespace: "default",
				},
			},
			isNewCommunity: true,
			expected:       nil,
			failValidate:   true,
		},
	}
	for _, test := range tests {
		var err error
		mock := &mockValidator{}
		Validator = mock
		mock.forceError = test.failValidate

		if test.isNewCommunity {
			err = validateCommunityCreate(test.commuinty)
		} else {
			err = validateCommunityUpdate(test.commuinty, nil)
		}
		if test.failValidate && err == nil {
			t.Fatalf("test %s failed, expecting error", test.desc)
		}
		if !cmp.Equal(test.expected, mock.communities) {
			t.Fatalf("test %s failed, %s", test.desc, cmp.Diff(test.expected, mock.communities))
		}
	}
}
