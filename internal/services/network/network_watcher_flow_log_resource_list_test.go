// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func testAccNetworkWatcherFlowLog_listByNetworkWatcherID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "list")
	r := NetworkWatcherFlowLogResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_network_watcher_flow_log.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_network_watcher_flow_log.list",
						map[string]knownvalue.Check{
							"name":                 knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"network_watcher_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":      knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r NetworkWatcherFlowLogResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%[1]d"
  location = "%[2]s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

# Provisioning 3 virtual networks as the network watcher flow log must target distinct resources
resource "azurerm_virtual_network" "test" {
  count               = 3
  name                = "acctestvn-%[1]d${count.index}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test[0].name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-NW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsa%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_tier               = "Standard"
  account_kind               = "StorageV2"
  account_replication_type   = "LRS"
  https_traffic_only_enabled = true
}

resource "azurerm_network_watcher_flow_log" "test" {
  count                = 3
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[1]d${count.index}"

  target_resource_id = azurerm_virtual_network.test[count.index].id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

func (r NetworkWatcherFlowLogResource) basicQuery() string {
	return `
list "azurerm_network_watcher_flow_log" "list" {
  provider = azurerm
  config {
    network_watcher_id = azurerm_network_watcher.test.id
  }
}
`
}
