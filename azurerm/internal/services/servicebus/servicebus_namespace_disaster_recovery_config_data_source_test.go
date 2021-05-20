package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceDisasterRecoveryDataSource struct {
}

func TestAccDataSourceServiceBusNamespaceDisasterRecoveryConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace_disaster_recovery_config", "test")
	r := ServiceBusNamespaceDisasterRecoveryDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("partner_namespace_id").Exists(),
				check.That(data.ResourceName).Key("alias_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("alias_secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_primary_key").Exists(),
				check.That(data.ResourceName).Key("default_secondary_key").Exists(),
			),
		},
	})
}

func (ServiceBusNamespaceDisasterRecoveryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace_disaster_recovery_config" "test" {
  name                = azurerm_servicebus_namespace_disaster_recovery_config.pairing_test.name
  resource_group_name = azurerm_resource_group.primary.name
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
}
`, ServiceBusNamespaceDisasterRecoveryConfigResource{}.basic(data))
}
