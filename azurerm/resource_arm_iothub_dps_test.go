package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotHubDPS_basic(t *testing.T) {
	resourceName := "azurerm_iothub_dps.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPS_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "device_provisioning_host_name"),
					resource.TestCheckResourceAttrSet(resourceName, "id_scope"),
					resource.TestCheckResourceAttrSet(resourceName, "service_operations_host_name"),
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

func TestAccAzureRMIotHubDPS_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iothub_dps.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPS_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubDPS_requiresImport(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_iothubdps"),
			},
		},
	})
}

func TestAccAzureRMIotHubDPS_update(t *testing.T) {
	resourceName := "azurerm_iothub_dps.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPS_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMIotHubDPS_update(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
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

func TestAccAzureRMIotHubDPS_linkedHubs(t *testing.T) {
	resourceName := "azurerm_iothub_dps.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPS_linkedHubs(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMIotHubDPS_linkedHubsUpdated(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSExists(resourceName),
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

func testCheckAzureRMIotHubDPSDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothubdps" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IoT Device Provisioning Service %s still exists in resource group %s", name, resourceGroup)
		}
	}
	return nil
}

func testCheckAzureRMIotHubDPSExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		iotdpsName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IoT Device Provisioning Service: %s", iotdpsName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient
		resp, err := client.Get(ctx, iotdpsName, resourceGroup)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IoT Device Provisioning Service %q (Resource Group %q) does not exist", iotdpsName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubDPSResourceClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMIotHubDPS_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDPS_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHubDPS_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_dps" "import" {
  name                = "${azurerm_iothub_dps.test.name}"
  resource_group_name = "${azurerm_iothub_dps.test.name}"
  location            = "${azurerm_iothub_dps.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}
`, template)
}

func testAccAzureRMIotHubDPS_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
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
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDPS_linkedHubs(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  linked_hub {
    connection_string       = "HostName=test.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=booo"
    location                = "${azurerm_resource_group.test.location}"
    allocation_weight       = 15
    apply_allocation_policy = true
  }

  linked_hub {
    connection_string = "HostName=test2.azure-devices.net;SharedAccessKeyName=iothubowner2;SharedAccessKey=key2"
    location          = "${azurerm_resource_group.test.location}"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDPS_linkedHubsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  linked_hub {
    connection_string = "HostName=test.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=booo"
    location          = "${azurerm_resource_group.test.location}"
    allocation_weight = 150
  }
}
`, rInt, location, rInt)
}
