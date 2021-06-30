package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceDataSource struct {
}

func TestAccDataSourceServiceBusNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("capacity").Exists(),
				check.That(data.ResourceName).Key("default_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_primary_key").Exists(),
				check.That(data.ResourceName).Key("default_secondary_key").Exists(),
			),
		},
	})
}

func TestAccDataSourceServiceBusNamespace_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("capacity").Exists(),
				check.That(data.ResourceName).Key("default_primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("default_primary_key").Exists(),
				check.That(data.ResourceName).Key("default_secondary_key").Exists(),
			),
		},
	})
}

func (ServiceBusNamespaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServiceBusNamespaceResource{}.basic(data))
}

func (ServiceBusNamespaceDataSource) premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ServiceBusNamespaceResource{}.premium(data))
}
