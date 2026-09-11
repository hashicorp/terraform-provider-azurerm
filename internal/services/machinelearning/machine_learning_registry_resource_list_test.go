// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMachineLearningRegistry_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryResource{}
	resourceName := fmt.Sprintf("acctestMLR-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-ml-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.subscriptionListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_machine_learning_registry.list", 1),
					querycheck.ExpectIdentity("azurerm_machine_learning_registry.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_machine_learning_registry.list", 1),
					querycheck.ExpectIdentity("azurerm_machine_learning_registry.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r MachineLearningRegistryResource) subscriptionListQuery() string {
	return `
list "azurerm_machine_learning_registry" "list" {
  provider = azurerm
}
`
}

func (r MachineLearningRegistryResource) resourceGroupListQuery() string {
	return `
list "azurerm_machine_learning_registry" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
