package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorCustomDomainSecretValidatorResource struct{}

func TestAccCdnFrontdoorCustomDomainSecretValidator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain_secret_validator", "test")
	r := CdnFrontdoorCustomDomainSecretValidatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontdoorCustomDomainSecretValidatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	// This is not a real resource so it does not really exist in Azure, but if it exists in state that means that it was created successfully
	_, err := parse.FrontdoorCustomDomainSecretID(state.ID)
	if err != nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r CdnFrontdoorCustomDomainSecretValidatorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "accTestOriginGroup-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_count                       = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "accTestOrigin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id

  health_probes_enabled          = true
  certificate_name_check_enabled = false
  host_name                      = join(".", ["fabrikam", azurerm_dns_zone.test.name])
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "accTestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "accTestCustomDomain-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dns_zone_id = azurerm_dns_zone.test.id
  host_name   = join(".", ["fabrikam", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                            = "accTestRoute-%[1]d"
  cdn_frontdoor_endpoint_id       = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id   = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids        = [azurerm_cdn_frontdoor_origin.test.id]
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.test.id]
  enabled                         = true

  link_to_default_domain_enabled = true
  https_redirect_enabled         = true
  forwarding_protocol            = "HttpsOnly"
  patterns_to_match              = ["/*"]
  supported_protocols            = ["Http", "Https"]

  cache {
    compression_enabled       = true
    content_types_to_compress = ["text/html", "text/javascript", "text/xml"]
  }
}

resource "azurerm_dns_txt_record" "fabrikam" {
  name                = join(".", ["_dnsauth", "fabrikam"])
  zone_name           = azurerm_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 3600

  record {
    value = azurerm_cdn_frontdoor_custom_domain.test.validation_properties.0.validation_token
  }
}

resource "azurerm_cdn_frontdoor_custom_domain_txt_validator" "fabrikam" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.test.id
  dns_txt_record_id              = azurerm_dns_txt_record.fabrikam.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainSecretValidatorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cdn_frontdoor_custom_domain_secret_validator" "fabrikam" {
  cdn_frontdoor_route_id                        = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_custom_domain_ids               = [azurerm_cdn_frontdoor_custom_domain.test.id]
  cdn_frontdoor_custom_domain_txt_validator_ids = [azurerm_cdn_frontdoor_custom_domain_txt_validator.fabrikam.id]
}
`, template)
}
