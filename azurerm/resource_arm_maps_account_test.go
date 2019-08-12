package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMapsAccount_basic(t *testing.T) {
	resourceName := "azurerm_maps_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "s0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMapsAccount_sku(t *testing.T) {
	resourceName := "azurerm_maps_account.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_sku(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "x_ms_client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "s1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMapsAccount_tags(t *testing.T) {
	resourceName := "azurerm_maps_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMapsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMapsAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMapsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMapsAccount_tags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMapsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "testing"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMMapsAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		mapsAccountName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).maps.AccountsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, mapsAccountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on MapsAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Maps Account %q (resource group: %q) does not exist", mapsAccountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMapsAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).maps.AccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maps_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Error retrieving Maps Account %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	return nil

}

func testAccAzureRMMapsAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "s0"
}
`, rInt, location, rInt)
}

func testAccAzureRMMapsAccount_sku(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "s1"
}
`, rInt, location, rInt)
}

func testAccAzureRMMapsAccount_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "s0"

  tags = {
    environment = "testing"
  }
}
`, rInt, location, rInt)
}
