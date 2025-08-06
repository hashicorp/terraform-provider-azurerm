// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontDoorResource struct{}

func TestAccFrontDoor_createShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	expectedError := frontdoor.CreateDeprecationMessage
	if frontdoor.IsFrontDoorFullyRetired() {
		expectedError = frontdoor.FullyRetiredMessage
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(expectedError),
		},
	})
}

func (FrontDoorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := frontdoors.ParseFrontDoorIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Frontdoor.FrontDoorsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (FrontDoorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
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
