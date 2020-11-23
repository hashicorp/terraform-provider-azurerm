package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMServiceBusTopic_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_service_fabric_cluster"),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopic_basicDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basicDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusTopic_basicDisableEnable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_basicDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_update(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_batched_operations", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_express", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusTopic_enablePartitioningStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_enablePartitioningStandard(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "true"),
					// Ensure size is read back in its original value and not the x16 value returned by Azure
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_in_megabytes", "5120"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusTopic_enablePartitioningPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basicPremium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_enablePartitioningPremium(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_partitioning", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_in_megabytes", "81920"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusTopic_enableDuplicateDetection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusTopic_enableDuplicateDetection(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "requires_duplicate_detection", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusTopic_isoTimeSpanAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusTopic_isoTimeSpanAttributes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusTopicExists(data.ResourceName),
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

func testCheckAzureRMServiceBusTopicDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.TopicsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_topic" {
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
			return fmt.Errorf("ServiceBus Topic still exists:\n%+v", resp.SBTopicProperties)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusTopicExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.TopicsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		topicName := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for topic: %s", topicName)
		}

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusTopicsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Topic %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusTopic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_topic" "import" {
  name                = azurerm_servicebus_topic.test.name
  namespace_name      = azurerm_servicebus_topic.test.namespace_name
  resource_group_name = azurerm_servicebus_topic.test.resource_group_name
}
`, testAccAzureRMServiceBusTopic_basic(data))
}

func testAccAzureRMServiceBusTopic_basicDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  status              = "disabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                      = "acctestservicebustopic-%d"
  namespace_name            = azurerm_servicebus_namespace.test.name
  resource_group_name       = azurerm_resource_group.test.name
  enable_batched_operations = true
  enable_express            = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_basicPremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  enable_partitioning = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_enablePartitioningStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = azurerm_servicebus_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  enable_partitioning   = true
  max_size_in_megabytes = 5120
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_enablePartitioningPremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "premium"
  capacity            = 1
}

resource "azurerm_servicebus_topic" "test" {
  name                  = "acctestservicebustopic-%d"
  namespace_name        = azurerm_servicebus_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  enable_partitioning   = false
  max_size_in_megabytes = 81920
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_enableDuplicateDetection(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                         = "acctestservicebustopic-%d"
  namespace_name               = azurerm_servicebus_namespace.test.name
  resource_group_name          = azurerm_resource_group.test.name
  requires_duplicate_detection = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMServiceBusTopic_isoTimeSpanAttributes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                                    = "acctestservicebustopic-%d"
  namespace_name                          = azurerm_servicebus_namespace.test.name
  resource_group_name                     = azurerm_resource_group.test.name
  auto_delete_on_idle                     = "PT10M"
  default_message_ttl                     = "PT30M"
  requires_duplicate_detection            = true
  duplicate_detection_history_time_window = "PT15M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
