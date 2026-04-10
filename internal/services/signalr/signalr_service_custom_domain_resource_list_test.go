package signalr_test

import (
	"context"
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

func TestAccSignalrServiceCustomDomainResource_listBySignalRServiceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_custom_domain", "testlist")
	r := SignalrServiceCustomDomainResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_signalr_service_custom_domain.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_signalr_service_custom_domain.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile("signalrcustom-domain-")),
							"signal_r_name":       knownvalue.StringRegexp(regexp.MustCompile("acctestSignalR-")),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile("acctestRG-")),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
					querycheck.ExpectResourceKnownValues(
						"azurerm_signalr_service_custom_domain.list",
						queryfilter.ByDisplayName(knownvalue.StringRegexp(regexp.MustCompile("signalrcustom-domain-"))),
						[]querycheck.KnownValueCheck{
							{
								Path:       tfjsonpath.New("domain_name"),
								KnownValue: knownvalue.StringExact("signalr.wpstftestzone.com"),
							},
							{
								Path:       tfjsonpath.New("signalr_custom_certificate_id"),
								KnownValue: knownvalue.StringRegexp(regexp.MustCompile("/customCertificates/signalr-cert-")),
							},
						},
					),
				},
			},
		},
	})
}

func (r SignalrServiceCustomDomainResource) basicListQuery() string {
	return `
list "azurerm_signalr_service_custom_domain" "list" {
  provider         = azurerm
  include_resource = true
  config {
    signalr_service_id = azurerm_signalr_service.test.id
  }
}
`
}
