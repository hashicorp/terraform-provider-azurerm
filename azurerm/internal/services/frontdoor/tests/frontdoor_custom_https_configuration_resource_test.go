package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFrontDoorCustomHttpsConfiguration_CustomHttps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_custom_https_configuration", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoorCustomHttpsConfiguration_CustomHttpsEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorCustomHttpsConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_https_provisioning_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_https_configuration.0.certificate_source", "FrontDoor"),
				),
			},
			{
				Config: testAccAzureRMFrontDoorCustomHttpsConfiguration_CustomHttpsDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorCustomHttpsConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_https_provisioning_enabled", "false"),
				),
			},
		},
	})
}

func testCheckAzureRMFrontDoorCustomHttpsConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsFrontendClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Front Door Custom Https Configuration not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["frontend_endpoint_id"])
		if err != nil {
			return fmt.Errorf("Bad: cannot parse frontend_endpoint_id for %q", resourceName)
		}
		frontDoorName := id.Path["frontdoors"]
		// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
		if frontDoorName == "" {
			frontDoorName = id.Path["Frontdoors"]
		}
		frontendEndpointName := id.Path["frontendendpoints"]
		if frontendEndpointName == "" {
			frontDoorName = id.Path["FrontendEndpoints"]
		}

		resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Front Door (%q) Frontend Endpoint %q (Resource Group %q) does not exist", frontDoorName, frontendEndpointName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on FrontDoorsFrontendClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMFrontDoorCustomHttpsConfiguration_CustomHttpsEnabled(data acceptance.TestData) string {
	template := testAccAzureRMFrontDoorCustomHttpsConfiguration_template(data)
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
`, template)
}

func testAccAzureRMFrontDoorCustomHttpsConfiguration_CustomHttpsDisabled(data acceptance.TestData) string {
	template := testAccAzureRMFrontDoorCustomHttpsConfiguration_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_custom_https_configuration" "test" {
  frontend_endpoint_id              = azurerm_frontdoor.test.frontend_endpoint[0].id
  resource_group_name               = azurerm_resource_group.test.name
  custom_https_provisioning_enabled = false
}
`, template)
}

func testAccAzureRMFrontDoorCustomHttpsConfiguration_template(data acceptance.TestData) string {
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
