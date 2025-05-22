// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusTopicDataSource struct{}

func TestAccDataSourceServiceBusTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_topic", "test")
	r := ServiceBusTopicDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("auto_delete_on_idle").Exists(),
				check.That(data.ResourceName).Key("default_message_ttl").Exists(),
				check.That(data.ResourceName).Key("duplicate_detection_history_time_window").Exists(),
				check.That(data.ResourceName).Key("batched_operations_enabled").Exists(),
				check.That(data.ResourceName).Key("express_enabled").Exists(),
				check.That(data.ResourceName).Key("partitioning_enabled").Exists(),
				check.That(data.ResourceName).Key("max_size_in_megabytes").Exists(),
				check.That(data.ResourceName).Key("requires_duplicate_detection").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
				check.That(data.ResourceName).Key("support_ordering").Exists(),
			),
		},
	})
}

func (ServiceBusTopicDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_topic" "test" {
  name         = azurerm_servicebus_topic.test.name
  namespace_id = azurerm_servicebus_namespace.test.id
}
`, ServiceBusTopicResource{}.basic(data))
}
