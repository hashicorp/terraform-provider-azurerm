package mysql_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMySqlFlexibleDatabase_list_by_resource_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "testlist")
	r := MysqlFlexibleDatabaseResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// mysql will have several admin db by default, so there will be more than just the one provisioned
					querycheck.ExpectLengthAtLeast("azurerm_mysql_flexible_database.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_database.list",
						map[string]knownvalue.Check{
							"name":                 knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"flexible_server_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":      knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r MysqlFlexibleDatabaseResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_mysql_flexible_server" "list" {
  provider         = azurerm
  include_resource = true
  config {
    subscription_id     = "%s"
    resource_group_name = "acctestRG-%d"
  }
}

list "azurerm_mysql_flexible_database" "list" {
  provider = azurerm
  config {
    flexible_server_id = list.azurerm_mysql_flexible_server.list.data[0].state.id
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
