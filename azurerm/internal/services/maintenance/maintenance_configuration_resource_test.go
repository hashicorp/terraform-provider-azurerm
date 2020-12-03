package maintenance_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccMaintenanceConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "All"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccMaintenanceConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccMaintenanceConfiguration_requiresImport),
		},
	})
}

func TestAccMaintenanceConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "Host"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TesT"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccMaintenanceConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaintenanceConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "All"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMaintenanceConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "Host"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TesT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMaintenanceConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMaintenanceConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "All"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckMaintenanceConfigurationDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maintenance_configuration" {
			continue
		}

		id, err := parse.MaintenanceConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckMaintenanceConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.MaintenanceConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on maintenanceConfigurationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Maintenance Configuration %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccMaintenanceConfiguration_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "All"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccMaintenanceConfiguration_requiresImport(data acceptance.TestData) string {
	template := testAccMaintenanceConfiguration_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_configuration" "import" {
  name                = azurerm_maintenance_configuration.test.name
  resource_group_name = azurerm_maintenance_configuration.test.resource_group_name
  location            = azurerm_maintenance_configuration.test.location
  scope               = azurerm_maintenance_configuration.test.scope
}
`, template)
}

func testAccMaintenanceConfiguration_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%d"
  location = "%s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "Host"

  tags = {
    env = "TesT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
