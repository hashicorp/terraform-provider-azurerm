// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dns_test

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

func TestAccDnsARecord_listByDnsZoneID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "list")
	r := DnsARecordResource{}

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
					querycheck.ExpectLength("azurerm_dns_a_record.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_dns_a_record.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"dns_zone_name":       knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
							"record_type":         knownvalue.StringExact("A"),
						},
					),
				},
			},
		},
	})
}

func (DnsARecordResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_a_record" "test" {
  count               = 3
  name                = "myarecord%d${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.${count.index}"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r DnsARecordResource) basicQuery() string {
	return `
list "azurerm_dns_a_record" "list" {
  provider = azurerm
  config {
    dns_zone_id = azurerm_dns_zone.test.id
  }
}
`
}
