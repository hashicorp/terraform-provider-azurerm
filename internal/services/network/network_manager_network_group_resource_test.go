// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerNetworkGroupResource struct{}

func testAccNetworkManagerNetworkGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_network_group", "test")
	r := ManagerNetworkGroupResource{}
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

func testAccNetworkManagerNetworkGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_network_group", "test")
	r := ManagerNetworkGroupResource{}
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

func testAccNetworkManagerNetworkGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_network_group", "test")
	r := ManagerNetworkGroupResource{}
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

func testAccNetworkManagerNetworkGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_network_group", "test")
	r := ManagerNetworkGroupResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (r ManagerNetworkGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkgroups.ParseNetworkGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.NetworkGroups
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerNetworkGroupResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ManagerNetworkGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s
resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
}
`, template, data.RandomInteger)
}

func (r ManagerNetworkGroupResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s
resource "azurerm_network_manager_network_group" "import" {
  name               = azurerm_network_manager_network_group.test.name
  network_manager_id = azurerm_network_manager_network_group.test.network_manager_id
}
`, config)
}

func (r ManagerNetworkGroupResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s
resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
  description        = "test complete"
}
`, template, data.RandomInteger)
}

func (r ManagerNetworkGroupResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s
resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
  description        = "test update"
}
`, template, data.RandomInteger)
}
