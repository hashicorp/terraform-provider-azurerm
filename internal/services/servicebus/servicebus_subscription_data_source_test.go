// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusSubscriptionDataSource struct{}

func TestAccDataSourceServiceBusSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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
  name         = "acctestservicebustopic-%d"
  namespace_id = azurerm_servicebus_namespace.test.id
}
resource "azurerm_servicebus_subscription" "test" {
  name               = "acctestservicebussubscription-%d"
  topic_id           = azurerm_servicebus_topic.test.id
  max_delivery_count = 10
}
data "azurerm_servicebus_subscription" "test" {
  name     = azurerm_servicebus_subscription.test.name
  topic_id = azurerm_servicebus_topic.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
