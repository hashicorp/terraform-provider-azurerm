// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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

func TestAccMySqlFlexibleServer_list_no_config(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist1")
	r := MySqlFlexibleServerResource{}

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_mysql_flexible_server.list", 1), // expect at least the 1 we created
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_server.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary)},
					),
				},
			},
		},
	})
}

func TestAccMySqlFlexibleServer_list_by_resource_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist")
	r := MySqlFlexibleServerResource{}

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
					querycheck.ExpectLength("azurerm_mysql_flexible_server.list", 1), // only 1 should be returned
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_server.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary)},
					),
				},
			},
		},
	})
}

func (r MySqlFlexibleServerResource) basicListQuery() string {
	return `
list "azurerm_mysql_flexible_server" "list" {
  provider = azurerm
  config {}
}
`
}

func (r MySqlFlexibleServerResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_mysql_flexible_server" "list" {
  provider = azurerm
  config {
	subscription_id     = "%s"
	resource_group_name = "acctestRG-mysql-%d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
