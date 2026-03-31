// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FrontDoorCustomHttpsConfigurationResource struct{}

func TestAccFrontDoorCustomHttpsConfiguration_createShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_https_configuration", "test")
	r := FrontDoorCustomHttpsConfigurationResource{}

	expectedError := frontdoor.CreateDeprecationMessage
	if frontdoor.IsFrontDoorFullyRetired() {
		expectedError = frontdoor.FullyRetiredMessage
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.Enabled(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(expectedError),
		},
	})
}

func (FrontDoorCustomHttpsConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomHttpsConfigurationIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	frontdoorId, err := frontdoors.ParseFrontendEndpointID(state.Attributes["frontend_endpoint_id"])
	if err != nil {
		return nil, err
	}

	resp, err := clients.Frontdoor.FrontDoorsClient.FrontendEndpointsGet(ctx, *frontdoorId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r FrontDoorCustomHttpsConfigurationResource) Enabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoints[local.endpoint_name]
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source = "FrontDoor"
  }
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
  name                = "acctest-FD-%d"
  resource_group_name = azurerm_resource_group.test.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

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
