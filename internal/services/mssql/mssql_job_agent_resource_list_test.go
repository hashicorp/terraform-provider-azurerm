package mssql_test

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

func TestAccMssqlJobAgent_listByServerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_agent", "testlist1")
	r := MssqlJobAgentResource{}

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
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_mssql_job_agent.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_mssql_job_agent.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"server_name":         knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r MssqlJobAgentResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_mssql_job_agent" "list" {
  provider = azurerm
  config {
    mssql_server_id = azurerm_mssql_server.test.id
  }
}
`
}
