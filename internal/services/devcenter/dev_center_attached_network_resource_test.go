// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterAttachedNetworkTestResource struct{}

func TestAccDevCenterAttachedNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_attached_network", "test")
	r := DevCenterAttachedNetworkTestResource{}

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

func TestAccDevCenterAttachedNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_attached_network", "test")
	r := DevCenterAttachedNetworkTestResource{}

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

func (r DevCenterAttachedNetworkTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := attachednetworkconnections.ParseDevCenterAttachedNetworkID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.AttachedNetworkConnections.AttachedNetworksGetByDevCenter(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterAttachedNetworkTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_attached_network" "test" {
  name                  = "acctest-dcet-%d"
  dev_center_id         = azurerm_dev_center.test.id
  network_connection_id = azurerm_dev_center_network_connection.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterAttachedNetworkTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_attached_network" "import" {
  name                  = azurerm_dev_center_attached_network.test.name
  dev_center_id         = azurerm_dev_center_attached_network.test.dev_center_id
  network_connection_id = azurerm_dev_center_attached_network.test.network_connection_id
}
`, r.basic(data))
}

func (r DevCenterAttachedNetworkTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-dcan-%d"
  location = "%s"
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

resource "azurerm_dev_center" "test" {
  name                = "acctest-dc-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_dev_center_network_connection" "test" {
  name                = "acctest-dcnc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  domain_join_type    = "AzureADJoin"
  subnet_id           = azurerm_subnet.test.id

  depends_on = [azurerm_dev_center.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}
