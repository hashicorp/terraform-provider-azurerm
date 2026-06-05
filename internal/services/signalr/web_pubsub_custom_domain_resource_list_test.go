package signalr_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccWebPubsubCustomDomainResource_listByWebPubSubID(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_custom_domain", "testlist")
	r := WebPubsubCustomDomainResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_web_pubsub_custom_domain.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_web_pubsub_custom_domain.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(`(?i)webpubsubcustom-domain-`)),
							"web_pubsub_name":     knownvalue.StringRegexp(regexp.MustCompile(`(?i)acctestwebpubsub-`)),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(`(?i)acctestrg-`)),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r WebPubsubCustomDomainResource) basicListQuery() string {
	return `
list "azurerm_web_pubsub_custom_domain" "list" {
  provider = azurerm
  config {
    web_pubsub_id = azurerm_web_pubsub.test.id
  }
}
`
}
