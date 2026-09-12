// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

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

func TestAccKubernetesFleetAutoUpgradeProfile_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "testlist")
	r := KubernetesFleetAutoUpgradeProfileResource{}
	resourceName := fmt.Sprintf("acctestfaup-%d", data.RandomInteger)
	fleetName := fmt.Sprintf("acctestkfm-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctest-rg-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.listConfig(data)},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_kubernetes_fleet_auto_upgrade_profile.list", 1),
					querycheck.ExpectIdentity("azurerm_kubernetes_fleet_auto_upgrade_profile.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"fleet_name":          knownvalue.StringExact(fleetName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r KubernetesFleetAutoUpgradeProfileResource) listConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "test" {
  name                        = "acctestfaup-%[2]d"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) listQuery() string {
	return `
list "azurerm_kubernetes_fleet_auto_upgrade_profile" "list" {
  provider = azurerm
  config {
    kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  }
}
`
}
