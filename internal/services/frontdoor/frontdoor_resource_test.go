// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontDoorResource struct{}

func TestAccFrontDoor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_pool_health_probe.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("backend_pool_health_probe.0.probe_method").HasValue("GET"),
			),
		},
		{
			Config: r.basicDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_pool_health_probe.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("backend_pool_health_probe.0.probe_method").HasValue("HEAD"),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}

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

func TestAccFrontDoor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_multiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePools(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_pool.#").HasValue("2"),
				check.That(data.ResourceName).Key("backend_pool_health_probe.#").HasValue("2"),
				check.That(data.ResourceName).Key("backend_pool_load_balancing.#").HasValue("2"),
				check.That(data.ResourceName).Key("routing_rule.#").HasValue("2"),
			),
		},
		// I am ignoring the returned values here because the ImportStep does not respect my state file reording code,
		// and the API is currently returning these value intermentently out of order causing false positive failures
		// of the test case.
		data.ImportStep("explicit_resource_order", "routing_rule", "backend_pool_load_balancing", "backend_pool_health_probe", "backend_pool"),
	})
}

func TestAccFrontDoor_routingRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routingRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_waf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.waf(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_RulesEngine(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.testRulesEngineConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
	})
}

func TestAccFrontDoor_EnableDisableCache(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	r := FrontDoorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.EnableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
		{
			Config: r.DisableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
		{
			Config: r.EnableCache(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
		{
			Config: r.testWithCachingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("explicit_resource_order"),
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

func (FrontDoorResource) basicDisabled(data acceptance.TestData) string {
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
    name         = local.health_probe_name
    enabled      = false
    probe_method = "HEAD"
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

func (r FrontDoorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor" "import" {
  name                = azurerm_frontdoor.test.name
  resource_group_name = azurerm_frontdoor.test.resource_group_name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = azurerm_frontdoor.test.backend_pool_settings.0.enforce_backend_pools_certificate_name_check
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
`, r.basic(data), data.RandomInteger)
}

func (FrontDoorResource) complete(data acceptance.TestData) string {
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
  name                  = "acctest-FD-%d"
  resource_group_name   = azurerm_resource_group.test.name
  friendly_name         = "TestGroup"
  load_balancer_enabled = false

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
    backend_pools_send_receive_timeout_seconds   = 45
  }


  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    enabled            = false
    forwarding_configuration {
      forwarding_protocol    = "MatchRequest"
      backend_pool_name      = local.backend_name
      custom_forwarding_path = "/"
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name                = local.health_probe_name
    interval_in_seconds = 30
    path                = "/"
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
      priority    = 2
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name                         = local.endpoint_name
    host_name                    = "acctest-FD-%d.azurefd.net"
    session_affinity_enabled     = true
    session_affinity_ttl_seconds = 2
  }
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (FrontDoorResource) routingRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%[1]d"
  location = "%[2]s"
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
}

resource "azurerm_frontdoor" "test" {
  name                = "acctest-FD-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    redirect_configuration {
      redirect_type       = "Moved"
      redirect_protocol   = "HttpOnly"
      custom_query_string = false
      custom_path         = "/"
      custom_host         = "127.0.0.1"
      custom_fragment     = "5fgdfg"
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
    host_name = "acctest-FD-%[1]d.azurefd.net"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FrontDoorResource) waf(data acceptance.TestData) string {
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

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                = "acctestwafp%d"
  resource_group_name = azurerm_resource_group.test.name
  mode                = "Prevention"
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
    name                                    = local.endpoint_name
    host_name                               = "acctest-FD-%d.azurefd.net"
    web_application_firewall_policy_link_id = azurerm_frontdoor_firewall_policy.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FrontDoorResource) DisableCache(data acceptance.TestData) string {
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

func (FrontDoorResource) EnableCache(data acceptance.TestData) string {
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
      cache_enabled       = true
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

func (FrontDoorResource) multiplePools(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%[1]d"
  location = "%s"
}

resource "azurerm_frontdoor" "test" {
  name                = "acctest-FD-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

  frontend_endpoint {
    name      = "acctest-FD-%[1]d-default-FE"
    host_name = "acctest-FD-%[1]d.azurefd.net"
  }

  # --- Pool 1

  routing_rule {
    name               = "acctest-FD-%[1]d-bing-RR"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/poolBing/*"]
    frontend_endpoints = ["acctest-FD-%[1]d-default-FE"]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "acctest-FD-%[1]d-pool-bing"
      cache_enabled       = true
    }
  }

  backend_pool_load_balancing {
    name                            = "acctest-FD-%[1]d-bing-LB"
    additional_latency_milliseconds = 0
    sample_size                     = 4
    successful_samples_required     = 2
  }

  backend_pool_health_probe {
    name         = "acctest-FD-%[1]d-bing-HP"
    protocol     = "Https"
    enabled      = true
    probe_method = "HEAD"
  }

  backend_pool {
    name                = "acctest-FD-%[1]d-pool-bing"
    load_balancing_name = "acctest-FD-%[1]d-bing-LB"
    health_probe_name   = "acctest-FD-%[1]d-bing-HP"

    backend {
      host_header = "bing.com"
      address     = "bing.com"
      http_port   = 80
      https_port  = 443
      weight      = 75
      enabled     = true
    }
  }

  # --- Pool 2

  routing_rule {
    name               = "acctest-FD-%[1]d-google-RR"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/poolGoogle/*"]
    frontend_endpoints = ["acctest-FD-%[1]d-default-FE"]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "acctest-FD-%[1]d-pool-google"
      cache_enabled       = true
    }
  }

  backend_pool_load_balancing {
    name                            = "acctest-FD-%[1]d-google-LB"
    additional_latency_milliseconds = 0
    sample_size                     = 4
    successful_samples_required     = 2
  }

  backend_pool_health_probe {
    name     = "acctest-FD-%[1]d-google-HP"
    protocol = "Https"
  }

  backend_pool {
    name                = "acctest-FD-%[1]d-pool-google"
    load_balancing_name = "acctest-FD-%[1]d-google-LB"
    health_probe_name   = "acctest-FD-%[1]d-google-HP"

    backend {
      host_header = "google.com"
      address     = "google.com"
      http_port   = 80
      https_port  = 443
      weight      = 75
      enabled     = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FrontDoorResource) testWithCachingEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%[1]d"
  location = "%s"
}

resource "azurerm_frontdoor" "test" {
  name                = "acctest-FD-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }


  frontend_endpoint {
    name      = "acctest-FD-%[1]d-default-FE"
    host_name = "acctest-FD-%[1]d.azurefd.net"
  }

  routing_rule {
    name               = "acctest-FD-%[1]d-bing-RR"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/poolBing/*"]
    frontend_endpoints = ["acctest-FD-%[1]d-default-FE"]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "acctest-FD-%[1]d-pool-bing"
      cache_enabled       = true

      cache_use_dynamic_compression = false

      cache_query_parameter_strip_directive = "StripAllExcept"

      cache_duration = "P90D"

      cache_query_parameters = [
        "width",
        "height"
      ]
    }
  }

  backend_pool_load_balancing {
    name                            = "acctest-FD-%[1]d-bing-LB"
    additional_latency_milliseconds = 0
    sample_size                     = 4
    successful_samples_required     = 2
  }

  backend_pool_health_probe {
    name         = "acctest-FD-%[1]d-bing-HP"
    protocol     = "Https"
    enabled      = true
    probe_method = "HEAD"
  }

  backend_pool {
    name                = "acctest-FD-%[1]d-pool-bing"
    load_balancing_name = "acctest-FD-%[1]d-bing-LB"
    health_probe_name   = "acctest-FD-%[1]d-bing-HP"

    backend {
      host_header = "bing.com"
      address     = "bing.com"
      http_port   = 80
      https_port  = 443
      weight      = 75
      enabled     = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FrontDoorResource) testRulesEngineConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%[1]d"
  location = "%s"
}

resource "azurerm_frontdoor" "test" {
  name                = "acctest-FD-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

  frontend_endpoint {
    name      = "acctest-FD-%[1]d-default-FE"
    host_name = "acctest-FD-%[1]d.azurefd.net"
  }

  routing_rule {
    name               = "acctest-FD-%[1]d-bing-RR"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/poolBing/*"]
    frontend_endpoints = ["acctest-FD-%[1]d-default-FE"]

    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "acctest-FD-%[1]d-pool-bing"
      cache_enabled       = false
    }
  }

  backend_pool_load_balancing {
    name                            = "acctest-FD-%[1]d-bing-LB"
    additional_latency_milliseconds = 0
    sample_size                     = 4
    successful_samples_required     = 2
  }

  backend_pool_health_probe {
    name         = "acctest-FD-%[1]d-bing-HP"
    protocol     = "Https"
    enabled      = true
    probe_method = "HEAD"
  }

  backend_pool {
    name                = "acctest-FD-%[1]d-pool-bing"
    load_balancing_name = "acctest-FD-%[1]d-bing-LB"
    health_probe_name   = "acctest-FD-%[1]d-bing-HP"

    backend {
      host_header = "bing.com"
      address     = "bing.com"
      http_port   = 80
      https_port  = 443
      weight      = 75
      enabled     = true
    }
  }
}

resource "azurerm_frontdoor_rules_engine" "sample_engine_config" {
  name                = "CORSAPI"
  frontdoor_name      = azurerm_frontdoor.test.name
  resource_group_name = azurerm_resource_group.test.name

  rule {
    name     = "debug"
    priority = 1

    match_condition {
      variable = "RequestMethod"
      operator = "Equal"
      value    = ["GET", "POST"]
    }

    action {
      response_header {
        header_action_type = "Append"
        header_name        = "X-TEST-HEADER"
        value              = "CORS API Rule"
      }
    }
  }

  rule {
    name     = "origin"
    priority = 2

    action {
      request_header {
        header_action_type = "Overwrite"
        header_name        = "Origin"
        value              = "*"
      }
      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Origin"
        value              = "*"
      }
      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Credentials"
        value              = "true"
      }
    }
  }

  depends_on = [
    azurerm_frontdoor.test
  ]
}
`, data.RandomInteger, data.Locations.Primary)
}
