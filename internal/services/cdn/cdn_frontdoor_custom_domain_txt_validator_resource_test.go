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

type CdnFrontdoorCustomDomainTxtValidatorResource struct{}

func TestAccCdnFrontdoorCustomDomainTxtValidator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainTxtValidatorResource{}
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

func TestAccCdnFrontdoorCustomDomainTxtValidator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainTxtValidatorResource{}
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

func TestAccCdnFrontdoorCustomDomainTxtValidator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainTxtValidatorResource{}
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

func TestAccCdnFrontdoorCustomDomainTxtValidator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontdoorCustomDomainTxtValidatorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorCustomDomainsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-c-%d"
  frontdoor_cdn_profile_id = azurerm_cdn_frontdoor_profile.test.id

  azure_dns_zone_id                       = ""
  host_name                               = ""
  pre_validated_custom_domain_resource_id = ""

  tls {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_custom_domain" "import" {
  name                     = azurerm_cdn_frontdoor_custom_domain.test.name
  frontdoor_cdn_profile_id = azurerm_cdn_frontdoor_profile.test.id

  azure_dns_zone_id                       = ""
  host_name                               = ""
  pre_validated_custom_domain_resource_id = ""

  tls {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, config)
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-c-%d"
  frontdoor_cdn_profile_id = azurerm_cdn_frontdoor_profile.test.id

  azure_dns_zone_id                       = ""
  host_name                               = ""
  pre_validated_custom_domain_resource_id = ""

  tls {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorCustomDomainTxtValidatorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-c-%d"
  frontdoor_cdn_profile_id = azurerm_cdn_frontdoor_profile.test.id
  azure_dns_zone {
    id = ""
  }
  pre_validated_custom_domain_resource_id {
    id = ""
  }
  tls {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}
