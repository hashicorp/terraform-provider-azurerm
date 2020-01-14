package tests

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_geographic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMTrafficManagerProfile_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_geographic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMTrafficManagerProfile_requiresImport),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weighted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_weighted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_weightedTCP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_weightedTCP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performance(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_performance(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Performance"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_priority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_priority(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Priority"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerProfile_withTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_performanceToGeographic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_performance(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Performance"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerProfile_geographic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_priorityToWeighted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	fqdn, err := getTrafficManagerFQDN(fmt.Sprintf("acctesttmp%d", data.RandomInteger))
	if err != nil {
		t.Fatalf("Error obtaining Azure Region: %+v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_priority(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Priority"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
			{
				Config: testAccAzureRMTrafficManagerProfile_weighted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fqdn),
				),
			},
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_fastEndpointFailoverSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_failover(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_failoverUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_fastEndpointFailoverSettingsError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMTrafficManagerProfile_failoverError(data),
				ExpectError: regexp.MustCompile("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`"),
			},
		},
	})
}

func testCheckAzureRMTrafficManagerProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).TrafficManager.ProfilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_traffic_manager_profile" {
			continue
		}

		log.Printf("[TRACE] test_profile %#v", rs)

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
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

func testAccAzureRMTrafficManagerProfile_geographic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
func testAccAzureRMTrafficManagerProfile_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_geographic(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_weighted(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_weightedTCP(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_performance(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_priority(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_withTags(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_withTagsUpdated(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_failover(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_failoverUpdate(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_failoverError(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
