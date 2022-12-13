package helpers_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestMergeUserAppSettings(t *testing.T) {
	cases := []struct {
		service  []web.NameValuePair
		user     map[string]string
		expected []web.NameValuePair // Note: The function doesn't preserve the order of the list
	}{
		{
			service: []web.NameValuePair{
				{
					Name:  utils.String("test"),
					Value: utils.String("ServiceValue"),
				},
			},
			user: map[string]string{
				"test": "UserValue",
			},
			expected: []web.NameValuePair{
				{
					Name:  utils.String("test"),
					Value: utils.String("UserValue"),
				},
			},
		},
		{
			service: []web.NameValuePair{
				{
					Name:  utils.String("test"),
					Value: utils.String("ServiceValue"),
				},
				{
					Name:  utils.String("test2"),
					Value: utils.String("ServiceValue2"),
				},
				{
					Name:  utils.String("test3"),
					Value: utils.String("ServiceValue3"),
				},
				{
					Name:  utils.String("test4"),
					Value: utils.String("ServiceValue4"),
				},
			},
			user: map[string]string{
				"test":  "UserValue",
				"test4": "UserValue4",
			},
			expected: []web.NameValuePair{
				{
					Name:  utils.String("test"),
					Value: utils.String("UserValue"),
				},
				{
					Name:  utils.String("test2"),
					Value: utils.String("ServiceValue2"),
				},
				{
					Name:  utils.String("test3"),
					Value: utils.String("ServiceValue3"),
				},
				{
					Name:  utils.String("test4"),
					Value: utils.String("UserValue4"),
				},
			},
		},
	}

	for _, v := range cases {
		actualRaw := helpers.MergeUserAppSettings(&v.service, v.user)
		actual := *actualRaw
		sort.Slice(actual, func(i, j int) bool {
			return *actual[i].Name < *actual[j].Name
		})
		if !reflect.DeepEqual(actual, v.expected) {
			t.Fatalf("expected %+v, got %+v", v.expected, actual)
		}
	}
}
