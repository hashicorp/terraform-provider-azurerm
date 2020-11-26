package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceBusSubscriptionRule_requiresImport),
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_sqlFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sql_filter", "2=2"),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilterUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sql_filter", "3=3"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_correlationFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "correlation_filter.0.message_id", "test_message_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "correlation_filter.0.reply_to", ""),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "correlation_filter.0.message_id", "test_message_id_updated"),
					resource.TestCheckResourceAttr(data.ResourceName, "correlation_filter.0.reply_to", "test_reply_to_added"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_updateSqlFilterToCorrelationFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ServiceBus.SubscriptionRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ruleName := rs.Primary.Attributes["name"]
		subscriptionName := rs.Primary.Attributes["subscription_name"]
		topicName := rs.Primary.Attributes["topic_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Subscription Rule: %q", ruleName)
		}

		resp, err := client.Get(ctx, resourceGroup, namespaceName, topicName, subscriptionName, ruleName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusRulesClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Subscription Rule %q (resource group: %q) does not exist", ruleName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "SqlFilter"
  sql_filter          = "2=2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "import" {
  name                = azurerm_servicebus_subscription_rule.test.name
  namespace_name      = azurerm_servicebus_subscription_rule.test.namespace_name
  topic_name          = azurerm_servicebus_subscription_rule.test.topic_name
  subscription_name   = azurerm_servicebus_subscription_rule.test.subscription_name
  resource_group_name = azurerm_servicebus_subscription_rule.test.resource_group_name
  filter_type         = azurerm_servicebus_subscription_rule.test.filter_type
  sql_filter          = azurerm_servicebus_subscription_rule.test.sql_filter
}
`, testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(data))
}

func testAccAzureRMServiceBusSubscriptionRule_basicSqlFilterUpdated(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "SqlFilter"
  sql_filter          = "3=3"
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "SqlFilter"
  sql_filter          = "2=2"
  action              = "SET Test='true'"
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id      = "test_correlation_id"
    message_id          = "test_message_id"
    to                  = "test_to"
    reply_to            = "test_reply_to"
    label               = "test_label"
    session_id          = "test_session_id"
    reply_to_session_id = "test_reply_to_session_id"
    content_type        = "test_content_type"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilter(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id"
    properties = {
      test_key = "test_value"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id_updated"
    reply_to       = "test_reply_to_added"
    properties = {
      test_key = "test_value_updated"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  subscription_name   = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  action              = "SET Test='true'"
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceBusSubscriptionRule_template(data acceptance.TestData) string {
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

resource "azurerm_servicebus_subscription" "test" {
  name                = "acctestservicebussubscription-%d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
  resource_group_name = azurerm_resource_group.test.name
  max_delivery_count  = 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
