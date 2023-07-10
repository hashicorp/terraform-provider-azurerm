// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusQueueAuthorizationRuleDataSource struct{}

func TestAccDataSourceServiceBusQueueAuthorizationRule_byId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.byId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("queue_name").Exists(),
				check.That(data.ResourceName).Key("queue_id").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").HasValue(""),
			),
		},
	})
}
func TestAccDataSourceServiceBusQueueAuthorizationRule_byName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.byName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("queue_name").Exists(),
				check.That(data.ResourceName).Key("queue_id").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceServiceBusQueueAuthorizationRule_queueName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.queueName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("queue_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceServiceBusQueueAuthorizationRule_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.queueAliasPolicy(data),
		},
		{
			Config: r.queueAliasPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
	})
}

func (ServiceBusQueueAuthorizationRuleDataSource) byId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name     = azurerm_servicebus_queue_authorization_rule.test.name
  queue_id = azurerm_servicebus_queue.test.id
}
`, ServiceBusQueueAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusQueueAuthorizationRuleDataSource) byName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = azurerm_servicebus_queue_authorization_rule.test.name
  resource_group_name = azurerm_servicebus_namespace.test.resource_group_name
  namespace_name      = azurerm_servicebus_namespace.test.name
  queue_name          = azurerm_servicebus_queue.test.name
}
`, ServiceBusQueueAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusQueueAuthorizationRuleDataSource) queueName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = azurerm_servicebus_queue_authorization_rule.test.name
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  queue_name          = azurerm_servicebus_queue.test.name
}
`, ServiceBusQueueAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusQueueAuthorizationRuleDataSource) queueAliasPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name     = azurerm_servicebus_queue_authorization_rule.test.name
  queue_id = azurerm_servicebus_queue.example.id
}
`, ServiceBusQueueAuthorizationRuleResource{}.withAliasConnectionString(data))
}
