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

type AfdOriginGroupResource struct{}

func TestAccCdnAfdOriginGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin_group", "test")
	r := AfdOriginGroupResource{}

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

func TestAccCdnAfdOriginGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin_group", "test")
	r := AfdOriginGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r AfdOriginGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdOriginGroupsID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.AFDOriginGroupsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving CDN Front Door Origin Group %q (Resource Group %q / Profile Name %q): %+v", id.OriginGroupName, id.ResourceGroup, id.ProfileName, err)
	}
	return utils.Bool(true), nil
}

func (r AfdOriginGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AfdOriginGroupsID(state.ID)
	if err != nil {
		return nil, err
	}

	originGroupsClient := client.Cdn.AFDOriginGroupsClient
	future, err := originGroupsClient.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
	if err != nil {
		return nil, fmt.Errorf("deleting CDN Front Door origins group %q (Resource Group %q / Profile %q): %+v", id.OriginGroupName, id.ResourceGroup, id.ProfileName, err)
	}
	if err := future.WaitForCompletionRef(ctx, originGroupsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of CDN Front Door origins group %q (Resource Group %q / Profile %q): %+v", id.OriginGroupName, id.ResourceGroup, id.ProfileName, err)
	}

	return utils.Bool(true), nil
}

func (r AfdOriginGroupResource) basic(data acceptance.TestData) string {
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

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name = "acctestcdnorigingroup%[1]d"

  profile_id = azurerm_cdn_frontdoor_profile.test.id

  health_probe {
    protocol            = "Http"
    path                = "/*"
    request_type        = "GET"
    interval_in_seconds = 240
  }

  load_balancing {
    sample_size                 = 6
    successful_samples_required = 3
    additional_latency_in_ms    = 10
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
