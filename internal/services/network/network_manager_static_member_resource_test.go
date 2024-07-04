// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/staticmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerStaticMemberResource struct{}

func testAccNetworkManagerStaticMember_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func testAccNetworkManagerStaticMember_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_static_member", "test")
	r := ManagerStaticMemberResource{}
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

func (r ManagerStaticMemberResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := staticmembers.ParseStaticMemberID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.StaticMembers
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerStaticMemberResource) template(data acceptance.TestData) string {
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
  name                = "acctest-nm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
}

resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/22"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ManagerStaticMemberResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_static_member" "test" {
  name                      = "acctest-nmsm-%d"
  network_group_id          = azurerm_network_manager_network_group.test.id
  target_virtual_network_id = azurerm_virtual_network.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerStaticMemberResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_static_member" "import" {
  name                      = azurerm_network_manager_static_member.test.name
  network_group_id          = azurerm_network_manager_static_member.test.network_group_id
  target_virtual_network_id = azurerm_network_manager_static_member.test.target_virtual_network_id
}
`, config)
}
