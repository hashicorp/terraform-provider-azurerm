package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorRouteDataSource struct{}

func TestAccCdnFrontDoorRouteDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_route", "test")
	d := CdnFrontDoorRouteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("forwarding_protocol").HasValue("HttpsOnly"),
			),
		},
	})
}

func TestAccCdnFrontDoorRouteDataSource_useCase(t *testing.T) {
	// this test case was added to validate customer use case in issue #21639
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_route", "test")
	d := CdnFrontDoorRouteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.stepOne(data),
		},
		{
			Config: d.stepTwo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("forwarding_protocol").HasValue("HttpsOnly"),
			),
		},
	})
}

func (CdnFrontDoorRouteDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_route" "test" {
  name                      = azurerm_cdn_frontdoor_route.test.name
  cdn_frontdoor_endpoint_id = azurerm_cdn_frontdoor_endpoint.test.id
}
`, CdnFrontDoorRouteResource{}.complete(data))
}

func (CdnFrontDoorRouteDataSource) stepOne(data acceptance.TestData) string {
	return fmt.Sprintf(`
# NOTE: This simulates a configuration being created in a different config
%s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctestRoute-%[2]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]

  # cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.test.id]

  enabled                    = true
  forwarding_protocol        = "HttpsOnly"
  https_redirect_enabled     = false
  patterns_to_match          = ["/*"]
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.test.id]
  supported_protocols        = ["Https"]

  cache {
    query_strings                 = ["bar"]
    query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
  }
}
`, CdnFrontDoorRouteResource{}.template(data), data.RandomInteger)
}

func (r CdnFrontDoorRouteDataSource) stepTwo(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-contoso-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  host_name                = "test.%[2]d.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

data "azurerm_cdn_frontdoor_endpoint" "test" {
  name                = "accTestEndpoint-%[2]d"
  profile_name        = "accTestProfile-%[2]d"
  resource_group_name = "acctestRG-cdn-afdx-%[2]d"
}

data "azurerm_cdn_frontdoor_route" "test" {
  name                      = "acctestRoute-%[2]d"
  cdn_frontdoor_endpoint_id = data.azurerm_cdn_frontdoor_endpoint.test.id
}

resource "azurerm_cdn_frontdoor_custom_domain_association" "test" {
  cdn_frontdoor_route_id          = data.azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.test.id]
  link_to_default_domain          = false
}
`, r.stepOne(data), data.RandomInteger)
}
