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
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

func TestAccKubernetesFleetAutoUpgradeProfile_listMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_auto_upgrade_profile", "test")
	r := KubernetesFleetAutoUpgradeProfileResource{}
	listAddress := "azurerm_kubernetes_fleet_auto_upgrade_profile.list"
	checks := []querycheck.QueryResultCheck{querycheck.ExpectLength(listAddress, 2)}
	for name, channel := range map[string]string{"stable": "Stable", "rapid": "Rapid"} {
		profileName := fmt.Sprintf("acctestfaup-%d-%s", data.RandomInteger, name)
		checks = append(checks,
			querycheck.ExpectIdentity(listAddress, map[string]knownvalue.Check{
				"name":                knownvalue.StringExact(profileName),
				"fleet_name":          knownvalue.StringExact(fmt.Sprintf("acctestkfm-%d", data.RandomInteger)),
				"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctest-rg-%d", data.RandomInteger)),
				"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
			}),
			querycheck.ExpectResourceKnownValues(listAddress, queryfilter.ByDisplayName(knownvalue.StringExact(profileName)), []querycheck.KnownValueCheck{
				{Path: tfjsonpath.New("name"), KnownValue: knownvalue.StringExact(profileName)},
				{Path: tfjsonpath.New("channel"), KnownValue: knownvalue.StringExact(channel)},
				{Path: tfjsonpath.New("enabled"), KnownValue: knownvalue.Bool(false)},
				{Path: tfjsonpath.New("node_image_selection_type"), KnownValue: knownvalue.StringExact("")},
				{Path: tfjsonpath.New("update_strategy_id"), KnownValue: knownvalue.StringExact("")},
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.listMultipleConfig(data)},
			{
				Query:             true,
				Config:            r.listQueryWithResource(),
				QueryResultChecks: checks,
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

func (r KubernetesFleetAutoUpgradeProfileResource) listMultipleConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "stable" {
  name                        = "acctestfaup-%[2]d-stable"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Stable"
  enabled                     = false
}

resource "azurerm_kubernetes_fleet_auto_upgrade_profile" "rapid" {
  name                        = "acctestfaup-%[2]d-rapid"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  channel                     = "Rapid"
  enabled                     = false
}
`, r.template(data), data.RandomInteger)
}

func (r KubernetesFleetAutoUpgradeProfileResource) listQueryWithResource() string {
	return `
list "azurerm_kubernetes_fleet_auto_upgrade_profile" "list" {
  provider         = azurerm
  include_resource = true
  config {
    kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.test.id
  }
}
`
}
