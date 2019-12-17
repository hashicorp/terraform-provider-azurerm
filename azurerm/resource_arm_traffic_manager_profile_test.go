package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func getTrafficManagerFQDN(hostname string) (string, error) {
	environment, err := acceptance.Environment()
	if err != nil {
		return "", err
	}
	dnsSuffix := environment.TrafficManagerDNSSuffix
	return fmt.Sprintf("%s.%s", hostname, dnsSuffix), nil
}

func TestAccAzureRMTrafficManagerProfile_geographic(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_geographic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Geographic"),
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
func TestAccAzureRMTrafficManagerProfile_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_geographic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_method", "Geographic"),
				),
			},
			{
				Config:      testAccAzureRMTrafficManagerProfile_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_traffic_manager_profile"),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weighted(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMTrafficManagerProfile_weighted(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weightedTCP(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMTrafficManagerProfile_weightedTCP(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performance(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMTrafficManagerProfile_performance(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_priority(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMTrafficManagerProfile_priority(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_withTags(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMTrafficManagerProfile_withTags(ri, acceptance.Location())
	postConfig := testAccAzureRMTrafficManagerProfile_withTagsUpdated(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performanceToGeographic(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMTrafficManagerProfile_performance(ri, acceptance.Location())
	postConfig := testAccAzureRMTrafficManagerProfile_geographic(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMTrafficManagerProfile_priority(ri, acceptance.Location())
	postConfig := testAccAzureRMTrafficManagerProfile_weighted(ri, acceptance.Location())

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", ri))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func TestAccAzureRMTrafficManagerProfile_fastEndpointFailoverSettings(t *testing.T) {
	resourceName := "azurerm_traffic_manager_profile.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_failover(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMTrafficManagerProfile_failoverUpdate(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(resourceName),
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

func TestAccAzureRMTrafficManagerProfile_fastEndpointFailoverSettingsError(t *testing.T) {
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMTrafficManagerProfile_failoverError(rInt, location),
				ExpectError: regexp.MustCompile("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`"),
			},
		},
	})
}

func testCheckAzureRMTrafficManagerProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Traffic Manager Profile: %s", name)
		}

		// Ensure resource group/virtual network combination exists in API
		conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.ProfilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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
	conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.ProfilesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_traffic_manager_profile" {
			continue
		}

		log.Printf("[TRACE] test_profile %#v", rs)

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}
`, rInt, location, rInt, rInt)
}
func testAccAzureRMTrafficManagerProfile_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "import" {
  name                   = "${azurerm_traffic_manager_profile.test.name}"
  resource_group_name    = "${azurerm_traffic_manager_profile.test.resource_group_name}"
  traffic_routing_method = "${azurerm_traffic_manager_profile.test.traffic_routing_method}"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}
`, testAccAzureRMTrafficManagerProfile_geographic(rInt, location), rInt)
}

func testAccAzureRMTrafficManagerProfile_weighted(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_weightedTCP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "tcp"
    port     = 443
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_performance(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_priority(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_withTagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "https"
    port     = 443
    path     = "/"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_failover(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol                     = "https"
    port                         = 443
    path                         = "/"
    interval_in_seconds          = 30
    timeout_in_seconds           = 6
    tolerated_number_of_failures = 3
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_failoverUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol                     = "https"
    port                         = 443
    path                         = "/"
    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 6
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMTrafficManagerProfile_failoverError(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctesttmp%d"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctesttmp%d"
    ttl           = 30
  }

  monitor_config {
    protocol                     = "https"
    port                         = 443
    path                         = "/"
    interval_in_seconds          = 10
    timeout_in_seconds           = 10
    tolerated_number_of_failures = 3
  }
}
`, rInt, location, rInt, rInt)
}
