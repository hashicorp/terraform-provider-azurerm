package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMdigitaltwinsDigitalTwin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMdigitaltwinsDigitalTwinDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMdigitaltwinsDigitalTwin_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMdigitaltwinsDigitalTwinDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMdigitaltwinsDigitalTwin_requiresImport),
		},
	})
}

func TestAccAzureRMdigitaltwinsDigitalTwin_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMdigitaltwinsDigitalTwinDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMdigitaltwinsDigitalTwin_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMdigitaltwinsDigitalTwinDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMdigitaltwinsDigitalTwin_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMdigitaltwinsDigitalTwinExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMdigitaltwinsDigitalTwinExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Digitaltwins.DigitalTwinClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("digitaltwins DigitalTwin not found: %s", resourceName)
		}
		id, err := parse.DigitaltwinsDigitalTwinID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Digitaltwins DigitalTwin %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Digitaltwins.DigitalTwinClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMdigitaltwinsDigitalTwinDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Digitaltwins.DigitalTwinClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_digital_twins" {
			continue
		}
		id, err := parse.DigitaltwinsDigitalTwinID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Digitaltwins.DigitalTwinClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMdigitaltwinsDigitalTwin_template(data acceptance.TestData) string {
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

func testAccAzureRMdigitaltwinsDigitalTwin_basic(data acceptance.TestData) string {
	template := testAccAzureRMdigitaltwinsDigitalTwin_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMdigitaltwinsDigitalTwin_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMdigitaltwinsDigitalTwin_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins" "import" {
  name                = azurerm_digital_twins.test.name
  resource_group_name = azurerm_digital_twins.test.resource_group_name
  location            = azurerm_digital_twins.test.location
}
`, config)
}

func testAccAzureRMdigitaltwinsDigitalTwin_complete(data acceptance.TestData) string {
	template := testAccAzureRMdigitaltwinsDigitalTwin_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMdigitaltwinsDigitalTwin_updateTags(data acceptance.TestData) string {
	template := testAccAzureRMdigitaltwinsDigitalTwin_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}
