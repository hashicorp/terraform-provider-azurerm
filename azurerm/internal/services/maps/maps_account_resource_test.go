package maps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MapsAccountResource struct {
}

func TestAccMapsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("x_ms_client_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("S0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMapsAccount_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sku(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("x_ms_client_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("S1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMapsAccount_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("testing"),
			),
		},
		data.ImportStep(),
	})
}

func (MapsAccountResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maps.AccountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Maps Account %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (MapsAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MapsAccountResource) sku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MapsAccountResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S0"

  tags = {
    environment = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
