// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SubscriptionZonePeersDataSource struct{}

func TestAccSubscriptionZonePeersDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription_zone_peers", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SubscriptionZonePeersDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.#").Exists(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.availability_zone").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.#").Exists(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.0.subscription_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("availability_zone_peers.0.peers.0.availability_zone").IsNotEmpty(),
			),
		},
	})
}

func (d SubscriptionZonePeersDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription_zone_peers" "test" {
  location             = %q
  peer_subscription_id = data.azurerm_client_config.current.subscription_id
}
`, data.Locations.Primary)
}
