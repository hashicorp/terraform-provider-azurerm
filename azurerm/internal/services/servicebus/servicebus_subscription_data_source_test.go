package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ServiceBusSubscriptionDataSource struct {
}

func TestAccDataSourceServiceBusSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("max_delivery_count").Exists(),
			),
		},
	})
}

func (ServiceBusSubscriptionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_servicebus_subscription" "test" {
  name                = "acctestservicebussubscription-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  max_delivery_count  = 10
}
data "azurerm_servicebus_subscription" "test" {
  name                = azurerm_servicebus_subscription.test.name
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name
  topic_name          = azurerm_servicebus_topic.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
