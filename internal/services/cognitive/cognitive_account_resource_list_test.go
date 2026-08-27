// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccCognitiveAccount_list_basic(t *testing.T) {
	r := CognitiveAccountResource{}
	listResourceAddress := "azurerm_cognitive_account.list"

	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test1")

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
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryWithIncludeResource(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectResourceKnownValues(listResourceAddress, queryfilter.ByDisplayName(knownvalue.StringRegexp(regexp.MustCompile("acctestcogacclist-"))), []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("primary_access_key"),
							KnownValue: knownvalue.NotNull(),
						},
						{
							Path:       tfjsonpath.New("secondary_access_key"),
							KnownValue: knownvalue.NotNull(),
						},
					}),
				},
			},
		},
	})
}

func (r CognitiveAccountResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-list-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test1" {
  name                = "acctestcogacclist-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account" "test2" {
  name                = "acctestcogacclist2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account" "test3" {
  name                = "acctestcogacclist3-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"
}
    `, data.RandomInteger, data.Locations.Primary)
}

func (r CognitiveAccountResource) basicQuery() string {
	return `
list "azurerm_cognitive_account" "list" {
  provider = azurerm
  config {}
}
    `
}

func (r CognitiveAccountResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_cognitive_account" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-cognitive-list-%[1]d"
  }
}
    `, data.RandomInteger)
}

func (r CognitiveAccountResource) basicQueryWithIncludeResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_cognitive_account" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-cognitive-list-%[1]d"
  }
  include_resource = true
}
    `, data.RandomInteger)
}
