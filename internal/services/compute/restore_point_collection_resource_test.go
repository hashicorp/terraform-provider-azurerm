// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RestorePointCollectionResource struct{}

func TestAccRestorePointCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_restore_point_collection", "test")
	r := RestorePointCollectionResource{}

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

func TestAccRestorePointCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_restore_point_collection", "test")
	r := RestorePointCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func (r RestorePointCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := restorepointcollections.ParseRestorePointCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.RestorePointCollectionsClient.Get(ctx, *id, restorepointcollections.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r RestorePointCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%[1]s
resource "azurerm_restore_point_collection" "test" {
  name                      = "acctestRPC-%[2]s"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_linux_virtual_machine.test.location
  source_virtual_machine_id = azurerm_linux_virtual_machine.test.id

  tags = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomString)
}

func (r RestorePointCollectionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%[1]s
resource "azurerm_restore_point_collection" "test" {
  name                      = "acctestRPC-%[2]s"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_linux_virtual_machine.test.location
  source_virtual_machine_id = azurerm_linux_virtual_machine.test.id

  tags = {
    bar = "foo"
  }
}
`, r.template(data), data.RandomString)
}

func (r RestorePointCollectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Compute-%[1]d"
  location = "%[2]s"
}


resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestinternal-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
