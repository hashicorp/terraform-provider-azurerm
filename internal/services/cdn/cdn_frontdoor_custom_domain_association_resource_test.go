// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorCustomDomainAssociationResource struct{}

// NOTE: There isn't a complete test case because the basic and the
// update together equals what the complete test case would be...
func TestAccCdnFrontDoorCustomDomainAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain_association", "test")
	r := CdnFrontDoorCustomDomainAssociationResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

// NOTE: the 'requiresImport' test is not possible on this resource

func TestAccCdnFrontDoorCustomDomainAssociation_removeAssociation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain_association", "test")
	r := CdnFrontDoorCustomDomainAssociationResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.remove(data),
			Check:              acceptance.ComposeTestCheckFunc(),
			ExpectNonEmptyPlan: true, // since deleting this resource actually removes the linked custom domain from the route resource(s)
		},
	})
}

func TestAccCdnFrontDoorCustomDomainAssociation_removeAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain_association", "test")
	r := CdnFrontDoorCustomDomainAssociationResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.remove(data),
			Check:              acceptance.ComposeTestCheckFunc(),
			ExpectNonEmptyPlan: true, // since deleting this resource actually removes the linked custom domain from the route resource(s)
		},
	})
}

func (r CdnFrontDoorCustomDomainAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorCustomDomainAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.AFDCustomDomainsClient
	customDomainId := afdcustomdomains.NewCustomDomainID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AssociationName)
	resp, err := client.Get(ctx, customDomainId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (r CdnFrontDoorCustomDomainAssociationResource) preCheck(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME` must be set for acceptance tests!")
	}
	if os.Getenv("ARM_TEST_DNS_ZONE_NAME") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_NAME` must be set for acceptance tests!")
	}
}

func (r CdnFrontDoorCustomDomainAssociationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain_association" "test" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.contoso.id
  cdn_frontdoor_route_ids        = [azurerm_cdn_frontdoor_route.contoso.id]
}
`, template)
}

func (r CdnFrontDoorCustomDomainAssociationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_route" "fabrikam" {
  name                          = "acctest-fabrikam-%[2]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  enabled                       = true

  https_redirect_enabled = true
  forwarding_protocol    = "HttpsOnly"
  patterns_to_match      = ["/sub-%[3]s"]
  supported_protocols    = ["Http", "Https"]

  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.contoso.id]
  link_to_default_domain          = true

  cache {
    compression_enabled           = true
    content_types_to_compress     = ["text/html", "text/javascript", "text/xml"]
    query_strings                 = ["account", "settings", "foo", "bar"]
    query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
  }
}

resource "azurerm_cdn_frontdoor_custom_domain_association" "test" {
  cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.contoso.id
  cdn_frontdoor_route_ids        = [azurerm_cdn_frontdoor_route.contoso.id, azurerm_cdn_frontdoor_route.fabrikam.id]
}
`, template, data.RandomInteger, data.RandomStringOfLength(10))
}

func (r CdnFrontDoorCustomDomainAssociationResource) remove(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
`, template)
}

func (r CdnFrontDoorCustomDomainAssociationResource) template(data acceptance.TestData) string {
	dnsZoneName := os.Getenv("ARM_TEST_DNS_ZONE_NAME")
	dnsZoneRG := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

data "azurerm_dns_zone" "test" {
  name                = "%[4]s"
  resource_group_name = "%[5]s"
}

locals {
  # Create a delegated child zone inside the test RG.
  # NOTE: ARM_TEST_DNS_ZONE_NAME / ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME must refer to a real, delegated parent zone.
  child_zone_label = "acctest%[1]d"
  child_zone_name  = join(".", [local.child_zone_label, data.azurerm_dns_zone.test.name])
}

resource "azurerm_dns_zone" "child" {
  name                = local.child_zone_name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "delegation" {
  name                = local.child_zone_label
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300

  records = azurerm_dns_zone.child.name_servers
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-profile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctest-origin-group-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  session_affinity_enabled = true

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "GET"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-origin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = true
  host_name                      = azurerm_storage_account.test.primary_blob_host
  origin_host_header             = azurerm_storage_account.test.primary_blob_host
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-endpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  enabled                  = true
}

resource "azurerm_cdn_frontdoor_route" "contoso" {
  name                          = "acctest-contoso-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  enabled                       = true

  depends_on = [azurerm_dns_txt_record.validation]

  https_redirect_enabled = true
  forwarding_protocol    = "HttpsOnly"
  patterns_to_match      = ["/%[3]s"]
  supported_protocols    = ["Http", "Https"]

  cdn_frontdoor_custom_domain_ids = [azurerm_cdn_frontdoor_custom_domain.contoso.id]
  link_to_default_domain          = false

  cache {
    compression_enabled           = true
    content_types_to_compress     = ["text/html", "text/javascript", "text/xml"]
    query_strings                 = ["account", "settings", "foo", "bar"]
    query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
  }
}

resource "azurerm_cdn_frontdoor_custom_domain" "contoso" {
  name                     = "acctest-contoso-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  depends_on = [azurerm_dns_ns_record.delegation]

  dns_zone_id = azurerm_dns_zone.child.id
  host_name   = join(".", ["%[3]s", azurerm_dns_zone.child.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_dns_txt_record" "validation" {
  depends_on = [azurerm_dns_ns_record.delegation]

  name                = join(".", ["_dnsauth", split(".", azurerm_cdn_frontdoor_custom_domain.contoso.host_name)[0]])
  zone_name           = azurerm_dns_zone.child.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300

  record {
    value = azurerm_cdn_frontdoor_custom_domain.contoso.validation_token
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(10), dnsZoneName, dnsZoneRG)
}
