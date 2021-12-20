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

type AfdCustomDomainResource struct{}

func TestAccCdnAfdCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := AfdCustomDomainResource{}

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

func TestAccCdnAfdCustomDomain_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := AfdCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r AfdCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.AFDCustomDomainsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving CDN Front Door Custom Domain %q (Resource Group %q / Profile Name %q): %+v", id.CustomDomainName, id.ResourceGroup, id.ProfileName, err)
	}
	return utils.Bool(true), nil
}

func (r AfdCustomDomainResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	customDomainsClient := client.Cdn.AFDCustomDomainsClient
	future, err := customDomainsClient.Delete(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		return nil, fmt.Errorf("deleting CDN Front Door custom domain %q (Resource Group %q / Profile %q): %+v", id.CustomDomainName, id.ResourceGroup, id.ProfileName, err)
	}
	if err := future.WaitForCompletionRef(ctx, customDomainsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of CDN Endpoint %q (Resource Group %q / Profile %q): %+v", id.CustomDomainName, id.ResourceGroup, id.ProfileName, err)
	}

	return utils.Bool(true), nil
}

func (r AfdCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name       = "acctestcdnprof%[1]d"
  profile_id = azurerm_cdn_frontdoor_profile.test.id
  host_name  = "custom.domain.tld"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name       = "acctestcdnend%[1]d"
  profile_id = azurerm_cdn_frontdoor_profile.test.id

  origin_response_timeout_in_seconds = 60
}
`, data.RandomInteger, data.Locations.Primary)
}
