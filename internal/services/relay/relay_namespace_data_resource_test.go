// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RelayNamespaceDataResource struct{}

func TestAccRelayNamespaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_relay_namespace", "test")
	r := RelayNamespaceDataResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("metric_id").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
			),
		},
	})
}

func (RelayNamespaceDataResource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_relay_namespace" "test" {
  depends_on = [azurerm_relay_namespace.test]

  name                = azurerm_relay_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, RelayNamespaceResource{}.basic(data))
}
