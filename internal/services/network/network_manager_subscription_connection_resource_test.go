// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerSubscriptionConnectionResource struct{}

func testAccNetworkSubscriptionNetworkManagerConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func testAccNetworkSubscriptionNetworkManagerConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func testAccNetworkSubscriptionNetworkManagerConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
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

func testAccNetworkSubscriptionNetworkManagerConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_subscription_connection", "test")
	r := ManagerSubscriptionConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerSubscriptionConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkmanagerconnections.ParseNetworkManagerConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkManagerConnections
	resp, err := client.SubscriptionNetworkManagerConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerSubscriptionConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "test" {
  name                = "acctest-networkmanager-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagerSubscriptionConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name               = "acctest-nmsc-%d"
  subscription_id    = data.azurerm_subscription.current.id
  network_manager_id = azurerm_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerSubscriptionConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_subscription_connection" "import" {
  name               = "acctest-nmsc-%d"
  subscription_id    = data.azurerm_subscription.current.id
  network_manager_id = azurerm_network_manager.test.id
}
`, config, data.RandomInteger)
}

func (r ManagerSubscriptionConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name               = "acctest-nmsc-%d"
  subscription_id    = data.azurerm_subscription.current.id
  network_manager_id = azurerm_network_manager.test.id
  description        = "complete"
}
`, template, data.RandomInteger)
}

func (r ManagerSubscriptionConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_subscription_connection" "test" {
  name               = "acctest-nmsc-%d"
  subscription_id    = data.azurerm_subscription.current.id
  network_manager_id = azurerm_network_manager.test.id
  description        = "update"
}
`, template, data.RandomInteger)
}
