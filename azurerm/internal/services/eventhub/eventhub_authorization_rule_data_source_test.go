package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type EventHubAuthorizationRuleDataSource struct {
}

func TestAccEventHubAuthorizationRuleDataSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.base(data, true, true, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("eventhub_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func TestAccEventHubAuthorizationRuleDataSource_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
	})
}

func (EventHubAuthorizationRuleDataSource) base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_authorization_rule" "test" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventHubAuthorizationRuleResource{}.base(data, listen, send, manage))
}

func (EventHubAuthorizationRuleDataSource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_authorization_rule" "test" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventHubAuthorizationRuleResource{}.withAliasConnectionString(data))
}
