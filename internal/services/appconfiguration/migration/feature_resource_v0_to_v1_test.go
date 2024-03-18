// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestFeatureResourceV0ToV1(t *testing.T) {
	testData := []struct {
		name                        string
		input                       map[string]interface{}
		expected                    *string
		appConfigurationEnvironment environments.Api
	}{
		{
			name: "old id (normal)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/labelName",
			},
			expected:                    utils.String("https://appConf1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label=labelName"),
			appConfigurationEnvironment: environments.AzurePublic().AppConfiguration,
		},
		{
			name: "old id (complicated)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/key:name/test/Label/test:label/name",
			},
			expected:                    utils.String("https://appConf1.azconfig.io/kv/.appconfig.featureflag%2Fkey:name%2Ftest?label=test%3Alabel%2Fname"),
			appConfigurationEnvironment: environments.AzurePublic().AppConfiguration,
		},
		{
			name: "old id (no label)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/%00",
			},
			expected:                    utils.String("https://appConf1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label="),
			appConfigurationEnvironment: environments.AzurePublic().AppConfiguration,
		},
		{
			name: "old id (\000 label)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/\000",
			},
			expected:                    utils.String("https://appConf1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label="),
			appConfigurationEnvironment: environments.AzurePublic().AppConfiguration,
		},
		{
			name: "old id (empty label)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/",
			},
			expected:                    utils.String("https://appConf1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label="),
			appConfigurationEnvironment: environments.AzurePublic().AppConfiguration,
		},
		{
			name: "old id (empty label - china)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/",
			},
			expected:                    utils.String("https://appConf1.azconfig.azure.cn/kv/.appconfig.featureflag%2FkeyName?label="),
			appConfigurationEnvironment: environments.AzureChina().AppConfiguration,
		},
		{
			name: "old id (empty label - usgov)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/keyName/Label/",
			},
			expected:                    utils.String("https://appConf1.azconfig.azure.us/kv/.appconfig.featureflag%2FkeyName?label="),
			appConfigurationEnvironment: environments.AzureUSGovernment().AppConfiguration,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		client := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: environments.Environment{
					AppConfiguration: test.appConfigurationEnvironment,
				},
			},
		}
		result, err := FeatureResourceV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, client)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		actualId := result["id"].(string)
		if *test.expected != actualId {
			t.Fatalf("expected %q but got %q", *test.expected, actualId)
		}
	}
}
