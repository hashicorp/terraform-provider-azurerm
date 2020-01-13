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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMServiceBusQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "false"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceBusQueue_requiresImport),
		},
	})
}

func TestAccAzureRMServiceBusQueue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_update(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_in_megabytes", "2048"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_enablePartitioningStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enablePartitioningStandard(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "true"),
					// Ensure size is read back in it's original value and not the x16 value returned by Azure
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_in_megabytes", "5120"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusQueue_defaultEnablePartitioningPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_Premium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableDuplicateDetection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "requires_duplicate_detection", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enableDuplicateDetection(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "requires_duplicate_detection", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableRequiresSession(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "requires_session", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enableRequiresSession(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "requires_session", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "dead_lettering_on_message_expiration", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "dead_lettering_on_message_expiration", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_lockDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_lockDuration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "lock_duration", "PT40S"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_lockDurationUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "lock_duration", "PT2M"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_isoTimeSpanAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_isoTimeSpanAttributes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_delete_on_idle", "PT10M"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_message_ttl", "PT30M"),
					resource.TestCheckResourceAttr(data.ResourceName, "requires_duplicate_detection", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "duplicate_detection_history_time_window", "PT15M"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusQueue_maxDeliveryCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "max_delivery_count", "10"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueue_maxDeliveryCount(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "max_delivery_count", "20"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMServiceBusQueueDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.QueuesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.QueuesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMServiceBusQueue_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusQueue_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue" "import" {
  name                = "${azurerm_servicebus_queue.test.name}"
  resource_group_name = "${azurerm_servicebus_queue.test.resource_group_name}"
  namespace_name      = "${azurerm_servicebus_queue.test.namespace_name}"
}
`, template)
}

func testAccAzureRMServiceBusQueue_Premium(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_update(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_enablePartitioningStandard(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_enableDuplicateDetection(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_enableRequiresSession(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_enableDeadLetteringOnMessageExpiration(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_lockDuration(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_lockDurationUpdated(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_isoTimeSpanAttributes(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusQueue_maxDeliveryCount(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
