// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServiceBusNamespaceDisasterRecoveryDataSource struct{}

func TestAccDataSourceServiceBusNamespaceDisasterRecoveryConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace_disaster_recovery_config", "test")
	r := ServiceBusNamespaceDisasterRecoveryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("partner_namespace_id").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
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
  name         = azurerm_servicebus_namespace_disaster_recovery_config.pairing_test.name
  namespace_id = azurerm_servicebus_namespace.primary_namespace_test.id
}
`, ServiceBusNamespaceDisasterRecoveryConfigResource{}.basic(data))
}
