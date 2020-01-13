package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotDPS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_dps", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotDPS_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_iot_dps", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotDPS_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iotdps"),
			},
		},
	})
}

func TestAccAzureRMIotDPS_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_dps", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIotDPS_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotDPS_linkedHubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_dps", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPS_linkedHubs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIotDPS_linkedHubsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIotDPSDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iotdps" {
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

func testCheckAzureRMIotDPSExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		iotdpsName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IoT Device Provisioning Service: %s", iotdpsName)
		}

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

func testAccAzureRMIotDPS_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotDPS_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotDPS_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iot_dps" "import" {
  name                = "${azurerm_iot_dps.test.name}"
  resource_group_name = "${azurerm_iot_dps.test.name}"
  location            = "${azurerm_iot_dps.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}
`, template)
}

func testAccAzureRMIotDPS_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotDPS_linkedHubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotDPS_linkedHubsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
