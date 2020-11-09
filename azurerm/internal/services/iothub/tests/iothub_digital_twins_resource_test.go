package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDigitalTwins_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwins_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDigitalTwins_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwins_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDigitalTwins_requiresImport),
		},
	})
}

func TestAccAzureRMDigitalTwins_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwins_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDigitalTwins_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwins_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwins_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwins_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwins_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDigitalTwinsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DigitalTwinsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("digital Twins not found: %s", resourceName)
		}
		id, err := parse.DigitalTwinsID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Digital Twins %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Digitaltwins Client: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMDigitalTwinsDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DigitalTwinsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_digital_twins" {
			continue
		}
		id, err := parse.DigitalTwinsID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Digitaltwins.DigitalTwinsClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMDigitalTwins_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-DigitalTwins-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMDigitalTwins_basic(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwins_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMDigitalTwins_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDigitalTwins_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_digital_twins" "import" {
  name                = azurerm_iothub_digital_twins.test.name
  resource_group_name = azurerm_iothub_digital_twins.test.resource_group_name
  location            = azurerm_iothub_digital_twins.test.location
}
`, config)
}

func testAccAzureRMDigitalTwins_complete(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwins_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDigitalTwins_updateTags(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwins_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}
