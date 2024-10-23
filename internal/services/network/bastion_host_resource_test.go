// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-01-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BastionHostResource struct{}

func TestAccBastionHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccBastionHost_standardSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccBastionHost_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBastionHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_bastion_host"),
		},
	})
}

func TestAccBastionHost_scaleUnits(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleUnits(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.scaleUnits(data, 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccBastionHost_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.sku(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBastionHost_developerSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.developerSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBastionHost_premiumSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")
	r := BastionHostResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (BastionHostResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := bastionhosts.ParseBastionHostID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.BastionHostsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Bastion Host (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (BastionHostResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}

func (BastionHostResource) standardSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                   = "acctestBastion%s"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  sku                    = "Standard"
  file_copy_enabled      = true
  ip_connect_enabled     = true
  kerberos_enabled       = true
  shareable_link_enabled = true
  tunneling_enabled      = true

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}

func (BastionHostResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  copy_paste_enabled  = false

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}

func (r BastionHostResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_bastion_host" "import" {
  name                = azurerm_bastion_host.test.name
  resource_group_name = azurerm_bastion_host.test.resource_group_name
  location            = azurerm_bastion_host.test.location

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, r.basic(data))
}

func (BastionHostResource) scaleUnits(data acceptance.TestData, scaleUnits int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  scale_units         = %d

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString, scaleUnits)
}

func (BastionHostResource) sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "%s"

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString, sku)
}

func (BastionHostResource) developerSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Developer"
  virtual_network_id  = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Ternary, data.RandomString, data.RandomString)
}

func (BastionHostResource) premiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                      = "acctestBastion%s"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "Premium"
  file_copy_enabled         = true
  ip_connect_enabled        = true
  kerberos_enabled          = true
  shareable_link_enabled    = true
  tunneling_enabled         = true
  session_recording_enabled = true

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}
