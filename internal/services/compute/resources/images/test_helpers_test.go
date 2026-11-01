// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package images_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ImageResource struct{}

func (ImageResource) generalizeVirtualMachine() func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := virtualmachines.ParseVirtualMachineID(state.ID)
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		log.Printf("[DEBUG] Deallocating VM..")
		if err := client.Compute.VirtualMachinesClient.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
			return fmt.Errorf("Bad: deallocating %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Generalizing VM..")
		if _, err = client.Compute.VirtualMachinesClient.Generalize(ctx, *id); err != nil {
			return fmt.Errorf("Bad: Generalizing %s: %+v", *id, err)
		}

		return nil
	}
}

func (r ImageResource) setupManagedDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_linux_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  size                  = "Standard_D1_v2"

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk1${local.number}"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  computer_name                   = "mdimagetestsource"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_managed_disk" "testsource" {
  name                = azurerm_linux_virtual_machine.testsource.os_disk.0.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (ImageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  number            = "%d"
  location          = %q
  domain_name_label = "acctestvm-%s"
  random_string     = %q
  admin_username    = "testadmin%d"
  admin_password    = "Password1234!%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-${local.number}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = local.domain_name_label
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (ImageResource) virtualMachineExists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := virtualmachines.ParseVirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}
	resp, err := client.Compute.VirtualMachinesClient.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return nil
}
