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

func TestAccAzureRMDigitalTwinsInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwinsInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDigitalTwinsInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwinsInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDigitalTwinsInstance_requiresImport),
		},
	})
}

func TestAccAzureRMDigitalTwinsInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwinsInstance_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDigitalTwinsInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDigitalTwinsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDigitalTwinsInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwinsInstance_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwinsInstance_updateTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDigitalTwinsInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDigitalTwinsInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "host_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDigitalTwinsInstanceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DigitalTwins.InstanceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Digital Twins Instance not found: %s", resourceName)
		}
		id, err := parse.DigitalTwinsInstanceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Digital Twins Instance %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Digital Twins Instance Client: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMDigitalTwinsInstanceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DigitalTwins.InstanceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_digital_twins_instance" {
			continue
		}
		id, err := parse.DigitalTwinsInstanceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on DigitalTwins.InstanceClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMDigitalTwinsInstance_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-DigitalTwins-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMDigitalTwinsInstance_basic(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwinsInstance_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMDigitalTwinsInstance_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDigitalTwinsInstance_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "import" {
  name                = azurerm_digital_twins_instance.test.name
  resource_group_name = azurerm_digital_twins_instance.test.resource_group_name
  location            = azurerm_digital_twins_instance.test.location
}
`, config)
}

func testAccAzureRMDigitalTwinsInstance_complete(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwinsInstance_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDigitalTwinsInstance_updateTags(data acceptance.TestData) string {
	template := testAccAzureRMDigitalTwinsInstance_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}
