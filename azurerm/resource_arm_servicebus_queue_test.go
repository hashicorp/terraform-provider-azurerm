package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusQueue_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusQueue_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
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
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_update(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
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
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_enablePartitioningStandard(ri, location)

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
	config := testAccAzureRMServiceBusQueue_Premium(ri, testLocation())

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
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_enableDuplicateDetection(ri, location)

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

func TestAccAzureRMServiceBusQueue_enableRequiresSession(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	location := testLocation()
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "requires_session", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enableRequiresSession(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requires_session", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	location := testLocation()
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dead_lettering_on_message_expiration", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dead_lettering_on_message_expiration", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_lockDuration(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := acctest.RandInt()
	location := testLocation()

	config := testAccAzureRMServiceBusQueue_lockDuration(ri, location)
	updatedConfig := testAccAzureRMServiceBusQueue_lockDurationUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "lock_duration", "PT40S"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "lock_duration", "PT2M"),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceBusQueuesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_queue" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Queue still exists:\n%#v", resp.SBQueueProperties)
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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, namespaceName, queueName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusQueuesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Queue %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusQueue_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_Premium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "premium"
    capacity = 1
}

resource "azurerm_servicebus_queue" "test" {
    name = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_partitioning = true
    enable_express = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_express = true
    max_size_in_megabytes = 2048
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enablePartitioningStandard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    enable_partitioning = true
    max_size_in_megabytes = 5120
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enableDuplicateDetection(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    requires_duplicate_detection = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enableRequiresSession(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
    name                = "acctestservicebusnamespace-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name                = "acctestservicebusqueue-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    namespace_name      = "${azurerm_servicebus_namespace.test.name}"
    requires_session    = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
    name                = "acctestservicebusnamespace-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    sku = "standard"
}

resource "azurerm_servicebus_queue" "test" {
    name                				 = "acctestservicebusqueue-%d"
    resource_group_name 				 = "${azurerm_resource_group.test.name}"
    namespace_name      				 = "${azurerm_servicebus_namespace.test.name}"
    dead_lettering_on_message_expiration = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_lockDuration(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    lock_duration = "PT40S"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_lockDurationUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
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
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    lock_duration = "PT2M"
}
`, rInt, location, rInt, rInt)
}
