// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
)

var (
	ipAddressOne = "127.0.0.1/32"
	ipAddressTwo = "168.63.129.16/32"
)

func TestCosmosDBIpRulesToIpRangeFilter(t *testing.T) {
	testData := []struct {
		Name     string
		Input    *[]documentdb.IPAddressOrRange
		Expected []string
	}{
		{
			Name:     "Nil",
			Input:    nil,
			Expected: []string{},
		},
		{
			Name: "One element",
			Input: &[]documentdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
			},
			Expected: []string{"127.0.0.1/32"},
		},
		{
			Name: "Multiple elements",
			Input: &[]documentdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
				{IPAddressOrRange: &ipAddressTwo},
			},
			Expected: []string{"127.0.0.1/32", "168.63.129.16/32"},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := CosmosDBIpRulesToIpRangeFilter(v.Input)

		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %q but got %q", v.Expected, actual)
		}
	}
}

func TestCosmosDBIpRangeFilterToIpRules(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []string
		Expected *[]documentdb.IPAddressOrRange
	}{
		{
			Name:     "Empty",
			Input:    []string{},
			Expected: &[]documentdb.IPAddressOrRange{},
		},
		{
			Name:  "One element",
			Input: []string{ipAddressOne},
			Expected: &[]documentdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
			},
		},
		{
			Name:  "Multiple elements",
			Input: []string{ipAddressOne, ipAddressTwo},
			Expected: &[]documentdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
				{IPAddressOrRange: &ipAddressTwo},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual := CosmosDBIpRangeFilterToIpRules(v.Input)

		if !reflect.DeepEqual(actual, v.Expected) {
			t.Fatalf("Expected %+v but got %+v", v.Expected, actual)
		}
	}
}
