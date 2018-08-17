package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMIotHub_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
		},
	})
}

func TestAccAzureRMIotHub_requiresImport(t *testing.T) {
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
			{
				Config:      testAccAzureRMIotHub_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_iothub"),
			},
		},
	})
}

func TestAccAzureRMIotHub_standard(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_standard(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
		},
	})
}

func TestAccAzureRMIotHub_customRoutes(t *testing.T) {
	rInt := acctest.RandInt()
	rStr := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_customRoutes(rInt, rStr, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
		},
	})
}

func testCheckAzureRMIotHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).iothubResourceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IotHub %s still exists in resource group %s", name, resourceGroup)
		}
	}
	return nil
}

func testCheckAzureRMIotHubExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		iothubName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IotHub: %s", iothubName)
		}

		client := testAccProvider.Meta().(*ArmClient).iothubResourceClient
		resp, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IotHub %q (resource group: %q) does not exist", iothubName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		return nil

	}
}

func testAccAzureRMIotHub_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "foo" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.foo.name}"
  location            = "${azurerm_resource_group.foo.location}"
  sku {
    name = "B1"
    tier = "Basic"
    capacity = "1"
  }

  tags {
    "purpose" = "testing"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "import" {
  name                = "${azurerm_iothub.test.name}"
  resource_group_name = "${azurerm_iothub.test.resource_group_name}"
  location            = "${azurerm_iothub.test.location}"
  sku {
    name = "B1"
    tier = "Basic"
    capacity = "1"
  }

  tags {
    "purpose" = "testing"
  }
}
`, template)
}

func testAccAzureRMIotHub_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "foo" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.foo.name}"
  location            = "${azurerm_resource_group.foo.location}"
  sku {
    name = "S1"
    tier = "Standard"
    capacity = "1"
  }

  tags {
    "purpose" = "testing"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHub_customRoutes(rInt int, rStr string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "foo" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                      = "acctestsa%s"
  resource_group_name       = "${azurerm_resource_group.foo.name}"
  location                  = "${azurerm_resource_group.foo.location}"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                      = "test"
  resource_group_name       = "${azurerm_resource_group.foo.name}"
  storage_account_name      = "${azurerm_storage_account.test.name}"
  container_access_type     = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.foo.name}"
  location            = "${azurerm_resource_group.foo.location}"
  sku {
    name = "S1"
    tier = "Standard"
    capacity = "1"
  }

  endpoint {
    type                        = "AzureIotHub.StorageContainer"
    connection_string           = "${azurerm_storage_account.test.primary_blob_connection_string}"
    name                        = "export"
    batch_frequency_in_seconds  = 60
    max_chunk_size_in_bytes     = 10485760
    container_name              = "test"
    encoding                    = "Avro"
    file_name_format            = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  }

  route {
    name            = "export"
    source          = "DeviceMessages"
    condition       = "true"
    endpoint_names  = ["export"]
    enabled      = true
  }

  tags {
    "purpose" = "testing"
  }
}
`, rInt, location, rStr, rInt)
}
