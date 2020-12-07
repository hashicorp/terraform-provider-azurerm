package iothub_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMIotHubConsumerGroup_events(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(data, "events"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub_endpoint_name", "events"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubConsumerGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(data, "events"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub_endpoint_name", "events"),
				),
			},
			{
				Config:      testAccAzureRMIotHubConsumerGroup_requiresImport(data, "events"),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_consumer_group"),
			},
		},
	})
}

func TestAccAzureRMIotHubConsumerGroup_operationsMonitoringEvents(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(data, "operationsMonitoringEvents"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub_endpoint_name", "operationsMonitoringEvents"),
				),
			}, data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubConsumerGroup_withSharedAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_withSharedAccessPolicy(data, "events"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIotHubConsumerGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_consumer_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		iotHubName := rs.Primary.Attributes["iothub_name"]
		endpointName := rs.Primary.Attributes["eventhub_endpoint_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Consumer Group %q still exists in Endpoint %q / IotHub %q / Resource Group %q", name, endpointName, iotHubName, resourceGroup)
		}
	}
	return nil
}

func testCheckAzureRMIotHubConsumerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		iotHubName := rs.Primary.Attributes["iothub_name"]
		endpointName := rs.Primary.Attributes["eventhub_endpoint_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetEventHubConsumerGroup(ctx, resourceGroup, iotHubName, endpointName, name)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Consumer Group %q (Endpoint %q / IotHub %q / Resource Group: %q) does not exist", name, endpointName, iotHubName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMIotHubConsumerGroup_basic(data acceptance.TestData, eventName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_consumer_group" "test" {
  name                   = "test"
  iothub_name            = azurerm_iothub.test.name
  eventhub_endpoint_name = "%s"
  resource_group_name    = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, eventName)
}

func testAccAzureRMIotHubConsumerGroup_requiresImport(data acceptance.TestData, eventName string) string {
	template := testAccAzureRMIotHubConsumerGroup_basic(data, eventName)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_consumer_group" "import" {
  name                   = azurerm_iothub_consumer_group.test.name
  iothub_name            = azurerm_iothub_consumer_group.test.iothub_name
  eventhub_endpoint_name = azurerm_iothub_consumer_group.test.eventhub_endpoint_name
  resource_group_name    = azurerm_iothub_consumer_group.test.resource_group_name
}
`, template)
}

func testAccAzureRMIotHubConsumerGroup_withSharedAccessPolicy(data acceptance.TestData, eventName string) string {
	template := testAccAzureRMIotHubConsumerGroup_basic(data, eventName)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_shared_access_policy" "test" {
  name                = "acctestSharedAccessPolicy"
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  service_connect     = true
}
`, template)
}
