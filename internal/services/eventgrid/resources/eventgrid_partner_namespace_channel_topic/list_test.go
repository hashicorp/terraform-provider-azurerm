// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccEventGridPartnerNamespaceChannelTopic_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel_topic", "test")
	r := EventgridPartnerNamespaceChannelTopicResource{}
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
				Query: true,
				Config: fmt.Sprintf(`
list "azurerm_eventgrid_partner_namespace_channel_topic" "test" {
  provider         = azurerm
  include_resource = true
  config {
    partner_namespace_id = "/subscriptions/%s/resourceGroups/acctestRG-partner-%d/providers/Microsoft.EventGrid/partnerNamespaces/acctest-egpn-%d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_eventgrid_partner_namespace_channel_topic.test", 1),
					querycheck.ExpectIdentity("azurerm_eventgrid_partner_namespace_channel_topic.test", map[string]knownvalue.Check{
						"name":                   knownvalue.StringExact(fmt.Sprintf("acctest-channel-%d", data.RandomInteger)),
						"partner_namespace_name": knownvalue.StringExact(fmt.Sprintf("acctest-egpn-%d", data.RandomInteger)),
						"resource_group_name":    knownvalue.StringExact(fmt.Sprintf("acctestRG-partner-%d", data.RandomInteger)),
						"subscription_id":        knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}
