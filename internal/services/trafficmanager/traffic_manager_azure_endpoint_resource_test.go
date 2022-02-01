package trafficmanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/sdk/2018-08-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TrafficManagerAzureEndpointResource struct{}

func TestAccAzureRMTrafficManagerAzureEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_azure_endpoint", "test")
	r := TrafficManagerAzureEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMTrafficManagerAzureEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_azure_endpoint", "test")
	r := TrafficManagerAzureEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMTrafficManagerAzureEndpoint_subnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_azure_endpoint", "test")
	r := TrafficManagerAzureEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subnets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TrafficManagerAzureEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.TrafficManager.EndpointsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r TrafficManagerAzureEndpointResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.TrafficManager.EndpointsClient.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r TrafficManagerAzureEndpointResource) basic(data acceptance.TestData) string {
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

resource "azurerm_traffic_manager_azure_endpoint" "test" {
  name                = "acctestend-azure%d"
  target_resource_id  = azurerm_public_ip.test.id
  weight              = 3
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerAzureEndpointResource) complete(data acceptance.TestData) string {
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

resource "azurerm_traffic_manager_azure_endpoint" "test" {
  name                = "acctestend-azure%d"
  target_resource_id  = azurerm_public_ip.test.id
  weight              = 5
  profile_name        = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
  enabled             = false
  priority            = 4

  geo_mappings = ["WORLD"]

  custom_header {
    name  = "header"
    value = "www.bing.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerAzureEndpointResource) subnets(data acceptance.TestData) string {
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

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestpublicip-%d"
}

resource "azurerm_traffic_manager_azure_endpoint" "test" {
  name                = "acctestend-azure%d"
  target_resource_id  = azurerm_public_ip.test.id
  weight              = 5
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
