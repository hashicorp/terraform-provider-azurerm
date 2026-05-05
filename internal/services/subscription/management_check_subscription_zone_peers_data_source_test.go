// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagementCheckSubscriptionZonePeersDataSource struct{}

func TestAccManagementCheckSubscriptionZonePeersDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_check_subscription_zone_peers", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ManagementCheckSubscriptionZonePeersDataSource{}.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.#").HasValue("3"),
				check.That(data.ResourceName).Key("availability_zone_peers.0.availability_zone").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.#").HasValue("1"),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.0.subscription_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.0.availability_zone").IsNotEmpty(),
			),
		},
	})
}

func TestAccManagementCheckSubscriptionZonePeersDataSource_fullResourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_check_subscription_zone_peers", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ManagementCheckSubscriptionZonePeersDataSource{}.fullResourceId(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.#").HasValue("3"),
			),
		},
	})
}

func (d ManagementCheckSubscriptionZonePeersDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_check_subscription_zone_peers" "test" {
  location             = "eastus"
  peer_subscription_id = data.azurerm_client_config.current.subscription_id
}
`
}

func (d ManagementCheckSubscriptionZonePeersDataSource) fullResourceId() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_check_subscription_zone_peers" "test" {
  location             = "eastus"
  peer_subscription_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
}
`
}
