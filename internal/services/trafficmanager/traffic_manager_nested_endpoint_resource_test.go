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

type NestedEndpointResource struct{}

func TestAccNestedEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func TestAccNestedEndpoint_priority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func TestAccNestedEndpoint_multipleEndpointsWithDynamicPriority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func TestAccNestedEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func TestAccNestedEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func TestAccNestedEndpoint_subnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_nested_endpoint", "test")
	r := NestedEndpointResource{}

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

func (r NestedEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r NestedEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                    = "acctestend-parent%d"
  target_resource_id      = azurerm_traffic_manager_profile.child.id
  profile_id              = azurerm_traffic_manager_profile.parent.id
  minimum_child_endpoints = 5
}
`, r.template(data), data.RandomInteger)
}

func (r NestedEndpointResource) priority(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "parent" {
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

resource "azurerm_traffic_manager_profile" "child" {
  name                   = "acctesttmpchild%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmpchild%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                    = "acctestend-parent%[1]d"
  target_resource_id      = azurerm_traffic_manager_profile.child.id
  profile_id              = azurerm_traffic_manager_profile.parent.id
  minimum_child_endpoints = 5
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NestedEndpointResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_nested_endpoint" "import" {
  name                    = azurerm_traffic_manager_nested_endpoint.test.name
  target_resource_id      = azurerm_traffic_manager_nested_endpoint.test.target_resource_id
  profile_id              = azurerm_traffic_manager_nested_endpoint.test.profile_id
  minimum_child_endpoints = azurerm_traffic_manager_nested_endpoint.test.minimum_child_endpoints
  weight                  = azurerm_traffic_manager_nested_endpoint.test.weight
}
`, r.basic(data))
}

func (r NestedEndpointResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                                  = "acctestend-parent%d"
  target_resource_id                    = azurerm_traffic_manager_profile.child.id
  priority                              = 3
  profile_id                            = azurerm_traffic_manager_profile.parent.id
  weight                                = 5
  minimum_child_endpoints               = 9
  minimum_required_child_endpoints_ipv4 = 2
  minimum_required_child_endpoints_ipv6 = 2
  endpoint_location                     = azurerm_resource_group.test.location

  geo_mappings = ["WORLD"]
}
`, r.template(data), data.RandomInteger)
}

func (r NestedEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_traffic_manager_profile" "parent" {
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

resource "azurerm_traffic_manager_profile" "child" {
  name                   = "acctesttmpchild%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmpchild%[1]d"
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

func (r NestedEndpointResource) subnets(data acceptance.TestData) string {
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

resource "azurerm_traffic_manager_profile" "child" {
  name                   = "acctesttmpchild%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmpchild%[1]d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                    = "acctestend-parent%[1]d"
  target_resource_id      = azurerm_traffic_manager_profile.child.id
  profile_id              = azurerm_traffic_manager_profile.test.id
  minimum_child_endpoints = 5
  weight                  = 3

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

func (r NestedEndpointResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_traffic_manager_profile" "child2" {
  name                   = "acctesttmpchild%[2]d-2"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctesttmpchild%[2]d-2"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}

resource "azurerm_traffic_manager_nested_endpoint" "test" {
  name                    = "acctestend-parent%[2]d"
  target_resource_id      = azurerm_traffic_manager_profile.child.id
  profile_id              = azurerm_traffic_manager_profile.parent.id
  minimum_child_endpoints = 5
}

resource "azurerm_traffic_manager_nested_endpoint" "test2" {
  name                    = "acctestend-parent%[2]d-2"
  target_resource_id      = azurerm_traffic_manager_profile.child2.id
  profile_id              = azurerm_traffic_manager_profile.parent.id
  minimum_child_endpoints = 5
}
`, r.template(data), data.RandomInteger)
}
