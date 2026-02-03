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
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// mysql config properties always exist, even if not set in terraform, and they are all returned in the list
					// no need to check anything specific here, the list is always the entire set of available properties
					querycheck.ExpectLengthAtLeast("azurerm_mysql_flexible_server_configuration.list", 1),
				},
			},
		},
	})
}

func (r MysqlFlexibleServerConfigurationResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
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
  config {
    flexible_server_id = list.azurerm_mysql_flexible_server.list.data[0].state.id
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
