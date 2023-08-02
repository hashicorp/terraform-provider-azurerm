// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/scopeconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerScopeConnectionResource struct{}

func testAccNetworkManagerScopeConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func testAccNetworkManagerScopeConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func testAccNetworkManagerScopeConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func testAccNetworkManagerScopeConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_scope_connection", "test")
	r := ManagerScopeConnectionResource{}
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

func (r ManagerScopeConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scopeconnections.ParseScopeConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ScopeConnections
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerScopeConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
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

func (r ManagerScopeConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_scope_connection" "test" {
  name               = "acctest-nsc-%d"
  network_manager_id = azurerm_network_manager.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  target_scope_id    = data.azurerm_subscription.current.id
}
`, template, data.RandomInteger)
}

func (r ManagerScopeConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_scope_connection" "import" {
  name               = azurerm_network_manager_scope_connection.test.name
  network_manager_id = azurerm_network_manager_scope_connection.test.network_manager_id
  tenant_id          = azurerm_network_manager_scope_connection.test.tenant_id
  target_scope_id    = azurerm_network_manager_scope_connection.test.target_scope_id
}
`, config)
}

func (r ManagerScopeConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_scope_connection" "test" {
  name               = "acctest-nsc-%d"
  network_manager_id = azurerm_network_manager.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  target_scope_id    = data.azurerm_subscription.current.id
  description        = "complete"
}
`, template, data.RandomInteger)
}

func (r ManagerScopeConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_scope_connection" "test" {
  name               = "acctest-nsc-%d"
  network_manager_id = azurerm_network_manager.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  target_scope_id    = data.azurerm_subscription.current.id
  description        = "update"
}
`, template, data.RandomInteger)
}
