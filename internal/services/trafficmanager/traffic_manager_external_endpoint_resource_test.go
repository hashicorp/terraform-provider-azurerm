// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ExternalEndpointResource struct{}

func TestAccExternalEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

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

func TestAccExternalEndpoint_priority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.priority(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExternalEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccExternalEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

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

func TestAccExternalEndpoint_subnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

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

func TestAccExternalEndpoint_multipleEndpointsWithDynamicPriority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExternalEndpoint_performancePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_external_endpoint", "test")
	r := ExternalEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.performancePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ExternalEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
	return utils.Bool(resp.Model != nil), nil
}

func (r ExternalEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name       = "acctestend-azure%d"
  target     = "www.example.com"
  profile_id = azurerm_traffic_manager_profile.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ExternalEndpointResource) priority(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name       = "acctestend-azure%[1]d"
  target     = "www.example.com"
  profile_id = azurerm_traffic_manager_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ExternalEndpointResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_external_endpoint" "import" {
  name       = azurerm_traffic_manager_external_endpoint.test.name
  target     = azurerm_traffic_manager_external_endpoint.test.target
  weight     = azurerm_traffic_manager_external_endpoint.test.weight
  profile_id = azurerm_traffic_manager_external_endpoint.test.profile_id
}
`, r.basic(data))
}

func (r ExternalEndpointResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name                 = "acctestend-azure%d"
  target               = "www.example.com"
  weight               = 5
  profile_id           = azurerm_traffic_manager_profile.test.id
  enabled              = false
  always_serve_enabled = true
  priority             = 4
  endpoint_location    = azurerm_resource_group.test.location

  geo_mappings = ["WORLD"]

  custom_header {
    name  = "header"
    value = "www.bing.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ExternalEndpointResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name              = "acctestend-azure%[2]d"
  target            = "www.example.com"
  weight            = 5
  profile_id        = azurerm_traffic_manager_profile.test.id
  endpoint_location = azurerm_resource_group.test.location
}

resource "azurerm_traffic_manager_external_endpoint" "test2" {
  name              = "acctestend-azure2%[2]d"
  target            = "www.pandas.com"
  weight            = 5
  profile_id        = azurerm_traffic_manager_profile.test.id
  endpoint_location = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ExternalEndpointResource) subnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Subnet"

  dns_config {
    relative_name = "acctest-tmp-%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name       = "acctestend-azure%[1]d"
  target     = "www.example.com"
  weight     = 5
  profile_id = azurerm_traffic_manager_profile.test.id

  subnet {
    first = "1.2.3.0"
    scope = "24"
  }
  subnet {
    first = "11.12.13.14"
    last  = "11.12.13.14"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ExternalEndpointResource) performancePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctest-tmp-%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_external_endpoint" "test" {
  name              = "acctestend-azure%[1]d"
  target            = "www.example.com"
  weight            = 3
  profile_id        = azurerm_traffic_manager_profile.test.id
  endpoint_location = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ExternalEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
