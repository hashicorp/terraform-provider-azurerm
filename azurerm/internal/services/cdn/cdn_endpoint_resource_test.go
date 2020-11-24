package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
)

func TestAccAzureRMCdnEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMCdnEndpoint_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_cdn_endpoint"),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					testCheckAzureRMCdnEndpointDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_updateHostHeader(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_hostHeader(data, "www.contoso.com"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", "www.contoso.com"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_hostHeader(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", ""),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_hostHeader(data, "www.example2.com"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", "www.example2.com"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnEndpoint_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpoint_withoutCompression(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_withoutCompression(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "is_compression_enabled"),
				),
			}, data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnEndpoint_optimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_optimized(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "optimization_type", "GeneralWebDelivery"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_withGeoFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_geoFilters(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_filter.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_fullFields(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_fullFields(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_http_allowed", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_https_allowed", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_path", "/origin-path"),
					resource.TestCheckResourceAttr(data.ResourceName, "probe_path", "/origin-path/probe"),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", "www.contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "optimization_type", "GeneralWebDelivery"),
					resource.TestCheckResourceAttr(data.ResourceName, "querystring_caching_behaviour", "UseQueryString"),
					resource.TestCheckResourceAttr(data.ResourceName, "content_types_to_compress.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_compression_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_filter.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_isHttpAndHttpsAllowedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_isHttpAndHttpsAllowed(data, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_http_allowed", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_https_allowed", "false"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_isHttpAndHttpsAllowed(data, "false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_http_allowed", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_https_allowed", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_globalDeliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_globalDeliveryRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.0.behavior", "Override"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.0.duration", "5.04:44:23"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_globalDeliveryRuleUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.0.behavior", "SetIfMissing"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.cache_expiration_action.0.duration", "12.04:11:22"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.modify_response_header_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.modify_response_header_action.0.action", "Overwrite"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.modify_response_header_action.0.name", "Content-Type"),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.0.modify_response_header_action.0.value", "application/json"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_globalDeliveryRuleRemove(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "global_delivery_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_deliveryRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_deliveryRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.name", "http2https"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.order", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.redirect_type", "Found"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.protocol", "Https"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_deliveryRuleUpdate1(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.name", "http2https"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.order", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.0.negate_condition", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.redirect_type", "Found"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.protocol", "Https"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_deliveryRuleUpdate2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.name", "http2https"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.order", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.0.negate_condition", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.request_scheme_condition.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.redirect_type", "Found"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.0.url_redirect_action.0.protocol", "Https"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.name", "test"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.order", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.device_condition.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.device_condition.0.match_values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.modify_response_header_action.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.modify_response_header_action.0.action", "Delete"),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.1.modify_response_header_action.0.name", "Content-Language"),
				),
			},
			{
				Config: testAccAzureRMCdnEndpoint_deliveryRuleRemove(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "delivery_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnEndpoint_dnsAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnEndpoint_dnsAlias(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCdnEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.EndpointID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cdnEndpointsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CDN Endpoint %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMCdnEndpointDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for cdn endpoint: %s", name)
		}

		future, err := conn.Delete(ctx, resourceGroup, profileName, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on cdnEndpointsClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, conn.Client); err != nil {
			return fmt.Errorf("Bad: Delete on cdnEndpointsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMCdnEndpointDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.EndpointsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cdn_endpoint" {
			continue
		}

		id, err := parse.EndpointID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("CDN Endpoint still exists:\n%#v", resp.EndpointProperties)
		}
	}

	return nil
}

func testAccAzureRMCdnEndpoint_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMCdnEndpoint_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint" "import" {
  name                = azurerm_cdn_endpoint.test.name
  profile_name        = azurerm_cdn_endpoint.test.profile_name
  location            = azurerm_cdn_endpoint.test.location
  resource_group_name = azurerm_cdn_endpoint.test.resource_group_name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, template)
}

func testAccAzureRMCdnEndpoint_hostHeader(data acceptance.TestData, domain string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  origin_host_header  = "%s"

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, domain)
}

func testAccAzureRMCdnEndpoint_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_geoFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  origin_path         = "/origin-path"
  probe_path          = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  geo_filter {
    relative_path = "/some-example-endpoint"
    action        = "Allow"
    country_codes = ["GB"]
  }

  geo_filter {
    relative_path = "/some-other-endpoint"
    action        = "Block"
    country_codes = ["US"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_optimized(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  optimization_type   = "GeneralWebDelivery"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_withoutCompression(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = false
  is_https_allowed    = true
  optimization_type   = "GeneralWebDelivery"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_fullFields(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                          = "acctestcdnend%d"
  profile_name                  = azurerm_cdn_profile.test.name
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  is_http_allowed               = true
  is_https_allowed              = true
  content_types_to_compress     = ["text/html"]
  is_compression_enabled        = true
  querystring_caching_behaviour = "UseQueryString"
  origin_host_header            = "www.contoso.com"
  optimization_type             = "GeneralWebDelivery"
  origin_path                   = "/origin-path"
  probe_path                    = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  geo_filter {
    relative_path = "/some-example-endpoint"
    action        = "Allow"
    country_codes = ["GB"]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_isHttpAndHttpsAllowed(data acceptance.TestData, isHttpAllowed string, isHttpsAllowed string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_http_allowed     = %s
  is_https_allowed    = %s

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, isHttpAllowed, isHttpsAllowed)
}

func testAccAzureRMCdnEndpoint_globalDeliveryRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  global_delivery_rule {
    cache_expiration_action {
      behavior = "Override"
      duration = "5.04:44:23"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_globalDeliveryRuleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  global_delivery_rule {
    cache_expiration_action {
      behavior = "SetIfMissing"
      duration = "12.04:11:22"
    }

    modify_response_header_action {
      action = "Overwrite"
      name   = "Content-Type"
      value  = "application/json"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_globalDeliveryRuleRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_deliveryRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      match_values = ["HTTP"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_deliveryRuleUpdate1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      negate_condition = true
      match_values     = ["HTTPS"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_deliveryRuleUpdate2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }

  delivery_rule {
    name  = "http2https"
    order = 1

    request_scheme_condition {
      negate_condition = true
      match_values     = ["HTTPS"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }

  delivery_rule {
    name  = "test"
    order = 2

    device_condition {
      match_values = ["Mobile"]
    }

    modify_response_header_action {
      action = "Delete"
      name   = "Content-Language"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_deliveryRuleRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin_host_header = "www.contoso.com"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_dnsAlias(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestcdnep%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnep%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnep%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.contoso.com"
    https_port = 443
    http_port  = 80
  }
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_cdn_endpoint.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
