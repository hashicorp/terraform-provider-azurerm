package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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
				Config: testAccAzureRMCdnEndpoint_hostHeader(data, "www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", "www.example.com"),
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
					resource.TestCheckResourceAttr(data.ResourceName, "origin_host_header", "www.example.com"),
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

func testCheckAzureRMCdnEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, profileName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cdnEndpointsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CDN Endpoint %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMCdnEndpointDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		profileName := rs.Primary.Attributes["profile_name"]

		resp, err := conn.Get(ctx, resourceGroup, profileName, name)
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
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
  name                = "${azurerm_cdn_endpoint.test.name}"
  profile_name        = "${azurerm_cdn_endpoint.test.profile_name}"
  location            = "${azurerm_cdn_endpoint.test.location}"
  resource_group_name = "${azurerm_cdn_endpoint.test.resource_group_name}"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
    https_port = 443
    http_port  = 80
  }
}
`, template)
}

func testAccAzureRMCdnEndpoint_hostHeader(data acceptance.TestData, domain string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  origin_host_header  = "%s"

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.example.com"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.example.com"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  origin {
    name       = "acceptanceTestCdnOrigin2"
    host_name  = "www.example.com"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  is_http_allowed     = false
  is_https_allowed    = true
  origin_path         = "/origin-path"
  probe_path          = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  is_http_allowed     = false
  is_https_allowed    = true
  optimization_type   = "GeneralWebDelivery"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMCdnEndpoint_fullFields(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                          = "acctestcdnend%d"
  profile_name                  = "${azurerm_cdn_profile.test.name}"
  location                      = "${azurerm_resource_group.test.location}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  is_http_allowed               = true
  is_https_allowed              = true
  content_types_to_compress     = ["text/html"]
  is_compression_enabled        = true
  querystring_caching_behaviour = "UseQueryString"
  origin_host_header            = "www.example.com"
  optimization_type             = "GeneralWebDelivery"
  origin_path                   = "/origin-path"
  probe_path                    = "/origin-path/probe"

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  is_http_allowed     = %s
  is_https_allowed    = %s

  origin {
    name       = "acceptanceTestCdnOrigin1"
    host_name  = "www.example.com"
    https_port = 443
    http_port  = 80
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, isHttpAllowed, isHttpsAllowed)
}
