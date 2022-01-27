package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorProfileOriginResource struct{}

func TestAccFrontdoorProfileOrigin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_origin", "test")
	r := FrontdoorProfileOriginResource{}
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

func TestAccFrontdoorProfileOrigin_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_origin", "test")
	r := FrontdoorProfileOriginResource{}
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

func TestAccFrontdoorProfileOrigin_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_origin", "test")
	r := FrontdoorProfileOriginResource{}
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

func TestAccFrontdoorProfileOrigin_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_origin", "test")
	r := FrontdoorProfileOriginResource{}
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

func (r FrontdoorProfileOriginResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := afdorigins.ParseOriginGroupOriginID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorProfileOriginsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorProfileOriginResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-cdn-%d"
  location = "%s"
}
resource "azurerm_cdn_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_frontdoor_profile_origin_group" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_cdn_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r FrontdoorProfileOriginResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_profile_origin" "test" {
  name                    = "acctest-c-%d"
  cdn_afd_origin_group_id = azurerm_frontdoor_profile_origin_group.test.id
  azure_origin {
    id = ""
  }
  enabled_state                  = ""
  enforce_certificate_name_check = false
  host_name                      = ""
  http_port                      = 0
  https_port                     = 0
  origin_host_header             = ""
  priority                       = 0
  weight                         = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileOriginResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_origin" "import" {
  name                    = azurerm_cdn_afd_origin.test.name
  cdn_afd_origin_group_id = azurerm_frontdoor_profile_origin_group.test.id
  azure_origin {
    id = ""
  }
  enabled_state                  = ""
  enforce_certificate_name_check = false
  host_name                      = ""
  http_port                      = 0
  https_port                     = 0
  origin_host_header             = ""
  priority                       = 0
  weight                         = 0
}
`, config)
}

func (r FrontdoorProfileOriginResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_origin" "test" {
  name                    = "acctest-c-%d"
  cdn_afd_origin_group_id = azurerm_frontdoor_profile_origin_group.test.id
  azure_origin {
    id = ""
  }
  enabled_state                  = ""
  enforce_certificate_name_check = false
  host_name                      = ""
  http_port                      = 0
  https_port                     = 0
  origin_host_header             = ""
  priority                       = 0
  weight                         = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileOriginResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_origin" "test" {
  name                    = "acctest-c-%d"
  cdn_afd_origin_group_id = azurerm_frontdoor_profile_origin_group.test.id
  azure_origin {
    id = ""
  }
  enabled_state                  = ""
  enforce_certificate_name_check = false
  host_name                      = ""
  http_port                      = 0
  https_port                     = 0
  origin_host_header             = ""
  priority                       = 0
  weight                         = 0
}
`, template, data.RandomInteger)
}
