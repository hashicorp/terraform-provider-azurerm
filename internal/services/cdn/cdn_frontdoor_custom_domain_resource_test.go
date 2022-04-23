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

type CdnFrontdoorCustomDomainResource struct {
	DoNotRunFrontdooCustomDomainTests bool
}

func TestAccCdnFrontdoorCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainResource{DoNotRunFrontdooCustomDomainTests: true}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontdoorCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainResource{DoNotRunFrontdooCustomDomainTests: true}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontdoorCustomDomain_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainResource{DoNotRunFrontdooCustomDomainTests: true}
	r.preCheck(t)

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

// TODO: Due to the validation logic in the service you cannot update the custom domain until
// it has been approved. Need to add a txt validator to facilitate testing the update functionality.

func (r CdnFrontdoorCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorCustomDomainsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontdoorCustomDomainResource) preCheck(t *testing.T) {
	if r.DoNotRunFrontdooCustomDomainTests {
		t.Skipf("`azurerm_cdn_frontdoor_custom_domain` currently is not testable due to service requirements")
	}
}

func (r CdnFrontdoorCustomDomainResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "accTestCustomDomain-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dns_zone_id = azurerm_dns_zone.test.id
  host_name   = join(".", ["fabrikam", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_custom_domain" "import" {
  name                     = "accTestCustomDomain-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dns_zone_id = azurerm_dns_zone.test.id
  host_name   = join(".", ["fabrikam", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}
`, config, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "accTestCustomDomain-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dns_zone_id = azurerm_dns_zone.test.id
  host_name   = join(".", ["fabrikam", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS10"
  }
}
`, template, data.RandomInteger)
}

// TODO: Add test case that uses pre_validated_custom_domain_resource_id
// TODO: Add test case that uses CMK
