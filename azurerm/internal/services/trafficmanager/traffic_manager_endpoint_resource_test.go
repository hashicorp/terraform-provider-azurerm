package trafficmanager_test

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type TrafficManagerEndpointResource struct{}

func TestAccAzureRMTrafficManagerEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	azureResourceName := "azurerm_traffic_manager_endpoint.testAzure"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(azureResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(azureResourceName).Key("endpoint_status").HasValue("Enabled"),
				check.That(data.ResourceName).Key("endpoint_status").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMTrafficManagerEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(externalResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint_status").HasValue("Enabled"),
				check.That(externalResourceName).Key("endpoint_status").HasValue("Enabled"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMTrafficManagerEndpoint_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccAzureRMTrafficManagerEndpoint_basicDisableExternal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testAzure")
	externalResourceName := "azurerm_traffic_manager_endpoint.testExternal"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(externalResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint_status").HasValue("Enabled"),
				check.That(externalResourceName).Key("endpoint_status").HasValue("Enabled"),
			),
		},
		{
			Config: r.basicDisableExternal(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(externalResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint_status").HasValue("Enabled"),
				check.That(externalResourceName).Key("endpoint_status").HasValue("Disabled"),
			),
		},
	})
}

// Altering weight might be used to ramp up migration traffic
func TestAccAzureRMTrafficManagerEndpoint_updateWeight(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.weight(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("weight").HasValue("50"),
				check.That(secondResourceName).Key("weight").HasValue("50"),
			),
		},
		{
			Config: r.updateWeight(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("weight").HasValue("25"),
				check.That(secondResourceName).Key("weight").HasValue("75"),
			),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_updateSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet.#").HasValue("2"),
				check.That(data.ResourceName).Key("subnet.0.first").HasValue("1.2.3.0"),
				check.That(data.ResourceName).Key("subnet.0.scope").HasValue("24"),
				check.That(data.ResourceName).Key("subnet.1.first").HasValue("11.12.13.14"),
				check.That(data.ResourceName).Key("subnet.1.last").HasValue("11.12.13.14"),
				check.That(secondResourceName).Key("subnet.#").HasValue("1"),
				check.That(secondResourceName).Key("subnet.0.first").HasValue("21.22.23.24"),
				check.That(secondResourceName).Key("subnet.0.scope").HasValue("32"),
			),
		},
		{
			Config: r.updateSubnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet.#").HasValue("0"),
				check.That(secondResourceName).Key("subnet.#").HasValue("1"),
				check.That(secondResourceName).Key("subnet.0.first").HasValue("12.34.56.78"),
				check.That(secondResourceName).Key("subnet.0.last").HasValue("12.34.56.78"),
			),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_updateCustomeHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_header.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_header.0.name").HasValue("header"),
				check.That(data.ResourceName).Key("custom_header.0.value").HasValue("www.bing.com"),
				check.That(secondResourceName).Key("custom_header.#").HasValue("0"),
			),
		},
		{
			Config: r.updateHeaders(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_header.#").HasValue("0"),
				check.That(secondResourceName).Key("custom_header.#").HasValue("1"),
				check.That(secondResourceName).Key("custom_header.0.name").HasValue("header"),
				check.That(secondResourceName).Key("custom_header.0.value").HasValue("www.bing.com"),
			),
		},
	})
}

// Altering priority might be used to switch failover/active roles
func TestAccAzureRMTrafficManagerEndpoint_updatePriority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "testExternal")
	secondResourceName := "azurerm_traffic_manager_endpoint.testExternalNew"
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.priority(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("priority").HasValue("1"),
				check.That(secondResourceName).Key("priority").HasValue("2"),
			),
		},
		{
			Config: r.updatePriority(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("priority").HasValue("3"),
				check.That(secondResourceName).Key("priority").HasValue("2"),
			),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_nestedEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "nested")
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nestedEndpoints(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_traffic_manager_endpoint.nested").ExistsInAzure(r),
				check.That("azurerm_traffic_manager_endpoint.externalChild").ExistsInAzure(r),
			),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_location(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "test")
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.location(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.locationUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccAzureRMTrafficManagerEndpoint_withGeoMappings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_endpoint", "test")
	r := TrafficManagerEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.geoMappings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_mappings.#").HasValue("2"),
				check.That(data.ResourceName).Key("geo_mappings.0").HasValue("GB"),
				check.That(data.ResourceName).Key("geo_mappings.1").HasValue("FR"),
			),
		},
		{
			Config: r.geoMappingsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_mappings.#").HasValue("2"),
				check.That(data.ResourceName).Key("geo_mappings.0").HasValue("FR"),
				check.That(data.ResourceName).Key("geo_mappings.1").HasValue("DE"),
			),
		},
	})
}

func (r TrafficManagerEndpointResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	endpointType := state.Attributes["type"]
	profileName := state.Attributes["profile_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.TrafficManager.EndpointsClient.Get(ctx, resourceGroup, profileName, path.Base(endpointType), name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Traffic Manager Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r TrafficManagerEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	endpointType := state.Attributes["type"]
	profileName := state.Attributes["profile_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	if _, err := client.TrafficManager.EndpointsClient.Delete(ctx, resourceGroup, profileName, path.Base(endpointType), name); err != nil {
		return nil, fmt.Errorf("deleting Traffic Manager Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r TrafficManagerEndpointResource) basic(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
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

func (r TrafficManagerEndpointResource) basicDisableExternal(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) weight(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) updateWeight(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) priority(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) updatePriority(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) subnets(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) updateSubnets(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) headers(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) updateHeaders(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) nestedEndpoints(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) location(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) locationUpdated(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) geoMappings(data acceptance.TestData) string {
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

func (r TrafficManagerEndpointResource) geoMappingsUpdated(data acceptance.TestData) string {
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
