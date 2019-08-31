package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusQueue_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
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
func TestAccAzureRMServiceBusQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
				),
			},
			{
				Config:      testAccAzureRMServiceBusQueue_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_service_fabric_cluster"),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_update(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_update(ri, location)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enablePartitioningStandard(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_enablePartitioningStandard(ri, location)

	resource.ParallelTest(t, resource.TestCase{
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
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusQueue_Premium(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_partitioning", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_express", "false"),
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

func TestAccAzureRMServiceBusQueue_enableDuplicateDetection(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMServiceBusQueue_basic(ri, location)
	postConfig := testAccAzureRMServiceBusQueue_enableDuplicateDetection(ri, location)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableRequiresSession(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	location := testLocation()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	location := testLocation()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_lockDuration(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	config := testAccAzureRMServiceBusQueue_lockDuration(ri, location)
	updatedConfig := testAccAzureRMServiceBusQueue_lockDurationUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_isoTimeSpanAttributes(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusQueue_isoTimeSpanAttributes(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_delete_on_idle", "PT10M"),
					resource.TestCheckResourceAttr(resourceName, "default_message_ttl", "PT30M"),
					resource.TestCheckResourceAttr(resourceName, "requires_duplicate_detection", "true"),
					resource.TestCheckResourceAttr(resourceName, "duplicate_detection_history_time_window", "PT15M"),
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

func TestAccAzureRMServiceBusQueue_maxDeliveryCount(t *testing.T) {
	resourceName := "azurerm_servicebus_queue.test"
	location := testLocation()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "max_delivery_count", "10"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_maxDeliveryCount(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "max_delivery_count", "20"),
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

func testCheckAzureRMServiceBusQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).servicebus.QueuesClient
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

func testCheckAzureRMServiceBusQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		queueName := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for queue: %s", queueName)
		}

		client := testAccProvider.Meta().(*ArmClient).servicebus.QueuesClient
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue" "import" {
  name                = "${azurerm_servicebus_queue.test.name}"
  resource_group_name = "${azurerm_servicebus_queue.test.resource_group_name}"
  namespace_name      = "${azurerm_servicebus_queue.test.namespace_name}"
}
`, testAccAzureRMServiceBusQueue_basic(rInt, location))
}

func testAccAzureRMServiceBusQueue_Premium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "premium"
  capacity            = 1
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  enable_partitioning = false
  enable_express      = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                  = "acctestservicebusqueue-%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  namespace_name        = "${azurerm_servicebus_namespace.test.name}"
  enable_express        = true
  max_size_in_megabytes = 2048
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enablePartitioningStandard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                  = "acctestservicebusqueue-%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  namespace_name        = "${azurerm_servicebus_namespace.test.name}"
  enable_partitioning   = true
  max_size_in_megabytes = 5120
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_enableDuplicateDetection(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                         = "acctestservicebusqueue-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  namespace_name               = "${azurerm_servicebus_namespace.test.name}"
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
  name                                 = "acctestservicebusqueue-%d"
  resource_group_name                  = "${azurerm_resource_group.test.name}"
  namespace_name                       = "${azurerm_servicebus_namespace.test.name}"
  dead_lettering_on_message_expiration = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_lockDuration(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  lock_duration       = "PT40S"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_lockDurationUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  lock_duration       = "PT2M"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_isoTimeSpanAttributes(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                                    = "acctestservicebusqueue-%d"
  resource_group_name                     = "${azurerm_resource_group.test.name}"
  namespace_name                          = "${azurerm_servicebus_namespace.test.name}"
  auto_delete_on_idle                     = "PT10M"
  default_message_ttl                     = "PT30M"
  requires_duplicate_detection            = true
  duplicate_detection_history_time_window = "PT15M"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMServiceBusQueue_maxDeliveryCount(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctestservicebusqueue-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  max_delivery_count  = 20
}
`, rInt, location, rInt, rInt)
}
