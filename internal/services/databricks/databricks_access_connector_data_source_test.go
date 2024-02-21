// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DatabricksAccessConnectorDataSource struct{}

func TestAccDatabricksAccessConnectorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccDatabricksAccessConnectorDataSource_systemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.systemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDatabricksAccessConnectorDataSource_userIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.userIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.0").Exists(),
			),
		},
	})
}

func (DatabricksAccessConnectorDataSource) basic(data acceptance.TestData) string {
	template := DatabricksAccessConnectorResource{}.basic(data)
	return fmt.Sprintf(`
%s
data "azurerm_databricks_access_connector" "test" {
  name                = azurerm_databricks_access_connector.test.name
  resource_group_name = azurerm_databricks_access_connector.test.resource_group_name
}
`, template)
}

func (DatabricksAccessConnectorDataSource) systemIdentity(data acceptance.TestData) string {
	template := DatabricksAccessConnectorResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_databricks_access_connector" "test" {
  name                = azurerm_databricks_access_connector.test.name
  resource_group_name = azurerm_databricks_access_connector.test.resource_group_name
}
`, template)
}

func (DatabricksAccessConnectorDataSource) userIdentity(data acceptance.TestData) string {
	template := DatabricksAccessConnectorResource{}.identityUserAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_databricks_access_connector" "test" {
  name                = azurerm_databricks_access_connector.test.name
  resource_group_name = azurerm_databricks_access_connector.test.resource_group_name
}
`, template)
}
