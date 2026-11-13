// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMssqlDatabaseSensitivityLabel_listByDatabaseID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_sensitivity_label", "test")
	r := MssqlDatabaseSensitivityLabelResource{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_mssql_database_sensitivity_label.list", 3),
					querycheck.ExpectIdentity("azurerm_mssql_database_sensitivity_label.list", map[string]knownvalue.Check{
						"column_name":              knownvalue.StringExact("EmailAddress"),
						"database_name":            knownvalue.StringExact(fmt.Sprintf("acctest-db-%d", data.RandomInteger)),
						"resource_group_name":      knownvalue.StringExact(fmt.Sprintf("acctestRG-mssql-%d", data.RandomInteger)),
						"schema_name":              knownvalue.StringExact("SalesLT"),
						"sensitivity_label_source": knownvalue.StringExact("current"),
						"server_name":              knownvalue.StringExact(fmt.Sprintf("acctest-sqlserver-%d", data.RandomInteger)),
						"subscription_id":          knownvalue.StringExact(data.Subscriptions.Primary),
						"table_name":               knownvalue.StringExact("Customer"),
					}),
				},
			},
		},
	})
}

func (r MssqlDatabaseSensitivityLabelResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_sensitivity_label" "test" {
  count = 3

  column_name      = ["EmailAddress", "FirstName", "LastName"][count.index]
  database_id      = azurerm_mssql_database.test.id
  schema_name      = "SalesLT"
  table_name       = "Customer"
  information_type = ["Email", "Name", "Name"][count.index]
  label_name       = ["Confidential", "Confidential", null][count.index]
}
`, r.template(data))
}

func (MssqlDatabaseSensitivityLabelResource) listQuery() string {
	return `
list "azurerm_mssql_database_sensitivity_label" "list" {
  provider         = azurerm
  include_resource = true
  config {
    database_id = azurerm_mssql_database.test.id
  }
}
`
}
