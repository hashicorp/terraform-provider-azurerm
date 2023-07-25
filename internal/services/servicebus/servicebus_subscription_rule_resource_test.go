// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusSubscriptionRuleResource struct{}

func TestAccServiceBusSubscriptionRule_basicSqlFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccServiceBusSubscriptionRule_basicCorrelationFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicCorrelationFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_sqlFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sqlFilterWithAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_correlationFilterWithAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.correlationFilterWithAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_sqlFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sql_filter").HasValue("2=2"),
			),
		},
		{
			Config: r.basicSqlFilterUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sql_filter").HasValue("3=3"),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_correlationFilterUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.correlationFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("correlation_filter.0.message_id").HasValue("test_message_id"),
				check.That(data.ResourceName).Key("correlation_filter.0.reply_to").HasValue(""),
			),
		},
		{
			Config: r.correlationFilterUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSqlFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicCorrelationFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccServiceBusSubscriptionRule_correlationFilterWhiteSpace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServiceBusSubscriptionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.correlationFilterWhiteSpace(data),
			ExpectError: regexp.MustCompile("expanding `correlation_filter`: at least one property must not be empty in the `correlation_filter` block"),
		},
	})
}

func (t ServiceBusSubscriptionRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := subscriptions.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.SubscriptionsClient.RulesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ServiceBusSubscriptionRuleResource) basicSqlFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "SqlFilter"
  sql_filter      = "2=2"
}
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "import" {
  name            = azurerm_servicebus_subscription_rule.test.name
  subscription_id = azurerm_servicebus_subscription_rule.test.subscription_id
  filter_type     = azurerm_servicebus_subscription_rule.test.filter_type
  sql_filter      = azurerm_servicebus_subscription_rule.test.sql_filter
}
`, r.basicSqlFilter(data))
}

func (r ServiceBusSubscriptionRuleResource) basicSqlFilterUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "SqlFilter"
  sql_filter      = "3=3"
}
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) sqlFilterWithAction(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "SqlFilter"
  sql_filter      = "2=2"
  action          = "SET Test='true'"
}
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) basicCorrelationFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "CorrelationFilter"

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
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "CorrelationFilter"

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
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  filter_type     = "CorrelationFilter"

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
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  action          = "SET Test='true'"
  filter_type     = "CorrelationFilter"

  correlation_filter {
    correlation_id = "test_correlation_id"
    message_id     = "test_message_id"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ServiceBusSubscriptionRuleResource) correlationFilterWhiteSpace(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription_rule" "test" {
  name            = "acctestservicebusrule-%d"
  subscription_id = azurerm_servicebus_subscription.test.id
  action          = "SET Test='true'"
  filter_type     = "CorrelationFilter"

  correlation_filter {
    label = ""
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
  name         = "acctestservicebustopic-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_subscription" "test" {
  name               = "acctestservicebussubscription-%d"
  topic_id           = azurerm_servicebus_topic.test.id
  max_delivery_count = 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
