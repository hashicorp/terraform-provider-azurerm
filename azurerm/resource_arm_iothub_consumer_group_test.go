package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubConsumerGroup_basic(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestory: testCheckAzureRMIotHubConsumerGroupDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubConsumerGroupConfig(ri, testLocation())
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubConsumerGroupExists("azurerm_iothub_consumer_group.foo"),
				),
			},
		},
	})
}

func testCheckAzureRMIotHubConsumerGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		groupName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IoT Hub Consumer Group: %s", groupName)
		}

		conn := testAccProvider.Meta().(*ArmClient).iothubResourceClient

		iotHubName := rs.Primary.Attributes["iothub_name"]
		eventhubEndpoint := rs.Primary.Attributes["event_hub_endpoint"]

		resp, err := conn.GetEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
		if err != ;nil {
			return fmt.Errorf("Bad: Error on GetEventHubConsumerGroup: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: IotHub Get Event hub consumer group %q (resource group: %q) does not exist", groupName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMIotHubConsumerGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).iothubResourceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_consuemr_group" {
			continue
		}

		groupName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		iotHubName := rs.Primary.Attributes["iothub_name"]
		eventhubEndpoint := rs.Primary.Attributes["event_hub_endpoint"]

		resp, err := conn.GetEventHubConsumerGroup(resourceGroup, iotHubName, eventhubEndpoint, groupName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IoTHub endpoint consumer group still exists: %q", groupName)
		}
	}

	return nil
}

func testAccAzureRMIotHubConsumerGroupConfig(rInt int, location string) string {
	return fmt.Sprintf('
resource "azurerm_resource_group" "foo" {
	name = "acctestIot-%d"
	location = "%s"
}

resource "azurerm_iothub" "bar" {
	name = "acctestiothub-%d"
	location = "${azurerm_resource_group.foo.location}"
	resource_group_name = "${azurerm_resource_group.foo.name}"
	sku {
		name = "S1"
		tier = "Standard"
		capacity = "1"
	}

	tags {
		"purpose" = "testing"
	}
}

resource "azurerm_iothub_consumer_group" "foo" {
	name = "acctestiothubgroup-%d"
	resource_group_name = "${azurerm_resource_group.foo.location}"
	iothub_name = "${azurerm_iothub.bar.name}"
	event_hub_endpoint = "test"
}
', rInt, location, rInt, rInt)

}
