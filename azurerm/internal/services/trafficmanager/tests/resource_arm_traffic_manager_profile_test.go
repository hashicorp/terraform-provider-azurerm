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
)

func TestAccAzureRMTrafficManagerProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Geographic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_updateEnsureDoNotEraseEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_completeWithEndpoint(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_completeUpdatedWithEndpoint(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Geographic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMTrafficManagerProfile_requiresImport),
		},
	})
}

func TestAccAzureRMTrafficManagerProfile_cycleMethod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTrafficManagerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Geographic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Weighted"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Weighted"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Subnet"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Subnet"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Priority"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Priority"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTrafficManagerProfile_basic(data, "Performance"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTrafficManagerProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Performance"),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
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

func testAccAzureRMTrafficManagerProfile_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMTrafficManagerProfile_basic(data acceptance.TestData, method string) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "%s"

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
`, template, data.RandomInteger, method, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_basic(data, "Geographic")
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "import" {
  name                   = azurerm_traffic_manager_profile.test.name
  resource_group_name    = azurerm_traffic_manager_profile.test.resource_group_name
  traffic_routing_method = azurerm_traffic_manager_profile.test.traffic_routing_method

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
`, template, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_complete(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "100-101",
      "301-303",
    ]

    custom_header {
      name  = "foo"
      value = "bar"
    }

    protocol = "tcp"
    port     = 777

    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 6
  }

  tags = {
    Environment = "Production"
    cost_center = "acctest"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_completeUpdated(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "302-304",
    ]

    custom_header {
      name  = "foo2"
      value = "bar2"
    }

    protocol = "https"
    port     = 442
    path     = "/"

    interval_in_seconds          = 30
    timeout_in_seconds           = 6
    tolerated_number_of_failures = 3
  }

  tags = {
    Environment = "staging"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_endpointResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_traffic_manager_endpoint" "test" {
  name                = "acctestend-external%d"
  resource_group_name = azurerm_resource_group.test.name
  profile_name        = azurerm_traffic_manager_profile.test.name
  target              = "terraform.io"
  type                = "externalEndpoints"
  weight              = 100
}
`, data.RandomInteger)
}

func testAccAzureRMTrafficManagerProfile_completeWithEndpoint(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	endpoint := testAccAzureRMTrafficManagerProfile_endpointResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "100-101",
      "301-303",
    ]

    custom_header {
      name  = "foo"
      value = "bar"
    }

    protocol = "tcp"
    port     = 777

    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 6
  }

  tags = {
    Environment = "Production"
    cost_center = "acctest"
  }
}

%s
`, template, data.RandomInteger, data.RandomInteger, endpoint)
}

func testAccAzureRMTrafficManagerProfile_completeUpdatedWithEndpoint(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	endpoint := testAccAzureRMTrafficManagerProfile_endpointResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "302-304",
    ]

    custom_header {
      name  = "foo2"
      value = "bar2"
    }

    protocol = "https"
    port     = 442
    path     = "/"

    interval_in_seconds          = 30
    timeout_in_seconds           = 6
    tolerated_number_of_failures = 3
  }

  tags = {
    Environment = "staging"
  }
}

%s
`, template, data.RandomInteger, data.RandomInteger, endpoint)
}

func testAccAzureRMTrafficManagerProfile_failoverError(data acceptance.TestData) string {
	template := testAccAzureRMTrafficManagerProfile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctest-tmp-%d"
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
`, template, data.RandomInteger, data.RandomInteger)
}
