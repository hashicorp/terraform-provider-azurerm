// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RelayHybridConnectionDataResource struct{}

func TestAccRelayHybridConnectionDataDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_relay_hybrid_connection", "test")
	r := RelayHybridConnectionDataResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("default"),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("requires_client_authorization").Exists(),
				check.That(data.ResourceName).Key("user_metadata").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (RelayHybridConnectionDataResource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_relay_hybrid_connection" "test" {
  depends_on = [azurerm_relay_hybrid_connection.test]

  name       		   = azurerm_relay_hybrid_connection.test.name
  relay_namespace_name = azurerm_relay_namespace.test.name
  resource_group_name  = azurerm_resource_group.test.name
}
`, RelayHybridConnectionResource{}.basic(data))
}
