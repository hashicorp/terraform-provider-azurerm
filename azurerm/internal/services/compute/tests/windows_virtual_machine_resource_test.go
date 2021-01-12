package tests

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type WindowsVirtualMachineResource struct {
}

func (t WindowsVirtualMachineResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VMClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Windows Virtual Machine %q", id)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (WindowsVirtualMachineResource) templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "acctestvm%s"
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
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) template(data acceptance.TestData) string {
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
`, r.templateBase(data), data.RandomInteger)
}
