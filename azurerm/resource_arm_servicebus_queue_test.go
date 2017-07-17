package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMServiceBusQueue_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusQueue_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_batched_operations", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_update(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri)
	postConfig := testAccAzureRMServiceBusQueue_update(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_batched_operations", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_batched_operations", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_size_in_megabytes", "2048"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enablePartitioningStandard(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri)
	postConfig := testAccAzureRMServiceBusQueue_enablePartitioningStandard(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "true"),
					// Ensure size is read back in it's original value and not the x16 value returned by Azure
					resource.TestCheckResourceAttr(resourceName, "max_size_in_megabytes", "5120"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_defaultEnablePartitioningPremium(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusQueue_Premium(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableDuplicateDetection(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri)
	postConfig := testAccAzureRMServiceBusQueue_enableDuplicateDetection(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "requires_duplicate_detection", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requires_duplicate_detection", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceBusQueuesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_queue" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, namespaceName, name)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Queue still exists:\n%#v", resp.QueueProperties)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusQueueExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		queueName := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for queue: %s", queueName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceBusQueuesClient

		resp, err := client.Get(resourceGroup, namespaceName, queueName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusQueuesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Queue %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusQueue_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
}
`, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_Premium(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "premium"
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_partitioning = true
    enable_express = false
}
`, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_update(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_batched_operations = true
    enable_express = true
    max_size_in_megabytes = 2048
}
`, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enablePartitioningStandard(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_partitioning = true
    max_size_in_megabytes = 5120
}
`, rInt, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enableDuplicateDetection(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    requires_duplicate_detection = true
}
`, rInt, rInt, rInt)
}
