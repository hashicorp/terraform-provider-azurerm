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

type CdnFrontdoorRouteResource struct{}

func TestAccCdnFrontdoorRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_route", "test")
	r := CdnFrontdoorRouteResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cdn_frontdoor_origin_group_id", "cdn_frontdoor_origin_ids"),
		{
			// You must delete the route first to disassociate the endpoint from the origin and origin group
			Config: r.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccCdnFrontdoorRoute_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_route", "test")
	r := CdnFrontdoorRouteResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
		{
			// You must delete the route first to disassociate the endpoint from the origin and origin group
			Config: r.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccCdnFrontdoorRoute_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_route", "test")
	r := CdnFrontdoorRouteResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cdn_frontdoor_origin_group_id", "cdn_frontdoor_origin_ids"),
		{
			// You must delete the route first to disassociate the endpoint from the origin and origin group
			Config: r.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccCdnFrontdoorRoute_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_route", "test")
	r := CdnFrontdoorRouteResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cdn_frontdoor_origin_group_id", "cdn_frontdoor_origin_ids"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cdn_frontdoor_origin_group_id", "cdn_frontdoor_origin_ids"),
		{
			// You must delete the route first to disassociate the endpoint from the origin and origin group
			Config: r.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (r CdnFrontdoorRouteResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorRouteID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorRoutesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontdoorRouteResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "accTestOriginGroup-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "accTestOrigin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id

  health_probes_enabled          = true
  enforce_certificate_name_check = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "accTestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "accTestRuleSet%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorRouteResource) destroy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "accTestOriginGroup-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "accTestOrigin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id

  health_probes_enabled          = true
  enforce_certificate_name_check = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "accTestEndpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "accTestRuleSet%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorRouteResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "accTestRoute-%d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]

  patterns_to_match   = ["/*"]
  supported_protocols = ["Http", "Https"]
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorRouteResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_route" "import" {
  name                          = azurerm_cdn_frontdoor_route.test.name
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]

  patterns_to_match   = ["/*"]
  supported_protocols = ["Http", "Https"]
}
`, config)
}

func (r CdnFrontdoorRouteResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "accTestRoute-%d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]

  enabled             = true
  forwarding_protocol = "HttpsOnly"
  https_redirect      = true
  # Cannot set this value because Frontdoor RP validates that the path is reachable
  # cdn_frontdoor_origin_path  = "contoso.com/site/content"
  patterns_to_match          = ["/*"]
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.test.id]
  supported_protocols        = ["Http", "Https"]

  cache_configuration {
    query_strings                 = ["foo", "bar"]
    query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorRouteResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "accTestRoute-%d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]

  enabled                = true
  forwarding_protocol    = "HttpOnly"
  https_redirect         = false
  # Cannot set this value because Frontdoor RP validates that the path is reachable
  # cdn_frontdoor_origin_path  = "contoso.com/site/content"
  patterns_to_match          = ["/*"]
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.test.id]
  supported_protocols        = ["Https"]

  cache_configuration {
    query_strings                 = ["bar"]
    query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
  }
}
`, template, data.RandomInteger)
}
