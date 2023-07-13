// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DedicatedHardwareSecurityModuleResource struct{}

func TestAccDedicatedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")
	r := DedicatedHardwareSecurityModuleResource{}

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

func TestAccDedicatedHardwareSecurityModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")
	r := DedicatedHardwareSecurityModuleResource{}

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

func TestAccDedicatedHardwareSecurityModule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")
	r := DedicatedHardwareSecurityModuleResource{}

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

func TestAccDedicatedHardwareSecurityModule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")
	r := DedicatedHardwareSecurityModuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managementNetworkProfile(data),
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
	})
}

func (r DedicatedHardwareSecurityModuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "import" {
  name                = azurerm_dedicated_hardware_security_module.test.name
  resource_group_name = azurerm_dedicated_hardware_security_module.test.resource_group_name
  location            = azurerm_dedicated_hardware_security_module.test.location
  sku_name            = azurerm_dedicated_hardware_security_module.test.sku_name
  stamp_id            = azurerm_dedicated_hardware_security_module.test.stamp_id

  network_profile {
    network_interface_private_ip_addresses = azurerm_dedicated_hardware_security_module.test.network_profile[0].network_interface_private_ip_addresses
    subnet_id                              = azurerm_dedicated_hardware_security_module.test.network_profile[0].subnet_id
  }
}
`, r.basic(data))
}

func (DedicatedHardwareSecurityModuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dedicatedhsms.ParseDedicatedHSMID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HSM.DedicatedHsmClient.DedicatedHsmGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DedicatedHardwareSecurityModuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hsm-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.2.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-hsmsubnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.1.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.HardwareSecurityModules/dedicatedHSMs"

      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet" "test3" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.255.0/26"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-vnetgateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "ExpressRoute"
  vpn_type = "PolicyBased"
  sku      = "Standard"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test3.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DedicatedHardwareSecurityModuleResource) basic(data acceptance.TestData) string {
	template := DedicatedHardwareSecurityModuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "test" {
  name                = "acctest-hsm-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "SafeNet Luna Network HSM A790"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  stamp_id = "stamp2"

  depends_on = [azurerm_virtual_network_gateway.test]
}
`, template, data.RandomString)
}

func (DedicatedHardwareSecurityModuleResource) managementNetworkProfile(data acceptance.TestData) string {
	template := DedicatedHardwareSecurityModuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "test" {
  name                = "acctest-hsm-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "payShield10K_LMK1_CPS60"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  management_network_profile {
    network_interface_private_ip_addresses = ["10.2.1.9"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  stamp_id = "stamp2"

  depends_on = [azurerm_virtual_network_gateway.test]
}
`, template, data.RandomString)
}

func (DedicatedHardwareSecurityModuleResource) complete(data acceptance.TestData) string {
	template := DedicatedHardwareSecurityModuleResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "test" {
  name                = "acctest-hsm-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "payShield10K_LMK1_CPS60"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  management_network_profile {
    network_interface_private_ip_addresses = ["10.2.1.9"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  stamp_id = "stamp2"

  tags = {
    env = "Test"
  }

  depends_on = [azurerm_virtual_network_gateway.test]
}
`, template, data.RandomString)
}
