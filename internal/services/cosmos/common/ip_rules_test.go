// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

var (
	ipAddressOne = "127.0.0.1/32"
	ipAddressTwo = "168.63.129.16/32"
)

func TestCosmosDBIpRulesToIpRangeFilter(t *testing.T) {
	testData := []struct {
		Name     string
		Input    *[]cosmosdb.IPAddressOrRange
		Expected []string
	}{
		{
			Name:     "Nil",
			Input:    nil,
			Expected: []string{},
		},
		{
			Name: "One element",
			Input: &[]cosmosdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
			},
			Expected: []string{"127.0.0.1/32"},
		},
		{
			Name: "Multiple elements",
			Input: &[]cosmosdb.IPAddressOrRange{
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
		Expected *[]cosmosdb.IPAddressOrRange
	}{
		{
			Name:     "Empty",
			Input:    []string{},
			Expected: &[]cosmosdb.IPAddressOrRange{},
		},
		{
			Name:  "One element",
			Input: []string{ipAddressOne},
			Expected: &[]cosmosdb.IPAddressOrRange{
				{IPAddressOrRange: &ipAddressOne},
			},
		},
		{
			Name:  "Multiple elements",
			Input: []string{ipAddressOne, ipAddressTwo},
			Expected: &[]cosmosdb.IPAddressOrRange{
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
