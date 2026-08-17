package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMySqlFlexibleServerConfiguration_list_by_resource_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "testlist")
	r := MysqlFlexibleServerConfigurationResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.characterSetServer(data),
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data, 50),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// mysql config properties always exist, even if not set in terraform, and they are all returned in the list
					// since we limit to 50, should have exact 50 here (over 300 configs get returned)
					querycheck.ExpectLength("azurerm_mysql_flexible_server_configuration.list", 50),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data, 500),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// limit set high so all configs returned, around 390 normally so check at least 300
					querycheck.ExpectLengthAtLeast("azurerm_mysql_flexible_server_configuration.list", 300),
				},
			},
		},
	})
}

func (r MysqlFlexibleServerConfigurationResource) basicListQueryByResourceGroupName(data acceptance.TestData, limit int16) string {
	return fmt.Sprintf(`
list "azurerm_mysql_flexible_server" "list" {
  provider         = azurerm
  include_resource = true
  config {
    subscription_id     = "%s"
    resource_group_name = "acctestRG-%d"
  }
}

list "azurerm_mysql_flexible_server_configuration" "list" {
  provider = azurerm
  limit    = %d
  config {
    flexible_server_id = list.azurerm_mysql_flexible_server.list.data[0].state.id
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, limit)
}
