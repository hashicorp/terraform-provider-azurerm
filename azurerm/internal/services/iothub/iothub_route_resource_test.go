package iothub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubRouteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubRoute_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubRouteExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubRoute_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_route"),
			},
		},
	})
}

func TestAccAzureRMIotHubRoute_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubRouteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIotHubRoute_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubRouteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIotHubRouteDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_route" {
			continue
		}

		routeName := rs.Primary.Attributes["name"]
		iothubName := rs.Primary.Attributes["iothub_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}
		if iothub.Properties == nil || iothub.Properties.Routing == nil {
			return nil
		}
		routes := iothub.Properties.Routing.Routes

		if routes == nil {
			return nil
		}

		for _, route := range *routes {
			if strings.EqualFold(*route.Name, routeName) {
				return fmt.Errorf("Bad: route %s still exists on IoTHb %s", routeName, iothubName)
			}
		}
	}
	return nil
}

func testCheckAzureRMIotHubRouteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		parsedIothubId, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		iothubName := parsedIothubId.Path["IotHubs"]
		routeName := parsedIothubId.Path["Routes"]
		resourceGroup := parsedIothubId.ResourceGroup

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
			}

			return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}

		if iothub.Properties == nil || iothub.Properties.Routing == nil {
			return fmt.Errorf("Bad: No route %s defined for IotHub %s", routeName, iothubName)
		}
		routes := iothub.Properties.Routing.Routes

		if routes == nil {
			return fmt.Errorf("Bad: No route %s defined for IotHub %s", routeName, iothubName)
		}

		for _, route := range *routes {
			if strings.EqualFold(*route.Name, routeName) {
				return nil
			}
		}

		return fmt.Errorf("Bad: No route %s defined for IotHub %s", routeName, iothubName)
	}
}

func testAccAzureRMIotHubRoute_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHubRoute_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_route" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
  enabled        = true
}
`, template)
}

func testAccAzureRMIotHubRoute_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string          = azurerm_storage_account.test.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.test.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_route" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
  enabled        = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMIotHubRoute_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string          = azurerm_storage_account.test.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.test.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_route" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  source         = "DeviceLifecycleEvents"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
  enabled        = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
