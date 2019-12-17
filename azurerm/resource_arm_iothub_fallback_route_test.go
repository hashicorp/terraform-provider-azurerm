package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubFallbackRoute_basic(t *testing.T) {
	resourceName := "azurerm_iothub_fallback_route.test"
	rInt := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubFallbackRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubFallbackRoute_basic(rInt, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubFallbackRouteExists(resourceName),
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

func TestAccAzureRMIotHubFallbackRoute_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iothub_fallback_route.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubFallbackRoute_basic(rInt, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubFallbackRouteExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubFallbackRoute_requiresImport(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_fallback_route"),
			},
		},
	})
}

func TestAccAzureRMIotHubFallbackRoute_update(t *testing.T) {
	resourceName := "azurerm_iothub_fallback_route.test"
	rInt := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubFallbackRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubFallbackRoute_basic(rInt, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubFallbackRouteExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMIotHubFallbackRoute_update(rInt, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubFallbackRouteExists(resourceName),
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

func testCheckAzureRMIotHubFallbackRouteDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_fallback_route" {
			continue
		}

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
		if iothub.Properties.Routing.FallbackRoute != nil {
			return fmt.Errorf("Bad: fallback route still exists on IoTHb %s", iothubName)
		}
	}
	return nil
}

func testCheckAzureRMIotHubFallbackRouteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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
		resourceGroup := parsedIothubId.ResourceGroup

		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
			}

			return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}

		if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.FallbackRoute == nil {
			return fmt.Errorf("Bad: No fallbackroute defined for IotHub %s", iothubName)
		}

		return nil
	}
}

func testAccAzureRMIotHubFallbackRoute_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_fallback_route" "import" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = ["${azurerm_iothub_endpoint_storage_container.test.name}"]
  enabled        = true
}
`, template)
}

func testAccAzureRMIotHubFallbackRoute_basic(rInt int, rStr string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-%[1]d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"

  connection_string          = "${azurerm_storage_account.test.primary_blob_connection_string}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = "${azurerm_storage_container.test.name}"
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_fallback_route" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"

  condition      = "true"
  endpoint_names = ["${azurerm_iothub_endpoint_storage_container.test.name}"]
  enabled        = true
}
`, rInt, location, rStr)
}

func testAccAzureRMIotHubFallbackRoute_update(rInt int, rStr string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-%[1]d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"

  connection_string          = "${azurerm_storage_account.test.primary_blob_connection_string}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = "${azurerm_storage_container.test.name}"
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_fallback_route" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"

  condition      = "true"
  endpoint_names = ["${azurerm_iothub_endpoint_storage_container.test.name}"]
  enabled        = false
}
`, rInt, location, rStr)
}
