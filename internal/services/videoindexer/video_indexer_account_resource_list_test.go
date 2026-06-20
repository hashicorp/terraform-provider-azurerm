// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package videoindexer_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccVideoIndexerAccount_list_basic(t *testing.T) {
	r := AccountResource{}
	listResourceAddress := "azurerm_video_indexer_account.list"

	data := acceptance.BuildTestData(t, "azurerm_video_indexer_account", "testlist")

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(`acctestvi[0-2]-`)),
							"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestRG-VI-%d", data.RandomInteger)),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryIncludeResource(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r AccountResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VI-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestvi%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_video_indexer_account" "test" {
  count = 3

  name                = "acctestvi${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage {
    storage_account_id = azurerm_storage_account.test.id
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r AccountResource) basicQuery() string {
	return `
list "azurerm_video_indexer_account" "list" {
  provider = azurerm
  config {}
}
`
}

func (r AccountResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_video_indexer_account" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}

func (r AccountResource) basicQueryIncludeResource() string {
	return `
list "azurerm_video_indexer_account" "list" {
  provider         = azurerm
  include_resource = true
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
