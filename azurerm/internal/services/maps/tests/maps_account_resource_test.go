package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMapsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "S0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMapsAccount_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_sku(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "S1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMapsAccount_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMapsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMapsAccount_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMapsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "testing"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMapsAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Maps.AccountsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on MapsAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Maps Account %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMapsAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Maps.AccountsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maps_account" {
			continue
		}

		id, err := parse.AccountID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Error retrieving Maps Account %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func testAccAzureRMMapsAccount_basic(data acceptance.TestData) string {
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

func testAccAzureRMMapsAccount_sku(data acceptance.TestData) string {
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

func testAccAzureRMMapsAccount_tags(data acceptance.TestData) string {
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
