package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ServiceBusQueueDataSource struct {
}

func TestAccDataSourceServiceBusQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_queue", "test")
	r := ServiceBusQueueDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("auto_delete_on_idle").Exists(),
				check.That(data.ResourceName).Key("dead_lettering_on_message_expiration").Exists(),
				check.That(data.ResourceName).Key("default_message_ttl").Exists(),
				check.That(data.ResourceName).Key("duplicate_detection_history_time_window").Exists(),
				check.That(data.ResourceName).Key("enable_batched_operations").Exists(),
				check.That(data.ResourceName).Key("enable_express").Exists(),
				check.That(data.ResourceName).Key("enable_partitioning").Exists(),
				check.That(data.ResourceName).Key("lock_duration").Exists(),
				check.That(data.ResourceName).Key("max_delivery_count").Exists(),
				check.That(data.ResourceName).Key("max_size_in_megabytes").Exists(),
				check.That(data.ResourceName).Key("requires_duplicate_detection").Exists(),
				check.That(data.ResourceName).Key("requires_session").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
	})
}

func (ServiceBusQueueDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_queue" "test" {
  name                = azurerm_servicebus_queue.test.name
  namespace_name      = azurerm_servicebus_queue.test.namespace_name
  resource_group_name = azurerm_servicebus_queue.test.resource_group_name
}
`, ServiceBusQueueResource{}.basic(data))
}
