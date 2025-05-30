// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterProjectPoolTestResource struct{}

func TestAccDevCenterProjectPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_pool", "test")
	r := DevCenterProjectPoolTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevCenterProjectPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_pool", "test")
	r := DevCenterProjectPoolTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDevCenterProjectPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_pool", "test")
	r := DevCenterProjectPoolTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevCenterProjectPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project_pool", "test")
	r := DevCenterProjectPoolTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
			Config: r.update(data),
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

func (r DevCenterProjectPoolTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pools.ParsePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.Pools.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterProjectPoolTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project_pool" "test" {
  name                             = "acctest-dcpl-%d"
  location                         = azurerm_resource_group.test.location
  dev_center_project_id            = azurerm_dev_center_project.test.id
  dev_box_definition_name          = azurerm_dev_center_dev_box_definition.test.name
  local_administrator_enabled      = false
  dev_center_attached_network_name = azurerm_dev_center_attached_network.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterProjectPoolTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_project_pool" "import" {
  name                             = azurerm_dev_center_project_pool.test.name
  location                         = azurerm_dev_center_project_pool.test.location
  dev_center_project_id            = azurerm_dev_center_project_pool.test.dev_center_project_id
  dev_box_definition_name          = azurerm_dev_center_project_pool.test.dev_box_definition_name
  local_administrator_enabled      = azurerm_dev_center_project_pool.test.local_administrator_enabled
  dev_center_attached_network_name = azurerm_dev_center_project_pool.test.dev_center_attached_network_name
}
`, r.basic(data))
}

func (r DevCenterProjectPoolTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project_pool" "test" {
  name                                    = "acctest-dcpl-%d"
  location                                = azurerm_resource_group.test.location
  dev_center_project_id                   = azurerm_dev_center_project.test.id
  dev_box_definition_name                 = azurerm_dev_center_dev_box_definition.test.name
  local_administrator_enabled             = true
  dev_center_attached_network_name        = azurerm_dev_center_attached_network.test.name
  stop_on_disconnect_grace_period_minutes = 60

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterProjectPoolTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project_pool" "test" {
  name                                    = "acctest-dcpl-%d"
  location                                = azurerm_resource_group.test.location
  dev_center_project_id                   = azurerm_dev_center_project.test.id
  dev_box_definition_name                 = azurerm_dev_center_dev_box_definition.test2.name
  local_administrator_enabled             = false
  dev_center_attached_network_name        = azurerm_dev_center_attached_network.test2.name
  stop_on_disconnect_grace_period_minutes = 80

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterProjectPoolTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-dcet-%d"
  location = "%s"
}

resource "azurerm_dev_center" "test" {
  name                = "acctest-dc-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dev_center_dev_box_definition" "test" {
  name               = "acctest-dcet-%d"
  location           = azurerm_resource_group.test.location
  dev_center_id      = azurerm_dev_center.test.id
  image_reference_id = "${azurerm_dev_center.test.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_dev_center_network_connection" "test" {
  name                = "acctest-dcnc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  domain_join_type    = "AzureADJoin"
  subnet_id           = azurerm_subnet.test.id
}

resource "azurerm_dev_center_attached_network" "test" {
  name                  = "acctest-dcan-%d"
  dev_center_id         = azurerm_dev_center.test.id
  network_connection_id = azurerm_dev_center_network_connection.test.id

  depends_on = [azurerm_dev_center_dev_box_definition.test]
}

resource "azurerm_dev_center_dev_box_definition" "test2" {
  name               = "acctest-dcet2-%d"
  location           = azurerm_resource_group.test.location
  dev_center_id      = azurerm_dev_center.test.id
  image_reference_id = "${azurerm_dev_center.test.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"
}

resource "azurerm_dev_center_attached_network" "test2" {
  name                  = "acctest-dcan2-%d"
  dev_center_id         = azurerm_dev_center.test.id
  network_connection_id = azurerm_dev_center_network_connection.test.id

  depends_on = [azurerm_dev_center_dev_box_definition.test2]
}

resource "azurerm_dev_center_project" "test" {
  name                = "acctest-dcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id

  depends_on = [azurerm_dev_center_attached_network.test, azurerm_dev_center_attached_network.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
