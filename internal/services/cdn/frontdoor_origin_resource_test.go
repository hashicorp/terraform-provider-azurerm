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

type FrontdoorOriginResource struct{}

func TestAccFrontdoorOrigin_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin", "test")
	r := FrontdoorOriginResource{}
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

func TestAccFrontdoorOrigin_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin", "test")
	r := FrontdoorOriginResource{}
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

func TestAccFrontdoorOrigin_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin", "test")
	r := FrontdoorOriginResource{}
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

func TestAccFrontdoorOrigin_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin", "test")
	r := FrontdoorOriginResource{}
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

func (r FrontdoorOriginResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorOriginID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorOriginsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r FrontdoorOriginResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-afdx-%d"
  location = "%s"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r FrontdoorOriginResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_origin" "test" {
  name                      = "acctest-c-%d"
  frontdoor_origin_group_id = azurerm_frontdoor_origin_group.test.id
  azure_origin_id           = ""

  enable_health_probes           = true
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

func (r FrontdoorOriginResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin" "import" {
  name                      = azurerm_cdn_afd_origin.test.name
  frontdoor_origin_group_id = azurerm_frontdoor_origin_group.test.id
  azure_origin_id           = ""

  enable_health_probes           = true
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

func (r FrontdoorOriginResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin" "test" {
  name                      = "acctest-c-%d"
  frontdoor_origin_group_id = azurerm_frontdoor_origin_group.test.id
  azure_origin_id           = ""

  enable_health_probes           = true
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

func (r FrontdoorOriginResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin" "test" {
  name                      = "acctest-c-%d"
  frontdoor_origin_group_id = azurerm_frontdoor_origin_group.test.id
  azure_origin_id           = ""

  enable_health_probes           = true
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
