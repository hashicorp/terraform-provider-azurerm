package postgres_test

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

func TestAccPostgresqlFlexibleServerFirewallRule_listByFlexibleServerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rule", "testlist1")
	r := PostgresqlFlexibleServerFirewallRuleResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_postgresql_flexible_server_firewall_rule.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_postgresql_flexible_server_firewall_rule.list",
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

func (PostgresqlFlexibleServerFirewallRuleResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_firewall_rule" "test" {
  count            = 3
  name             = "acctest-FSFR-%d${count.index}"
  server_id        = azurerm_postgresql_flexible_server.test.id
  start_ip_address = "122.122.0.${count.index}"
  end_ip_address   = "122.122.0.${count.index}"
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerFirewallRuleResource) basicQuery() string {
	return `
list "azurerm_postgresql_flexible_server_firewall_rule" "list" {
  provider = azurerm
  config {
    postgresql_flexible_server_id = azurerm_postgresql_flexible_server.test.id
  }
}
`
}
