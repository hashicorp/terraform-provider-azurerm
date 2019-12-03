package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubEndpointServiceBusQueue_basic(t *testing.T) {
	resourceName := "azurerm_iothub_endpoint_servicebus_queue.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAzureRMIotHubEndpointStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEndpointServiceBusQueue_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMIotHubEndpointServiceBusQueueExists(resourceName),
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

func TestAccAzureRMIotHubEndpointServiceBusQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_iothub_endpoint_servicebus_queue.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAzureRMIotHubEndpointServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEndpointServiceBusQueue_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMIotHubEndpointServiceBusQueueExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubEndpointServiceBusQueue_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_iothub_endpoint_servicebus_queue"),
			},
		},
	})
}

func testAccAzureRMIotHubEndpointServiceBusQueue_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}


resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  queue_name          = "${azurerm_servicebus_queue.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_servicebus_queue" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"
  
  connection_string = "${azurerm_servicebus_queue_authorization_rule.test.primary_connection_string}"
}
`, rInt, location)
}

func testAccAzureRMIotHubEndpointServiceBusQueue_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHubEndpointServiceBusQueue_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_queue" "import" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"
    
  connection_string = "${azurerm_servicebus_queue_authorization_rule.test.primary_connection_string}"
}
`, template)
}

func testAccAzureRMIotHubEndpointServiceBusQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		client := testAccProvider.Meta().(*ArmClient).IoTHub.ResourceClient

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
		endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues

		if endpoints == nil {
			return fmt.Errorf("Bad: No ServiceBus Queue endpoint %s defined for IotHub %s", endpointName, iothubName)
		}

		for _, endpoint := range *endpoints {
			if strings.EqualFold(*endpoint.Name, endpointName) {
				return nil
			}
		}

		return fmt.Errorf("Bad: No ServiceBus Queue endpoint %s defined for IotHub %s", endpointName, iothubName)
	}
}

func testAccAzureRMIotHubEndpointServiceBusQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).IoTHub.ResourceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_endpoint_servicebus_queue" {
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
		endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues

		if endpoints == nil {
			return nil
		}

		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, endpointName) {
					return fmt.Errorf("Bad: ServiceBus Queue endpoint %s still exists on IoTHb %s", endpointName, iothubName)
				}
			}
		}
	}
	return nil
}
