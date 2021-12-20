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

type AfdEndpointResource struct{}

func TestAccCdnAfdEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_endpoint", "test")
	r := AfdEndpointResource{}

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

func TestAccCdnAfdEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_endpoint", "test")
	r := AfdEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccCdnAfdEndpoint_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_endpoint", "test")
	r := AfdEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AfdEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdEndpointsID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.AFDEndpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving CDN Front Door Endpoint %q (Resource Group %q / Profile Name %q): %+v", id.AfdEndpointName, id.ResourceGroup, id.ProfileName, err)
	}
	return utils.Bool(true), nil
}

func (r AfdEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdEndpointsID(state.ID)
	if err != nil {
		return nil, err
	}

	endpointsClient := client.Cdn.AFDEndpointsClient
	future, err := endpointsClient.Delete(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		return nil, fmt.Errorf("deleting CDN Endpoint %q (Resource Group %q / Profile %q): %+v", id.AfdEndpointName, id.ResourceGroup, id.ProfileName, err)
	}
	if err := future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of CDN Endpoint %q (Resource Group %q / Profile %q): %+v", id.AfdEndpointName, id.ResourceGroup, id.ProfileName, err)
	}

	return utils.Bool(true), nil
}

func (r AfdEndpointResource) basic(data acceptance.TestData) string {
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

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
  profile_id          = azurerm_cdn_frontdoor_profile.test.id

  origin_response_timeout_in_seconds = 60
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdEndpointResource) withTags(data acceptance.TestData) string {
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

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
  profile_id          = azurerm_cdn_frontdoor_profile.test.id

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdEndpointResource) withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                = "acctestcdnend%[1]d"
  profile_id          = azurerm_cdn_frontdoor_profile.test.id

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
