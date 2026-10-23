// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

type LinuxVirtualMachineResource struct{}

func (r LinuxVirtualMachineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualmachines.ParseVirtualMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VirtualMachinesClient.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LinuxVirtualMachineResource) imageFromExistingMachineDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  first_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9ddAwoR0XBoT9kixLX6atgX9dovt9fpR1HO/R9jYwYnuB+SZ845KSqat+U0m6oagZhpsfcEEwjGGjQz6Z1rB6mvffsKq6i74cmm0jO564nBnZQeh31q3sFNs+XdrDtFmnYRqdPHhhr1sw0C/rxbiaE6nYZWRfHW//81nEePKMpjiN8JsrYQNbzEpz8QOBSquwBmXO+LVx//zAbY4jGTa4hjGeNzIgMJZ8Jk/11XbcxSK1PK43BrejHg6kctmEkYvMH/o12RfAeB8okGCRW3scwOozxVrHwxaPgEf03jig+Ag9V+GXNBabL5AWtxcuPN63rUfaAXEIXTHmndwVOxlpLrUf5ox1+ddGyWbLMXzd7akPioof5MNJMq/yuFGC5dY0Z6/+yGRNtShQesVo/czhKEPGIcsIi5gnKdfDB4i9ay2yz8ystnW6jbabcyqejk1Qc61wapaFdhUHL0iD/GW/5ZujDs5C3BT7EIgKLIfAaAx5TBEJyE1KQ/GEOifB8ztDl/gp99o+i2HKABtmYv12y4JVlEUkRckeLrw6luEb3ColHshsQcQGfudGFFgdEdcgBrV4Ch7IkLxVYQl3pegzZiirMPnRKh10r/Hrg6uYxn7sLeTJoD5VOKmqmeK4kFXsZMVtA6/SnxQtUKkKlfLBwBSDrrdgLjBV+KOndiwC7Q=="
  second_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  vm_name           = "acctestsourcevm-%[1]d"
  admin_username    = "testadmin%[1]d"
  admin_password    = "Password1234!%[1]d"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
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

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = local.vm_name
}

resource "azurerm_network_interface" "public" {
  name                = "acctnicsource-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinuxVirtualMachineResource) imageFromExistingMachinePrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = local.admin_username
  disable_password_authentication = false
  admin_password                  = local.admin_password

  network_interface_ids = [
    azurerm_network_interface.public.id,
  ]

  admin_ssh_key {
    username   = local.admin_username
    public_key = local.first_public_key
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

  custom_data = base64encode(<<-EOT
#!/bin/bash
sudo waagent -verbose -deprovision+user -force
EOT
  )
}
`, r.imageFromExistingMachineDependencies(data), data.RandomInteger)
}
