package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotHubConsumerGroup_events(t *testing.T) {
	resourceName := "azurerm_iothub_consumer_group.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(rInt, location, "events"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "eventhub_endpoint_name", "events"),
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

func TestAccAzureRMIotHubConsumerGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iothub_consumer_group.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(rInt, location, "events"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "eventhub_endpoint_name", "events"),
				),
			},
			{
				Config:      testAccAzureRMIotHubConsumerGroup_requiresImport(rInt, location, "events"),
				ExpectError: testRequiresImportError("azurerm_iothub_consumer_group"),
			},
		},
	})
}

func TestAccAzureRMIotHubConsumerGroup_operationsMonitoringEvents(t *testing.T) {
	resourceName := "azurerm_iothub_consumer_group.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroup_basic(rInt, location, "operationsMonitoringEvents"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "eventhub_endpoint_name", "operationsMonitoringEvents"),
				),
			}, {
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMIotHubConsumerGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).iothub.ResourceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		iotHubName := rs.Primary.Attributes["iothub_name"]
		endpointName := rs.Primary.Attributes["eventhub_endpoint_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).iothub.ResourceClient
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

func testAccAzureRMIotHubConsumerGroup_basic(rInt int, location, eventName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
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

resource "azurerm_iothub_consumer_group" "test" {
  name                   = "test"
  iothub_name            = "${azurerm_iothub.test.name}"
  eventhub_endpoint_name = "%s"
  resource_group_name    = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, eventName)
}

func testAccAzureRMIotHubConsumerGroup_requiresImport(rInt int, location, eventName string) string {
	template := testAccAzureRMIotHubConsumerGroup_basic(rInt, location, eventName)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_consumer_group" "import" {
  name                   = "${azurerm_iothub_consumer_group.test.name}"
  iothub_name            = "${azurerm_iothub_consumer_group.test.iothub_name}"
  eventhub_endpoint_name = "${azurerm_iothub_consumer_group.test.eventhub_endpoint_name}"
  resource_group_name    = "${azurerm_iothub_consumer_group.test.resource_group_name}"
}
`, template)
}
