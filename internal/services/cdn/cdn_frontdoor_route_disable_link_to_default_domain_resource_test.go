// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorRouteDisableLinkToDefaultDomainResource struct{}

func TestAccCdnFrontDoorRouteDisableLinkToDefaultDomain_basic(t *testing.T) {
	if !features.FourPointOhBeta() {
		data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_route_disable_link_to_default_domain", "test")
		r := CdnFrontDoorRouteDisableLinkToDefaultDomainResource{}
		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.basic(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
				ExpectNonEmptyPlan: true, // since this resource actually modifies the routes 'linked to default domain' field
			},
		})
	} else {
		t.Skip("the 'azurerm_cdn_frontdoor_route_disable_link_to_default_domain' resource is no longer supported in the AzureRM provider v4.0, please remove this test case")
	}
}

func (r CdnFrontDoorRouteDisableLinkToDefaultDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorRoutesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontDoorRouteDisableLinkToDefaultDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_route_disable_link_to_default_domain" "test" {
  cdn_frontdoor_route_id          = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.test.id]
}
`, template)
}

func (r CdnFrontDoorRouteDisableLinkToDefaultDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctest-dns-zone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-profile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "accTest-origin-group-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-origin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-endpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctest-route-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]

  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.test.id]
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-custom-domain-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS10"
  }
}

resource "azurerm_cdn_frontdoor_custom_domain_association" "test" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.test.id
  cdn_frontdoor_route_ids        = [azurerm_cdn_frontdoor_route.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(8))
}
