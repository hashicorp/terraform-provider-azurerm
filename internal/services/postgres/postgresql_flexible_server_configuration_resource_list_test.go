package postgres_test

import (
	"context"
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

func TestAccPostgresqlFlexibleServerConfiguration_listByFlexibleServerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "testlist1")
	r := PostgresqlFlexibleServerConfigurationResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data, "backslash_quote", "on"),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_postgresql_flexible_server_configuration.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_postgresql_flexible_server_configuration.list",
						map[string]knownvalue.Check{
							"name":                 knownvalue.StringExact("backslash_quote"),
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

func (r PostgresqlFlexibleServerConfigurationResource) basicQuery() string {
	return `
list "azurerm_postgresql_flexible_server_configuration" "list" {
  provider = azurerm
  config {
    server_id = azurerm_postgresql_flexible_server.test.id
  }
}
`
}
