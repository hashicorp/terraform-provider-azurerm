// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

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

func TestAccPaloAltoNextGenerationFirewallMetrics_list_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := PaloAltoNextGenerationFirewallMetricsResource{}
	listResourceAddress := "azurerm_palo_alto_next_generation_firewall_metrics.list"

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
					// the Metrics Object is a singleton beneath the Firewall, so exactly one result is expected
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"firewall_name":       knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r PaloAltoNextGenerationFirewallMetricsResource) basicQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_palo_alto_next_generation_firewall_metrics" "list" {
  provider = azurerm
  config {
    firewall_id = "/subscriptions/%[1]s/resourceGroups/acctestRG-PANGFWMETRICS-%[2]d/providers/PaloAltoNetworks.Cloudngfw/firewalls/acctest-ngfwvn-%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
