package tests

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMTrafficManagerEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	azureResourceName := "azurerm_traffic_manager_endpoint.testAzure"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(azureResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(azureResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint_status", "Enabled"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMTrafficManagerEndpoint_requiresImport),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
					testCheckAzureRMTrafficManagerEndpointDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_basicDisableExternal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_basicDisableExternal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Disabled"),
				),
			},
		},
	})
}

// Altering weight might be used to ramp up migration traffic
func TestAccAzureRMTrafficManagerEndpoint_updateWeight(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_weight(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "weight", "50"),
					resource.TestCheckResourceAttr(secondResourceName, "weight", "50"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_updateWeight(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "weight", "25"),
					resource.TestCheckResourceAttr(secondResourceName, "weight", "75"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_updateSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_subnets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.0.first", "1.2.3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.0.scope", "24"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.1.first", "11.12.13.14"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.1.last", "11.12.13.14"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.0.first", "21.22.23.24"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.0.scope", "32"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_updateSubnets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet.#", "0"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.0.first", "12.34.56.78"),
					resource.TestCheckResourceAttr(secondResourceName, "subnet.0.last", "12.34.56.78"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_updateCustomeHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_headers(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_header.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_header.0.name", "header"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_header.0.value", "www.bing.com"),
					resource.TestCheckResourceAttr(secondResourceName, "custom_header.#", "0"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_updateHeaders(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_header.#", "0"),
					resource.TestCheckResourceAttr(secondResourceName, "custom_header.#", "1"),
					resource.TestCheckResourceAttr(secondResourceName, "custom_header.0.name", "header"),
					resource.TestCheckResourceAttr(secondResourceName, "custom_header.0.value", "www.bing.com"),
				),
			},
		},
	})
}

// Altering priority might be used to switch failover/active roles
func TestAccAzureRMTrafficManagerEndpoint_updatePriority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_priority(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "1"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "2"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_updatePriority(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "3"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_nestedEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "nested")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_nestedEndpoints(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists("azurerm_traffic_manager_endpoint.nested"),
					testCheckAzureRMTrafficManagerEndpointExists("azurerm_traffic_manager_endpoint.externalChild"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_location(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_location(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_locationUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_withGeoMappings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerEndpoint_geoMappings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.0", "GB"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.1", "FR"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerEndpoint_geoMappingsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.0", "FR"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_mappings.1", "DE"),
				),
			},
		},
	})
}

func testCheckAzureRMTrafficManagerEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		resp, err := conn.Get(ctx, resourceGroup, profileName, path.Base(endpointType), name)
		if err != nil {
			return fmt.Errorf("Bad: Get on trafficManagerEndpointsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Traffic Manager Endpoint %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTrafficManagerEndpointDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.EndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		if _, err := conn.Delete(ctx, resourceGroup, profileName, path.Base(endpointType), name); err != nil {
			return fmt.Errorf("Bad: Delete on trafficManagerEndpointsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMTrafficManagerEndpointDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.EndpointsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_traffic_manager_endpoint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, profileName, path.Base(endpointType), name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Traffic Manager Endpoint sitll exists:\n%#v", resp.EndpointProperties)
		}
	}

	return nil
}

func testAccAzureRMTrafficManagerEndpoint_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestpublicip-%d"
}

resource "azurerm_traffic_manager_endpoint" "testAzure" {
  name                = "acctestend-azure%d"
  type                = "azureEndpoints"
  target_resource_id  = azurerm_public_ip.test.id
  weight              = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  weight              = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerEndpoint_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_endpoint" "import" {
  name                = azurerm_traffic_manager_endpoint.testAzure.name
  type                = azurerm_traffic_manager_endpoint.testAzure.type
  target_resource_id  = azurerm_traffic_manager_endpoint.testAzure.target_resource_id
  weight              = azurerm_traffic_manager_endpoint.testAzure.weight
  profile_name        = azurerm_traffic_manager_endpoint.testAzure.profile_name
  resource_group_name = azurerm_traffic_manager_endpoint.testAzure.resource_group_name
}
`, template)
}

func testAccAzureRMTrafficManagerEndpoint_basicDisableExternal(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestpublicip-%d"
}

resource "azurerm_traffic_manager_endpoint" "testAzure" {
  name                = "acctestend-azure%d"
  type                = "azureEndpoints"
  target_resource_id  = azurerm_public_ip.test.id
  weight              = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  endpoint_status     = "Disabled"
  type                = "externalEndpoints"
  target              = "terraform.io"
  weight              = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_weight(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  weight              = 50
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  weight              = 50
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_updateWeight(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  weight              = 25
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  weight              = 75
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_priority(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  priority            = 1
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  priority            = 2
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_updatePriority(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  priority            = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  priority            = 2
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_subnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Subnet"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  subnet {
    first = "1.2.3.0"
    scope = "24"
  }
  subnet {
    first = "11.12.13.14"
    last  = "11.12.13.14"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  subnet {
    first = "21.22.23.24"
    scope = "32"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_updateSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Subnet"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  subnet {
    first = "12.34.56.78"
    last  = "12.34.56.78"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 1
  custom_header {
    name  = "header"
    value = "www.bing.com"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_updateHeaders(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 1
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
  name                = "acctestend-external%d-2"
  type                = "externalEndpoints"
  target              = "www.terraform.io"
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 2
  custom_header {
    name  = "header"
    value = "www.bing.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_nestedEndpoints(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "parent" {
  name                   = "acctesttmpparent%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctestparent%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_profile" "child" {
  name                   = "acctesttmpchild%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmpchild%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "nested" {
  name                = "acctestend-parent%d"
  type                = "nestedEndpoints"
  target_resource_id  = azurerm_traffic_manager_profile.child.id
  priority            = 1
  profile_name        = azurerm_traffic_manager_profile.parent.name
  resource_group_name = azurerm_resource_group.test.name
  min_child_endpoints = 1
}

resource "azurerm_traffic_manager_endpoint" "externalChild" {
  name                = "acctestend-child%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  priority            = 1
  profile_name        = azurerm_traffic_manager_profile.child.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_location(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmpparent%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctestparent%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  endpoint_location   = azurerm_resource_group.test.location
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_locationUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmpparent%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctestparent%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "acctestend-external%d"
  type                = "externalEndpoints"
  target              = "terraform.io"
  endpoint_location   = azurerm_resource_group.test.location
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_geoMappings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "example.com"
  resource_group_name = azurerm_resource_group.test.name
  profile_name        = azurerm_traffic_manager_profile.test.name
  target              = "example.com"
  type                = "externalEndpoints"
  geo_mappings        = ["GB", "FR"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerEndpoint_geoMappingsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "example.com"
  resource_group_name = azurerm_resource_group.test.name
  profile_name        = azurerm_traffic_manager_profile.test.name
  target              = "example.com"
  type                = "externalEndpoints"
  geo_mappings        = ["FR", "DE"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
