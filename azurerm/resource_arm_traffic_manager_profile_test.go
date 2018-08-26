package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func getTrafficManagerFQDN(hostname string) (string, error) {
	environment, err := testArmEnvironment()
	if err != nil {
		return "", err
	}
	dnsSuffix := environment.TrafficManagerDNSSuffix
	return fmt.Sprintf("%s.%s", hostname, dnsSuffix), nil
}

func TestAccAzureRMTrafficManagerProfile_geographic(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerProfile_geographic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Geographic"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weighted(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerProfile_weighted(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weightedTCP(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerProfile_weightedTCP(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performance(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerProfile_performance(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Performance"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_priority(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMTrafficManagerProfile_priority(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Priority"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_withTags(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMTrafficManagerProfile_withTags(ri, testLocation())
	postConfig := testAccAzureRMTrafficManagerProfile_withTagsUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performanceToGeographic(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMTrafficManagerProfile_performance(ri, testLocation())
	postConfig := testAccAzureRMTrafficManagerProfile_geographic(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Performance"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Geographic"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_priorityToWeighted(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMTrafficManagerProfile_priority(ri, testLocation())
	postConfig := testAccAzureRMTrafficManagerProfile_weighted(ri, testLocation())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Priority"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func testCheckAzureRMTrafficManagerProfileExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		conn := testAccProvider.Meta().(*ArmClient).trafficManagerProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on trafficManagerProfilesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Traffic Manager %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTrafficManagerProfileDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).trafficManagerProfilesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_traffic_manager_profile" {
			continue
		}

		log.Printf("[TRACE] test_profile %#v", rs)

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Traffic Manager profile sitll exists:\n%#v", resp.ProfileProperties)
		}
	}

	return nil
}

func testAccAzureRMTrafficManagerProfile_geographic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Geographic"

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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_weighted(rInt int, location string) string {
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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_weightedTCP(rInt int, location string) string {
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
        protocol = "tcp"
        port = 443
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_performance(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
    name = "acctesttmp%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    traffic_routing_method = "Performance"

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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_priority(rInt int, location string) string {
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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_withTags(rInt int, location string) string {
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

    tags {
        environment = "Production"
        cost_center = "MSFT"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_withTagsUpdated(rInt int, location string) string {
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

    tags {
        environment = "staging"
    }
}
`, rInt, location, rInt, rInt)
}
