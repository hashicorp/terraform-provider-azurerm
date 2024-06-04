// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagerconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerManagementGroupConnectionResource struct{}

func testAccNetworkManagerManagementGroupConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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

func testAccNetworkManagerManagementGroupConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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

func testAccNetworkManagerManagementGroupConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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

func testAccNetworkManagerManagementGroupConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_management_group_connection", "test")
	r := ManagerManagementGroupConnectionResource{}
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

func (r ManagerManagementGroupConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkmanagerconnections.ParseProviders2NetworkManagerConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkManagerConnections
	resp, err := client.ManagementGroupNetworkManagerConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerManagementGroupConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_management_group_subscription_association" "test" {
  management_group_id = azurerm_management_group.test.id
  subscription_id     = data.azurerm_subscription.alt.id
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}

data "azurerm_subscription" "alt" {
  subscription_id = %q
}

data "azurerm_subscription" "current" {
}

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "network_contributor" {
  scope                = azurerm_management_group.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_network_manager" "test" {
  name                = "acctest-networkmanager-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
  depends_on     = [azurerm_role_assignment.network_contributor]
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_SUBSCRIPTION_ID_ALT"))
}

func (r ManagerManagementGroupConnectionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmmgc-%d"
  management_group_id = azurerm_management_group.test.id
  network_manager_id  = azurerm_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerManagementGroupConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_management_group_connection" "import" {
  name                = azurerm_network_manager_management_group_connection.test.name
  management_group_id = azurerm_network_manager_management_group_connection.test.management_group_id
  network_manager_id  = azurerm_network_manager_management_group_connection.test.network_manager_id
}
`, config)
}

func (r ManagerManagementGroupConnectionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmmgc-%d"
  management_group_id = azurerm_management_group.test.id
  network_manager_id  = azurerm_network_manager.test.id
  description         = "complete"
}
`, template, data.RandomInteger)
}

func (r ManagerManagementGroupConnectionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_management_group_connection" "test" {
  name                = "acctest-nmmgc-%d"
  management_group_id = azurerm_management_group.test.id
  network_manager_id  = azurerm_network_manager.test.id
  description         = "update"
}
`, template, data.RandomInteger)
}
