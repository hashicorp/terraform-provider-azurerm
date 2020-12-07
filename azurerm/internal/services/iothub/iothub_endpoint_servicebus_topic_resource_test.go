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

func TestAccAzureRMIotHubEndpointServiceBusTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMIotHubEndpointStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEndpointServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMIotHubEndpointServiceBusTopicExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubEndpointServiceBusTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMIotHubEndpointServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEndpointServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMIotHubEndpointServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubEndpointServiceBusTopic_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_endpoint_servicebus_topic"),
			},
		},
	})
}

func testAccAzureRMIotHubEndpointServiceBusTopic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
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

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMIotHubEndpointServiceBusTopic_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHubEndpointServiceBusTopic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_topic" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, template)
}

func testAccAzureRMIotHubEndpointServiceBusTopicExists(resourceName string) resource.TestCheckFunc {
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
		endpointName := parsedIothubId.Path["Endpoints"]
		resourceGroup := parsedIothubId.ResourceGroup

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
			}

			return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}

		if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
			return fmt.Errorf("Bad: No endpoint %s defined for IotHub %s", endpointName, iothubName)
		}
		endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics

		if endpoints == nil {
			return fmt.Errorf("Bad: No ServiceBus Topic endpoint %s defined for IotHub %s", endpointName, iothubName)
		}

		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, endpointName) {
					return nil
				}
			}
		}
		return fmt.Errorf("Bad: No ServiceBus Topic endpoint %s defined for IotHub %s", endpointName, iothubName)
	}
}

func testAccAzureRMIotHubEndpointServiceBusTopicDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_endpoint_servicebus_topic" {
			continue
		}

		endpointName := rs.Primary.Attributes["name"]
		iothubName := rs.Primary.Attributes["iothub_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}
		if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
			return nil
		}
		endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics

		if endpoints == nil {
			return nil
		}

		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, endpointName) {
					return fmt.Errorf("Bad: ServiceBus Topic endpoint %s still exists on IoTHb %s", endpointName, iothubName)
				}
			}
		}
	}
	return nil
}
