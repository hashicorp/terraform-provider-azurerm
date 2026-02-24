// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

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

func TestAccPrivateDnsZone_list_basic(t *testing.T) {
	r := PrivateDnsZoneResource{}

	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test")
	listResourceAddress := "azurerm_private_dns_zone.list"

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
				Config: r.basicQueryByResourceGroupNameIncludeResource(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectResourceKnownValues(listResourceAddress, queryfilter.ByDisplayName(knownvalue.StringRegexp(regexp.MustCompile(`acctestzone-0-`))), []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("soa_record"),
							KnownValue: knownvalue.ListSizeExact(1),
						},
					}),
				},
			},
		},
	})
}

func (r PrivateDnsZoneResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  count = 3

  name                = "acctestzone-${count.index}-%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDnsZoneResource) basicQuery() string {
	return `
list "azurerm_private_dns_zone" "list" {
  provider = azurerm
  config {}
}
`
}

func (r PrivateDnsZoneResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_private_dns_zone" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}
`, data.RandomInteger)
}

func (r PrivateDnsZoneResource) basicQueryByResourceGroupNameIncludeResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_private_dns_zone" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
  include_resource = true
}
`, data.RandomInteger)
}
