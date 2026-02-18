package privatedns_test

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

func TestAccPrivateDnsARecord_listByPrivateDnsZoneID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_a_record", "testlist1")
	r := PrivateDnsARecordResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_private_dns_a_record.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_private_dns_a_record.list",
						map[string]knownvalue.Check{
							"name":                  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":   knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"private_dns_zone_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":       knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r PrivateDnsARecordResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_private_dns_a_record" "list" {
  provider = azurerm
  config {
    private_dns_zone_id = azurerm_private_dns_zone.test.id
  }
}
`
}
