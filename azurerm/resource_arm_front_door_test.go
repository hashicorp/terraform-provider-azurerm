package azurerm

import (
		"fmt"
		"strings"
    "testing"

    "github.com/hashicorp/terraform/helper/resource"
		"github.com/hashicorp/terraform/terraform"
		"github.com/hashicorp/terraform/helper/acctest"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFrontDoor_basic(t *testing.T) {
	resourceName := "azurerm_frontdoor.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	config := testAccAzureRMFrontDoor_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAccFrontDoor-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "friendly_name", "tafd"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "enforce_backend_pools_certificate_name_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.name", fmt.Sprintf("testAccBackendPool1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.backend.0.address", fmt.Sprintf("testaccfrontdoorsa%s.blob.core.windows.net", rs)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.backend.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.load_balancing_name", fmt.Sprintf("testAccLoadBalancingSettings1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.health_probe_name", fmt.Sprintf("testAccHealthProbeSetting1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.backend.0.http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.backend.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool.0.backend.0.weight", "50"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool_health_probe.0.name", fmt.Sprintf("testAccHealthProbeSetting1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool_health_probe.0.protocol", "Https"),
					resource.TestCheckResourceAttr(resourceName, "backend_pool_load_balancing.0.name", fmt.Sprintf("testAccLoadBalancingSettings1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "backend_pool_load_balancing.0.successful_samples_required", "2"),
					resource.TestCheckResourceAttr(resourceName, "frontend_endpoint.0.name", fmt.Sprintf("testAccFrontendEndpoint1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "frontend_endpoint.0.host_name", fmt.Sprintf("testAccFrontDoor-%d.azurefd.net", ri)),
					resource.TestCheckResourceAttr(resourceName, "frontend_endpoint.0.custom_https_configuration.0.certificate_source", "FrontDoor"),
					resource.TestCheckResourceAttr(resourceName, "frontend_endpoint.0.session_affinity_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "frontend_endpoint.0.session_affinity_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.name", fmt.Sprintf("testAccRoutingRulerule1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.accepted_protocols.0", "Http"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.accepted_protocols.1", "Https"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.forwarding_configuration.0.cache_use_dynamic_compression", "true"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.forwarding_configuration.0.forwarding_protocol", "MatchRequest"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.forwarding_configuration.0.cache_query_parameter_strip_directive", "StripNone"),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.frontend_endpoints.0", fmt.Sprintf("testAccFrontendEndpoint1-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.0.patterns_to_match.0", "/*"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMFrontDoorExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("Front Door not found: %s", resourceName)
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]

        client := testAccProvider.Meta().(*ArmClient).frontDoorsClient
        ctx := testAccProvider.Meta().(*ArmClient).StopContext

        if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
            if utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Front Door %q (Resource Group %q) does not exist", name, resourceGroup)
            }
            return fmt.Errorf("Bad: Get on frontDoorsClient: %+v", err)
        }

        return nil
    }
}

func testCheckAzureRMFrontDoorDestroy(s *terraform.State) error {
    client := testAccProvider.Meta().(*ArmClient).frontDoorsClient
    ctx := testAccProvider.Meta().(*ArmClient).StopContext

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "azurerm_front_door" {
            continue
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]

        if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Get on frontDoorsClient: %+v", err)
            }
        }

        return nil
    }

    return nil
}

func testAccAzureRMFrontDoor_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccfrontdoorsa%[2]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_frontdoor" "test" {
  name                      = "testAccFrontDoor-%[1]d"
  friendly_name             = "tafd"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  enabled                   = true

  routing_rule {
      name                  = "testAccRoutingRulerule1-%[1]d"
      enabled               = true
      accepted_protocols    = ["Http", "Https"]
      patterns_to_match     = ["/*"]
      frontend_endpoints    = ["testAccFrontendEndpoint1-%[1]d"]
      forwarding_configuration {
          forwarding_protocol = "MatchRequest"
          cache_use_dynamic_compression = true
          backend_pool_name   = "testAccBackendPool1-%[1]d"
      }
  }

  enforce_backend_pools_certificate_name_check = true

  backend_pool_load_balancing {
    name = "testAccLoadBalancingSettings1-%[1]d"
    sample_size = 4
    successful_samples_required = 2
    additional_latency_milliseconds = 0
  }

  backend_pool_health_probe {
    name = "testAccHealthProbeSetting1-%[1]d"
    path = "/"
    protocol = "Https"
    interval_in_seconds = 30
  }

  backend_pool {
      name = "testAccBackendPool1-%[1]d"
      backend {
          enabled   = true
          host_header = "${azurerm_storage_account.test.name}.blob.core.windows.net"
          address   = "${azurerm_storage_account.test.name}.blob.core.windows.net"
          http_port  = 80
          https_port = 443
          weight    = 50
          priority  = 1
      }

      load_balancing_name = "testAccLoadBalancingSettings1-%[1]d"
      health_probe_name   = "testAccHealthProbeSetting1-%[1]d"
  }

  frontend_endpoint {
    name = "testAccFrontendEndpoint1-%[1]d"
    host_name = "testAccFrontDoor-%[1]d.azurefd.net"
    session_affinity_enabled = true
    session_affinity_ttl_seconds = 0
    custom_https_configuration {
        certificate_source = "FrontDoor"
    }
  }
}
`, rInt, rString, location)
}