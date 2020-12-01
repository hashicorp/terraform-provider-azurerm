package frontdoor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FrontDoorCustomHttpsConfigurationResource struct {
}

func TestAccFrontDoorCustomHttpsConfiguration_CustomHttps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_https_configuration", "test")
	r := FrontDoorCustomHttpsConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.CustomHttpsEnabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_https_provisioning_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("custom_https_configuration.0.certificate_source").HasValue("FrontDoor"),
			),
		},
		{
			Config: r.CustomHttpsDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_https_provisioning_enabled").HasValue("false"),
			),
		},
	})
}

func (FrontDoorCustomHttpsConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FrontendEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Frontdoor.FrontDoorsFrontendClient.Get(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Frontend Endpoint %q (Front Door %q / Resource Group %q) does not exist", id.Name, id.FrontDoorName, id.ResourceGroup)
	}

	return utils.Bool(resp.FrontendEndpointProperties != nil), nil
}

func (r FrontDoorCustomHttpsConfigurationResource) CustomHttpsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoint[0].id
  resource_group_name               = azurerm_resource_group.test.name
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source = "FrontDoor"
  }
}
`, r.template(data))
}

func (r FrontDoorCustomHttpsConfigurationResource) CustomHttpsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoint[0].id
  resource_group_name               = azurerm_resource_group.test.name
  custom_https_provisioning_enabled = false
}
`, r.template(data))
}

func (FrontDoorCustomHttpsConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing-custom"
  endpoint_name       = "frontend-endpoint-custom"
  health_probe_name   = "health-probe-custom"
  load_balancing_name = "load-balancing-setting-custom"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctest-FD-%d"
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name      = local.endpoint_name
    host_name = "acctest-FD-%d.azurefd.net"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
