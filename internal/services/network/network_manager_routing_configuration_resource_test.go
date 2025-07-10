// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkmanagerroutingconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerRoutingConfigurationResource struct{}

func testAccNetworkManagerRoutingConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_configuration", "test")
	r := ManagerRoutingConfigurationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerRoutingConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_configuration", "test")
	r := ManagerRoutingConfigurationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerRoutingConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_configuration", "test")
	r := ManagerRoutingConfigurationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkManagerRoutingConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_configuration", "test")
	r := ManagerRoutingConfigurationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerRoutingConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkmanagerroutingconfigurations.ParseRoutingConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.NetworkManagerRoutingConfigurations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerRoutingConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_configuration" "test" {
  name               = "acctest-nmrc-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_routing_configuration" "import" {
  name               = azurerm_network_manager_routing_configuration.test.name
  network_manager_id = azurerm_network_manager.test.id
}
`, r.basic(data))
}

func (r ManagerRoutingConfigurationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_routing_configuration" "test" {
  name               = "acctest-nmrc-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  description        = "This is a test Routing Configuration"
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerRoutingConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-rc-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-rc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Routing"]
}
`, data.RandomInteger, data.Locations.Primary)
}
