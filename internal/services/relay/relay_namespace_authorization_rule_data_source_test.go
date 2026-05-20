// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RelayNamespaceAuthorizationRuleDataResource struct{}

func TestAccRelayNamespaceAuthorizationRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_relay_namespace_authorization_rule", "test")
	r := RelayNamespaceAuthorizationRuleDataResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestrnak-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("listen").Exists(),
				check.That(data.ResourceName).Key("send").Exists(),
				check.That(data.ResourceName).Key("manage").Exists(),
			),
		},
	})
}

func (RelayNamespaceAuthorizationRuleDataResource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_relay_namespace_authorization_rule" "test" {
  depends_on = [azurerm_relay_namespace_authorization_rule.test]

  name                = azurerm_relay_namespace_authorization_rule.test.name
  namespace_name      = azurerm_relay_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, RelayNamespaceAuthorizationRuleResource{}.basic(data))
}
