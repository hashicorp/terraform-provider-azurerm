package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceAuthorizationRuleDataSource struct {
}

func TestAccDataSourceServiceBusNamespaceRule_basic(t *testing.T) {
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
				check.That(data.ResourceName).Key("alias_primary_connection_string").HasValue(""),
				check.That(data.ResourceName).Key("alias_secondary_connection_string").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespaceRule_aliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// `primary_connection_string_alias` and `secondary_connection_string_alias` are still `nil` while `data.azurerm_servicebus_namespace_authorization_rule` is retrieving acceptance. since `azurerm_servicebus_namespace_disaster_recovery_config` hasn't been created.
			// So these two properties should be checked in the second run.
			// And `depends_on` cannot be applied to `azurerm_servicebus_namespace_authorization_rule`.
			// Because it would throw error message `BreakPairing operation is only allowed on primary namespace with valid secondary namespace.` while destroying `azurerm_servicebus_namespace_disaster_recovery_config` if `depends_on` is applied.
			Config: r.namespaceAliasPolicy(data),
		},
		{
			Config: r.namespaceAliasPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("alias_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("alias_secondary_connection_string").Exists(),
			),
		},
	})
}

func (ServiceBusNamespaceAuthorizationRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = azurerm_servicebus_namespace_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServiceBusNamespaceAuthorizationRuleResource{}.base(data, true, true, true))
}

func (ServiceBusNamespaceAuthorizationRuleDataSource) namespaceAliasPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = azurerm_servicebus_namespace_authorization_rule.example.name
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
  resource_group_name = azurerm_resource_group.primary.name
}
`, ServiceBusNamespaceAuthorizationRuleResource{}.withAliasConnectionString(data))
}
