// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusTopicAuthorizationRuleDataSource struct{}

func TestAccDataSourceServiceBusTopicAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_topic_authorization_rule", "test")
	r := ServiceBusTopicAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("topic_id").Exists(),
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

func TestAccDataSourceServiceBusTopicAuthorizationRule_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_topic_authorization_rule", "test")
	r := ServiceBusTopicAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.topicAliasPolicy(data),
		},
		{
			Config: r.topicAliasPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
	})
}

func (ServiceBusTopicAuthorizationRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = azurerm_servicebus_topic_authorization_rule.test.name
  topic_id = azurerm_servicebus_topic.test.id
}
`, ServiceBusTopicAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusTopicAuthorizationRuleDataSource) topicAliasPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = azurerm_servicebus_topic_authorization_rule.test.name
  topic_id = azurerm_servicebus_topic.example.id
}
`, ServiceBusTopicAuthorizationRuleResource{}.withAliasConnectionString(data))
}
