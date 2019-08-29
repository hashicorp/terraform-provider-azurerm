package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
		},
	})
}
func TestAccAzureRMServiceBusSubscriptionRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMServiceBusSubscriptionRule_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_servicebus_subscription_rule"),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_sqlFilterUpdated(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(ri, location)
	updatedConfig := testAccAzureRMServiceBusSubscriptionRule_basicSqlFilterUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sql_filter", "2=2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sql_filter", "3=3"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMServiceBusSubscriptionRule_correlationFilter(ri, location)
	updatedConfig := testAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "correlation_filter.0.message_id", "test_message_id"),
					resource.TestCheckResourceAttr(resourceName, "correlation_filter.0.reply_to", ""),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "correlation_filter.0.message_id", "test_message_id_updated"),
					resource.TestCheckResourceAttr(resourceName, "correlation_filter.0.reply_to", "test_reply_to_added"),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusSubscriptionRule_updateSqlFilterToCorrelationFilter(t *testing.T) {
	resourceName := "azurerm_servicebus_subscription_rule.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(ri, location)
	updatedConfig := testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMServiceBusSubscriptionRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		client := testAccProvider.Meta().(*ArmClient).servicebus.SubscriptionRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  filter_type         = "SqlFilter"
  sql_filter          = "2=2"
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "import" {
  name                = "${azurerm_servicebus_subscription_rule.test.name}"
  namespace_name      = "${azurerm_servicebus_subscription_rule.test.namespace_name}"
  topic_name          = "${azurerm_servicebus_subscription_rule.test.topic_name}"
  subscription_name   = "${azurerm_servicebus_subscription_rule.test.subscription_name}"
  resource_group_name = "${azurerm_servicebus_subscription_rule.test.resource_group_name}"
  filter_type         = "${azurerm_servicebus_subscription_rule.test.filter_type}"
  sql_filter          = "${azurerm_servicebus_subscription_rule.test.sql_filter}"
}
`, testAccAzureRMServiceBusSubscriptionRule_basicSqlFilter(rInt, location))
}

func testAccAzureRMServiceBusSubscriptionRule_basicSqlFilterUpdated(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  filter_type         = "SqlFilter"
  sql_filter          = "3=3"
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_sqlFilterWithAction(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  filter_type         = "SqlFilter"
  sql_filter          = "2=2"
  action              = "SET Test='true'"
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_basicCorrelationFilter(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilter(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id"
  }
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilterUpdated(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id_updated"
    reply_to       = "test_reply_to_added"
  }
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_correlationFilterWithAction(rInt int, location string) string {
	template := testAccAzureRMServiceBusSubscriptionRule_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name                = "acctestservicebusrule-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  subscription_name   = "${azurerm_servicebus_subscription.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  action              = "SET Test='true'"
  filter_type         = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id"
  }
}
`, template, rInt)
}

func testAccAzureRMServiceBusSubscriptionRule_template(rInt int, location string) string {
	return fmt.Sprintf(`
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
}
`, rInt, location, rInt, rInt, rInt)
}
