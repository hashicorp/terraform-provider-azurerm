package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusQueueAuthorizationRuleDataSource struct{}

func TestAccDataSourceServiceBusQueueAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
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

func (ServiceBusQueueAuthorizationRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name     = azurerm_servicebus_queue_authorization_rule.test.name
  queue_id = azurerm_servicebus_queue.test.id
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
