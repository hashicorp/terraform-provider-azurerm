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

func TestAccAzureRMServiceBusSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceBusSubscription_requiresImport),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_defaultTtl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_withDefaultTtl(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr("azurerm_servicebus_subscription.test", "default_message_ttl", "PT1H"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateEnableBatched(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_updateEnableBatched(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "enable_batched_operations", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateRequiresSession(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_updateRequiresSession(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "requires_session", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateForwardTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_to-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_updateForwardTo(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "forward_to", expectedValue),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_dl_messages_to-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "forward_dead_lettered_messages_to", expectedValue),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_updateDeadLetteringOnFilterEvaluationExceptions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_updateDeadLetteringOnFilterEvaluationExceptions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceBusSubscription_status(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Active"),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_status(data, "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Disabled"),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_status(data, "ReceiveDisabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "ReceiveDisabled"),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscription_status(data, "Active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Active"),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusSubscriptionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.SubscriptionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_subscription" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Subscription still exists:\n%+v", resp.SBSubscriptionProperties)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusSubscriptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.SubscriptionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subscriptionName := rs.Primary.Attributes["name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Subscription: %q", topicName)
		}

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusSubscriptionsClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Subscription %q (resource group: %q) does not exist", subscriptionName, resourceGroup)
		}

		return nil
	}
}

const testAccAzureRMServiceBusSubscription_tfTemplate = `
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
    name                = "acctestservicebusnamespace-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
    name                = "acctestservicebustopic-%d"
    namespace_name      = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_subscription" "test" {
    name                = "acctestservicebussubscription-%d"
    namespace_name      = "${azurerm_servicebus_namespace.test.name}"
    topic_name          = "${azurerm_servicebus_topic.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    max_delivery_count  = 10
	%s
}
`

func testAccAzureRMServiceBusSubscription_basic(data acceptance.TestData) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, "")
}

func testAccAzureRMServiceBusSubscription_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription" "import" {
  name                = azurerm_servicebus_subscription.test.name
  namespace_name      = azurerm_servicebus_subscription.test.namespace_name
  topic_name          = azurerm_servicebus_subscription.test.topic_name
  resource_group_name = azurerm_servicebus_subscription.test.resource_group_name
  max_delivery_count  = azurerm_servicebus_subscription.test.max_delivery_count
}
`, testAccAzureRMServiceBusSubscription_basic(data))
}

func testAccAzureRMServiceBusSubscription_withDefaultTtl(data acceptance.TestData) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"default_message_ttl = \"PT1H\"\n")
}

func testAccAzureRMServiceBusSubscription_updateEnableBatched(data acceptance.TestData) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"enable_batched_operations = true\n")
}

func testAccAzureRMServiceBusSubscription_updateRequiresSession(data acceptance.TestData) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"requires_session = true\n")
}

func testAccAzureRMServiceBusSubscription_updateForwardTo(data acceptance.TestData) string {
	forwardToTf := testAccAzureRMServiceBusSubscription_tfTemplate + `

resource "azurerm_servicebus_topic" "forward_to" {
    name = "acctestservicebustopic-forward_to-%d"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

`
	return fmt.Sprintf(forwardToTf, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"forward_to = \"${azurerm_servicebus_topic.forward_to.name}\"\n", data.RandomInteger)
}

func testAccAzureRMServiceBusSubscription_updateForwardDeadLetteredMessagesTo(data acceptance.TestData) string {
	forwardToTf := testAccAzureRMServiceBusSubscription_tfTemplate + `

resource "azurerm_servicebus_topic" "forward_dl_messages_to" {
    name = "acctestservicebustopic-forward_dl_messages_to-%d"
    namespace_name = "${azurerm_servicebus_namespace.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

`
	return fmt.Sprintf(forwardToTf, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"forward_dead_lettered_messages_to = \"${azurerm_servicebus_topic.forward_dl_messages_to.name}\"\n", data.RandomInteger)
}

func testAccAzureRMServiceBusSubscription_status(data acceptance.TestData, status string) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		fmt.Sprintf("status = \"%s\"", status))
}

func testAccAzureRMServiceBusSubscription_updateDeadLetteringOnFilterEvaluationExceptions(data acceptance.TestData) string {
	return fmt.Sprintf(testAccAzureRMServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"dead_lettering_on_filter_evaluation_error = false\n")
}
