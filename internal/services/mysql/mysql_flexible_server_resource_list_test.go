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
	data1 := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist1")
	data2 := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist2")
	r := MySqlFlexibleServerResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicListResources(data1, data2), // provision first server + RG
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_mysql_flexible_server.list", 2), // expect at least the 2 we created
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_server.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data1.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data1.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data1.Subscriptions.Primary)},
					),
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_server.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data2.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data2.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data2.Subscriptions.Primary)},
					),
				},
			},
		},
	})
}

func TestAccMySqlFlexibleServer_list_by_resource_group(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist1")
	data2 := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "testlist2")
	r := MySqlFlexibleServerResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicListResources(data1, data2), // provision first server + RG
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data1),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_mysql_flexible_server.list", 1), // only 1 should be returned
					querycheck.ExpectIdentity(
						"azurerm_mysql_flexible_server.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data1.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data1.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data1.Subscriptions.Primary)},
					),
				},
			},
		},
	})
}

func (r MySqlFlexibleServerResource) basicListResources(data1 acceptance.TestData, data2 acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-mysql-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-mysql-%d"
  location = "%s"
}

resource "azurerm_mysql_flexible_server" "test1" {
  name                   = "acctest-fs-%d-1"
  resource_group_name    = azurerm_resource_group.test1.name
  location               = azurerm_resource_group.test1.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
}

resource "azurerm_mysql_flexible_server" "test2" {
  name                   = "acctest-fs-%d-2"
  resource_group_name    = azurerm_resource_group.test2.name
  location               = azurerm_resource_group.test2.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx456"
  sku_name               = "B_Standard_B1ms"
}
`, data1.RandomInteger, data1.Locations.Ternary, data2.RandomInteger, data2.Locations.Ternary, data1.RandomInteger, data2.RandomInteger)
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
