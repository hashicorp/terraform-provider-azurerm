package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ServiceBusQueueAuthorizationRuleDataSource struct {
}

func TestAccDataSourceServiceBusQueueAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func (ServiceBusQueueAuthorizationRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = azurerm_servicebus_queue_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_queue_authorization_rule.test.namespace_name
  resource_group_name = azurerm_servicebus_queue_authorization_rule.test.resource_group_name
  queue_name          = azurerm_servicebus_queue_authorization_rule.test.queue_name
}
`, ServiceBusQueueAuthorizationRuleResource{}.base(data, true, true, true))
}
