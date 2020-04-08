package tests

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func checkLinuxVirtualMachineIsDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if rs.Type != "azurerm_linux_virtual_machine" {
			continue
		}

		id, err := parse.VirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func checkLinuxVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Linux Virtual Machine %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on VMClient: %+v", err)
		}

		return nil
	}
}

func testLinuxVirtualMachine_templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  first_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testLinuxVirtualMachine_template(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}
