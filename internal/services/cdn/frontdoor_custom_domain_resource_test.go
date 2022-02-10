package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorCustomDomainResource struct{}

func TestAccFrontdoorCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_domain", "test")
	r := FrontdoorCustomDomainResource{}
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

func TestAccFrontdoorCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_domain", "test")
	r := FrontdoorCustomDomainResource{}
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

func TestAccFrontdoorCustomDomain_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_domain", "test")
	r := FrontdoorCustomDomainResource{}
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

func TestAccFrontdoorCustomDomain_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_domain", "test")
	r := FrontdoorCustomDomainResource{}
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

func (r FrontdoorCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := afdcustomdomains.ParseCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorCustomDomainsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-cdn-%d"
  location = "%s"
}
resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r FrontdoorCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_custom_domain" "test" {
  name           = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
  azure_dns_zone {
    id = ""
  }
  host_name = ""
  pre_validated_custom_domain_resource_id {
    id = ""
  }
  tls_settings {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorCustomDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_custom_domain" "import" {
  name           = azurerm_frontdoor_custom_domain.test.name
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
  azure_dns_zone {
    id = ""
  }
  host_name = ""
  pre_validated_custom_domain_resource_id {
    id = ""
  }
  tls_settings {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, config)
}

func (r FrontdoorCustomDomainResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_custom_domain" "test" {
  name           = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
  azure_dns_zone {
    id = ""
  }
  host_name = ""
  pre_validated_custom_domain_resource_id {
    id = ""
  }
  tls_settings {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorCustomDomainResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_custom_domain" "test" {
  name           = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
  azure_dns_zone {
    id = ""
  }
  pre_validated_custom_domain_resource_id {
    id = ""
  }
  tls_settings {
    certificate_type    = ""
    minimum_tls_version = ""
    secret {
      id = ""
    }
  }
}
`, template, data.RandomInteger)
}
