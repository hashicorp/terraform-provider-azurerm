// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceAuthorizationRuleDataSource struct{}

func TestAccDataSourceServiceBusNamespaceAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").HasValue(""),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespaceAuthorizationRule_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// `primary_connection_string_alias` and `secondary_connection_string_alias` are still `nil` while `data.azurerm_servicebus_namespace_authorization_rule` is retrieving acceptance.
			// `depends_on` cannot be applied to `azurerm_servicebus_namespace_authorization_rule` as throws error message `BreakPairing operation is only allowed on primary namespace...` while destroying disaster recovery config, so these two properties should be checked in the second run.
			Config: r.namespaceAliasPolicy(data),
		},
		{
			Config: r.namespaceAliasPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
	})
}

func (ServiceBusNamespaceAuthorizationRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace_authorization_rule" "test" {
  name         = azurerm_servicebus_namespace_authorization_rule.test.name
  namespace_id = azurerm_servicebus_namespace.test.id
}
`, ServiceBusNamespaceAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusNamespaceAuthorizationRuleDataSource) namespaceAliasPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_servicebus_namespace_authorization_rule" "test" {
  name         = azurerm_servicebus_namespace_authorization_rule.test.name
  namespace_id = azurerm_servicebus_namespace.primary_namespace_test.id
}
`, ServiceBusNamespaceAuthorizationRuleResource{}.withAliasConnectionString(data))
}
