// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func TestMergeUserAppSettings(t *testing.T) {
	cases := []struct {
		service  []webapps.NameValuePair
		user     map[string]string
		expected []webapps.NameValuePair // Note: The function doesn't preserve the order of the list
	}{
		{
			service: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("ServiceValue"),
				},
			},
			user: map[string]string{
				"test": "UserValue",
			},
			expected: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("UserValue"),
				},
			},
		},
		{
			service: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("ServiceValue"),
				},
				{
					Name:  pointer.To("test2"),
					Value: pointer.To("ServiceValue2"),
				},
				{
					Name:  pointer.To("test3"),
					Value: pointer.To("ServiceValue3"),
				},
				{
					Name:  pointer.To("test4"),
					Value: pointer.To("ServiceValue4"),
				},
			},
			user: map[string]string{
				"test":  "UserValue",
				"test4": "UserValue4",
			},
			expected: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("UserValue"),
				},
				{
					Name:  pointer.To("test2"),
					Value: pointer.To("ServiceValue2"),
				},
				{
					Name:  pointer.To("test3"),
					Value: pointer.To("ServiceValue3"),
				},
				{
					Name:  pointer.To("test4"),
					Value: pointer.To("UserValue4"),
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
