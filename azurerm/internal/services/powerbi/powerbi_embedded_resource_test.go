package powerbi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPowerBIEmbedded_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbedded_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPowerBIEmbedded_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbedded_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMPowerBIEmbedded_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_powerbi_embedded"),
			},
		},
	})
}

func TestAccAzureRMPowerBIEmbedded_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbedded_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPowerBIEmbedded_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_powerbi_embedded", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPowerBIEmbeddedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPowerBIEmbedded_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPowerBIEmbedded_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPowerBIEmbedded_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPowerBIEmbeddedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "A1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPowerBIEmbeddedExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("PowerBI Embedded not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).PowerBI.CapacityClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.GetDetails(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PowerBI Embedded (PowerBI Embedded Name %q / Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on PowerBI.CapacityClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPowerBIEmbeddedDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).PowerBI.CapacityClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_powerbi_embedded" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.GetDetails(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on CapacityClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPowerBIEmbedded_basic(data acceptance.TestData) string {
	template := testAccAzureRMPowerBIEmbedded_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_embedded" "test" {
  name                = "acctestpowerbi%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A1"
  administrators      = ["${data.azurerm_client_config.test.object_id}"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMPowerBIEmbedded_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_embedded" "import" {
  name                = "${azurerm_powerbi_embedded.test.name}"
  location            = "${azurerm_powerbi_embedded.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A1"
  administrators      = ["${data.azurerm_client_config.test.object_id}"]
}
`, testAccAzureRMPowerBIEmbedded_basic(data))
}

func testAccAzureRMPowerBIEmbedded_complete(data acceptance.TestData) string {
	template := testAccAzureRMPowerBIEmbedded_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_powerbi_embedded" "test" {
  name                = "acctestpowerbi%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "A2"
  administrators      = ["${data.azurerm_client_config.test.object_id}"]

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMPowerBIEmbedded_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-powerbi-%d"
  location = "%s"
}

data "azurerm_client_config" "test" {}
`, data.RandomInteger, data.Locations.Primary)
}
