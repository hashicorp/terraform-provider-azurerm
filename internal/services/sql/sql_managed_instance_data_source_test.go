// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SqlManagedInstanceDataSource struct{}

func TestAccDataSourceSqlManagedInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_managed_instance", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SqlManagedInstanceDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (d SqlManagedInstanceDataSource) basic(data acceptance.TestData) string {
	template := SqlManagedInstanceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_managed_instance" "test" {
  name                = azurerm_sql_managed_instance.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
