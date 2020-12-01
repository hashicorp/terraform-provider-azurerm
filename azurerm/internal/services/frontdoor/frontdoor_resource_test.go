package frontdoor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFrontDoor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_health_probe.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_health_probe.0.probe_method", "GET"),
				),
			},
			{
				Config: testAccAzureRMFrontDoor_basicDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_health_probe.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_health_probe.0.probe_method", "HEAD"),
				),
			},
			data.ImportStep(),
		},
	})
}

// remove in 3.0
func TestAccAzureRMFrontDoor_global(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_global(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "location", "global"),
				),
				ExpectNonEmptyPlan: true,
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFrontDoor_requiresImport),
		},
	})
}

func TestAccAzureRMFrontDoor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFrontDoor_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_multiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_multiplePools(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_health_probe.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_pool_load_balancing.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_waf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_waf(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_EnableDisableCache(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_EnableCache(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_use_dynamic_compression", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_query_parameter_strip_directive", "StripAll"),
				),
			},
			{
				Config: testAccAzureRMFrontDoor_DisableCache(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_use_dynamic_compression", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_query_parameter_strip_directive", "StripAll"),
				),
			},
			{
				Config: testAccAzureRMFrontDoor_EnableCache(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_use_dynamic_compression", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_rule.0.forwarding_configuration.0.cache_query_parameter_strip_directive", "StripAll"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_CustomHttps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_CustomHttpsEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_provisioning_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_configuration.0.certificate_source", "FrontDoor"),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_configuration.0.minimum_tls_version", "1.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_configuration.0.provisioning_state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_configuration.0.provisioning_substate", "CertificateDeployed"),
				),
			},
			{
				Config: testAccAzureRMFrontDoor_CustomHttpsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "frontend_endpoint.0.custom_https_provisioning_enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFrontDoorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Front Door not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Front Door %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on FrontDoorsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFrontDoorDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_frontdoor" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on FrontDoorsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMFrontDoor_basic(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_basicDisabled(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

// remove in 3.0
func testAccAzureRMFrontDoor_global(data acceptance.TestData) string {
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
  name                                         = "acctest-FD-%d"
  resource_group_name                          = azurerm_resource_group.test.name
  location                                     = "%s"
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMFrontDoor_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFrontDoor_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor" "import" {
  name                                         = azurerm_frontdoor.test.name
  resource_group_name                          = azurerm_frontdoor.test.resource_group_name
  enforce_backend_pools_certificate_name_check = azurerm_frontdoor.test.enforce_backend_pools_certificate_name_check

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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMFrontDoor_complete(data acceptance.TestData) string {
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
  name                                         = "acctest-FD-%d"
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false
  backend_pools_send_receive_timeout_seconds   = 45

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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_waf(data acceptance.TestData) string {
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
    name                                    = local.endpoint_name
    host_name                               = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled       = false
    web_application_firewall_policy_link_id = azurerm_frontdoor_firewall_policy.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_DisableCache(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_EnableCache(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_CustomHttpsEnabled(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = true
    custom_https_configuration {
      certificate_source = "FrontDoor"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_CustomHttpsDisabled(data acceptance.TestData) string {
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
    name                              = local.endpoint_name
    host_name                         = "acctest-FD-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_multiplePools(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-frontdoor-%[1]d"
  location = "%s"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctest-FD-%[1]d"
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  frontend_endpoint {
    name                              = "acctest-FD-%[1]d-default-FE"
    host_name                         = "acctest-FD-%[1]d.azurefd.net"
    custom_https_provisioning_enabled = false
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
