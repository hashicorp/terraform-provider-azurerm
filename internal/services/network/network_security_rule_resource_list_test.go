package network_test

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

func TestAccNetworkSecurityRule_listByNetworkSecurityGroupID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "testlist1")
	r := NetworkSecurityRuleResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_network_security_rule.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_network_security_rule.list",
						map[string]knownvalue.Check{
							"name":                        knownvalue.StringRegexp(regexp.MustCompile("test123")),
							"resource_group_name":         knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"network_security_group_name": knownvalue.StringRegexp(regexp.MustCompile("acceptanceTestSecurityGroup1")),
							"subscription_id":             knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r NetworkSecurityRuleResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_network_security_rule" "list" {
  provider = azurerm
  config {
    network_security_group_id = azurerm_network_security_group.test.id
  }
}
`
}
