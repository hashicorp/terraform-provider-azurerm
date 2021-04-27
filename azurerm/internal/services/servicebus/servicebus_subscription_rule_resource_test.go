package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusSubscriptionRuleResource struct {
}

func TestAccServiceBusSubscriptionRule_basicSqlFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccServiceBusSubscriptionRule_basicCorrelationFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCorrelationFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_sqlFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sqlFilterWithAction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_correlationFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.correlationFilterWithAction(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_sqlFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sql_filter").HasValue("2=2"),
			),
		},
		{
			Config: r.basicSqlFilterUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sql_filter").HasValue("3=3"),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_correlationFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.correlationFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("correlation_filter.0.message_id").HasValue("test_message_id"),
				check.That(data.ResourceName).Key("correlation_filter.0.reply_to").HasValue(""),
			),
		},
		{
			Config: r.correlationFilterUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("correlation_filter.0.message_id").HasValue("test_message_id_updated"),
				check.That(data.ResourceName).Key("correlation_filter.0.reply_to").HasValue("test_reply_to_added"),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_updateSqlFilterToCorrelationFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicCorrelationFilter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t ServiceBusSubscriptionRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.SubscriptionRulesClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.TopicName, id.SubscriptionName, id.RuleName)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace Subscription Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ServiceBusSubscriptionRuleResource) basicSqlFilter(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) requiresImport(data acceptance.TestData) string {
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
`, r.basicSqlFilter(data))
}

func (r ServiceBusSubscriptionRuleResource) basicSqlFilterUpdated(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) sqlFilterWithAction(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) basicCorrelationFilter(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) correlationFilter(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) correlationFilterUpdated(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) correlationFilterWithAction(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (ServiceBusSubscriptionRuleResource) template(data acceptance.TestData) string {
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
