package azurerm

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMTrafficManagerEndpoint_basic(t *testing.T) {
	azureResourceName := "azurerm_traffic_manager_endpoint.testAzure"
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerEndpoint_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(azureResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(azureResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_disappears(t *testing.T) {
	azureResourceName := "azurerm_traffic_manager_endpoint.testAzure"
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerEndpoint_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(azureResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(azureResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
					testCheckAzureRMTrafficManagerEndpointDisappears(azureResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_basicDisableExternal(t *testing.T) {
	azureResourceName := "azurerm_traffic_manager_endpoint.testAzure"
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMTrafficManagerEndpoint_basic(ri, testLocation())
	postConfig := testAccAzureRMTrafficManagerEndpoint_basicDisableExternal(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(azureResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(azureResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Enabled"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(azureResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(externalResourceName),
					resource.TestCheckResourceAttr(azureResourceName, "endpoint_status", "Enabled"),
					resource.TestCheckResourceAttr(externalResourceName, "endpoint_status", "Disabled"),
				),
			},
		},
	})
}

// Altering weight might be used to ramp up migration traffic
func TestAccAzureRMTrafficManagerEndpoint_updateWeight(t *testing.T) {
	firstResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMTrafficManagerEndpoint_weight(ri, location)
	postConfig := testAccAzureRMTrafficManagerEndpoint_updateWeight(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(firstResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "weight", "50"),
					resource.TestCheckResourceAttr(secondResourceName, "weight", "50"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(firstResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "weight", "25"),
					resource.TestCheckResourceAttr(secondResourceName, "weight", "75"),
				),
			},
		},
	})
}

// Altering priority might be used to switch failover/active roles
func TestAccAzureRMTrafficManagerEndpoint_updatePriority(t *testing.T) {
	firstResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"

	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMTrafficManagerEndpoint_priority(ri, location)
	postConfig := testAccAzureRMTrafficManagerEndpoint_updatePriority(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(firstResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "1"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(firstResourceName),
					testCheckAzureRMTrafficManagerEndpointExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "3"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_nestedEndpoints(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerEndpoint_nestedEndpoints(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists("azurerm_traffic_manager_endpoint.nested"),
					testCheckAzureRMTrafficManagerEndpointExists("azurerm_traffic_manager_endpoint.externalChild"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_location(t *testing.T) {
	resourceName := "azurerm_traffic_manager_endpoint.test"
	ri := acctest.RandInt()
	location := testLocation()
	first := testAccAzureRMTrafficManagerEndpoint_location(ri, location)
	second := testAccAzureRMTrafficManagerEndpoint_locationUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: first,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(resourceName),
				),
			},
			{
				Config: second,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_withGeoMappings(t *testing.T) {
	resourceName := "azurerm_traffic_manager_endpoint.test"
	ri := acctest.RandInt()
	location := testLocation()
	first := testAccAzureRMTrafficManagerEndpoint_geoMappings(ri, location)
	second := testAccAzureRMTrafficManagerEndpoint_geoMappingsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: first,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.0", "GB"),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.1", "FR"),
				),
			},
			{
				Config: second,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.0", "FR"),
					resource.TestCheckResourceAttr(resourceName, "geo_mappings.1", "DE"),
				),
			},
		},
	})
}

func testCheckAzureRMTrafficManagerEndpointExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		conn := testAccProvider.Meta().(*ArmClient).trafficManagerEndpointsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMTrafficManagerEndpointDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		conn := testAccProvider.Meta().(*ArmClient).trafficManagerEndpointsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		_, err := conn.Delete(ctx, resourceGroup, profileName, path.Base(endpointType), name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on trafficManagerEndpointsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMTrafficManagerEndpointDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).trafficManagerEndpointsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_traffic_manager_endpoint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		endpointType := rs.Primary.Attributes["type"]
		profileName := rs.Primary.Attributes["profile_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMTrafficManagerEndpoint_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Weighted"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
    domain_name_label = "acctestpublicip-%d"
}

resource "azurerm_traffic_manager_endpoint" "testAzure" {
    name = "acctestend-azure%d"
    type = "azureEndpoints"
    target_resource_id = "${azurerm_public_ip.test.id}"
    weight = 3
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    type = "externalEndpoints"
    target = "terraform.io"
    weight = 3
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_basicDisableExternal(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Weighted"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
    domain_name_label = "acctestpublicip-%d"
}

resource "azurerm_traffic_manager_endpoint" "testAzure" {
    name = "acctestend-azure%d"
    type = "azureEndpoints"
    target_resource_id = "${azurerm_public_ip.test.id}"
    weight = 3
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    endpoint_status = "Disabled"
    type = "externalEndpoints"
    target = "terraform.io"
    weight = 3
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_weight(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Weighted"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    type = "externalEndpoints"
    target = "terraform.io"
    weight = 50
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
    name = "acctestend-external%d-2"
    type = "externalEndpoints"
    target = "www.terraform.io"
    weight = 50
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_updateWeight(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Weighted"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    type = "externalEndpoints"
    target = "terraform.io"
    weight = 25
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
    name = "acctestend-external%d-2"
    type = "externalEndpoints"
    target = "www.terraform.io"
    weight = 75
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
func testAccAzureRMTrafficManagerEndpoint_priority(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Priority"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    type = "externalEndpoints"
    target = "terraform.io"
    priority = 1
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
    name = "acctestend-external%d-2"
    type = "externalEndpoints"
    target = "www.terraform.io"
    priority = 2
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_updatePriority(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Priority"

    dns_config {
        relative_name = "acctesttmp%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "testExternal" {
    name = "acctestend-external%d"
    type = "externalEndpoints"
    target = "terraform.io"
    priority = 3
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_traffic_manager_endpoint" "testExternalNew" {
    name = "acctestend-external%d-2"
    type = "externalEndpoints"
    target = "www.terraform.io"
    priority = 2
    profile_name = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_nestedEndpoints(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "parent" {
    name = "acctesttmpparent%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Priority"

    dns_config {
        relative_name = "acctestparent%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_profile" "child" {
    name = "acctesttmpchild%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Priority"

    dns_config {
        relative_name = "acctesttmpchild%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "nested" {
    name = "acctestend-parent%d"
    type = "nestedEndpoints"
    target_resource_id = "${azurerm_traffic_manager_profile.child.id}"
    priority = 1
    profile_name = "${azurerm_traffic_manager_profile.parent.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    min_child_endpoints = 1
}

resource "azurerm_traffic_manager_endpoint" "externalChild" {
    name = "acctestend-child%d"
    type = "externalEndpoints"
    target = "terraform.io"
    priority = 1
    profile_name = "${azurerm_traffic_manager_profile.child.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_location(rInt int, location string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name                   = "acctesttmpparent%d"
    resource_group_name    = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Performance"

    dns_config {
        relative_name = "acctestparent%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "test" {
    name                = "acctestend-external%d"
    type                = "externalEndpoints"
    target              = "terraform.io"
    endpoint_location   = "${azurerm_resource_group.test.location}"
    profile_name        = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_locationUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name                   = "acctesttmpparent%d"
    resource_group_name    = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Performance"

    dns_config {
        relative_name = "acctestparent%d"
        ttl = 30
    }

    monitor_config {
        protocol = "https"
        port = 443
        path = "/"
    }
}

resource "azurerm_traffic_manager_endpoint" "test" {
    name                = "acctestend-external%d"
    type                = "externalEndpoints"
    target              = "terraform.io"
    endpoint_location   = "${azurerm_resource_group.test.location}"
    profile_name        = "${azurerm_traffic_manager_profile.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_geoMappings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "example.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
  profile_name        = "${azurerm_traffic_manager_profile.test.name}"
  target              = "example.com"
  type                = "externalEndpoints"
  geo_mappings        = ["GB", "FR"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerEndpoint_geoMappingsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "example.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
  profile_name        = "${azurerm_traffic_manager_profile.test.name}"
  target              = "example.com"
  type                = "externalEndpoints"
  geo_mappings        = ["FR", "DE"]
}
`, rInt, location, rInt, rInt)
}
