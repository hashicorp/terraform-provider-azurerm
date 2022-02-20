package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorRouteResource struct{}

func TestAccFrontdoorRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_route", "test")
	r := FrontdoorRouteResource{}
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

func TestAccFrontdoorRoute_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_route", "test")
	r := FrontdoorRouteResource{}
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

func TestAccFrontdoorRoute_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_route", "test")
	r := FrontdoorRouteResource{}
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

func TestAccFrontdoorRoute_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_route", "test")
	r := FrontdoorRouteResource{}
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

func (r FrontdoorRouteResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := routes.ParseRouteID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorRoutesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorRouteResource) template(data acceptance.TestData) string {
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

resource "azurerm_frontdoor_endpoint" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r FrontdoorRouteResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_route" "test" {
  name                  = "acctest-c-%d"
  frontdoor_endpoint_id = azurerm_frontdoor_endpoint.test.id

  cache_configuration {
    query_parameters              = ["foo", "bar"]
    query_string_caching_behavior = "IgnoreQueryString"
  }

  custom_domains {
    id = ""
  }

  enabled                = true
  forwarding_protocol    = true
  https_redirect         = true
  link_to_default_domain = true
  origin_group_id        = azurerm_frontdoor_group.test.id
  origin_path            = ""
  patterns_to_match      = []
  rule_set_ids           = [""]
  supported_protocols    = ["Http", "Https"]
}
`, template, data.RandomInteger)
}

func (r FrontdoorRouteResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_route" "import" {
  name                  = azurerm_frontdoor_route.test.name
  frontdoor_endpoint_id = azurerm_frontdoor_endpoint.test.id

  cache_configuration {
    query_parameters              = ""
    query_string_caching_behavior = ""
  }

  custom_domains {
    id = ""
  }

  enabled_state          = ""
  forwarding_protocol    = ""
  https_redirect         = ""
  link_to_default_domain = ""
  origin_group_id        = ""
  origin_path            = ""
  patterns_to_match      = []
  rule_set_ids           = [""]
  supported_protocols    = ["Http", "Https"]
}
`, config)
}

func (r FrontdoorRouteResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_route" "test" {
  name                  = "acctest-c-%d"
  frontdoor_endpoint_id = azurerm_frontdoor_endpoint.test.id

  cache_configuration {
    query_parameters              = ""
    query_string_caching_behavior = ""
  }

  custom_domains {
    id = ""
  }

  enabled_state          = ""
  forwarding_protocol    = ""
  https_redirect         = ""
  link_to_default_domain = ""
  origin_group_id        = ""
  origin_path            = ""
  patterns_to_match      = [""]
  rule_set_ids           = [""]
  supported_protocols    = ["Http", "Https"]
}
`, template, data.RandomInteger)
}

func (r FrontdoorRouteResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_route" "test" {
  name                  = "acctest-c-%d"
  frontdoor_endpoint_id = azurerm_frontdoor_endpoint.test.id

  cache_configuration {
    query_parameters              = ""
    query_string_caching_behavior = ""
  }

  custom_domains {
    id = ""
  }

  enabled_state          = ""
  forwarding_protocol    = ""
  https_redirect         = ""
  link_to_default_domain = ""
  origin_group_id        = ""
  origin_path            = ""
  patterns_to_match      = []
  rule_set_ids           = [""]
  supported_protocols    = ["Https"]
}
`, template, data.RandomInteger)
}
