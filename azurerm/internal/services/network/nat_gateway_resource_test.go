package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NatGatewayResource struct {
}

func TestAccNatGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_ip_address_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("public_ip_prefix_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("10"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_ip_address_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("public_ip_prefix_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("10"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t NatGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NatGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NatGatewayClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading NAT Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

// Using alt location because the resource currently in private preview and is only available in eastus2.
func (NatGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctestnatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

// Using alt location because the resource currently in private preview and is only available in eastus2.
func (NatGatewayResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicIPPrefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "test" {
  name                    = "acctestnatGateway-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  public_ip_address_ids   = [azurerm_public_ip.test.id]
  public_ip_prefix_ids    = [azurerm_public_ip_prefix.test.id]
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
