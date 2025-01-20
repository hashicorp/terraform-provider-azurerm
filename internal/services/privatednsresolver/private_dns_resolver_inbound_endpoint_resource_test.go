// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DNSResolverInboundEndpointResource struct{}

func TestAccDNSResolverInboundEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_inbound_endpoint", "test")
	r := DNSResolverInboundEndpointResource{}
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

func TestAccDNSResolverInboundEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_inbound_endpoint", "test")
	r := DNSResolverInboundEndpointResource{}
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

func TestAccDNSResolverInboundEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_inbound_endpoint", "test")
	r := DNSResolverInboundEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDNSResolverInboundEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_inbound_endpoint", "test")
	r := DNSResolverInboundEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDNSResolverInboundEndpoint_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_inbound_endpoint", "test")
	r := DNSResolverInboundEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DNSResolverInboundEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := inboundendpoints.ParseInboundEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.InboundEndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r DNSResolverInboundEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}


resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "inbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r DNSResolverInboundEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  ip_configurations {
    subnet_id = azurerm_subnet.test.id
  }
}
`, template, data.RandomInteger)
}

func (r DNSResolverInboundEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "import" {
  name                    = azurerm_private_dns_resolver_inbound_endpoint.test.name
  private_dns_resolver_id = azurerm_private_dns_resolver_inbound_endpoint.test.private_dns_resolver_id
  location                = azurerm_private_dns_resolver_inbound_endpoint.test.location
  ip_configurations {
    subnet_id = azurerm_subnet.test.id
  }
}
`, config)
}

func (r DNSResolverInboundEndpointResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  ip_configurations {
    private_ip_allocation_method = "Dynamic"
    subnet_id                    = azurerm_subnet.test.id
  }
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r DNSResolverInboundEndpointResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  ip_configurations {
    subnet_id = azurerm_subnet.test.id
  }
  tags = {
    key = "updated value"
  }
}
`, template, data.RandomInteger)
}

func (r DNSResolverInboundEndpointResource) static(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location

  ip_configurations {
    subnet_id                    = azurerm_subnet.test.id
    private_ip_allocation_method = "Static"
    private_ip_address           = "10.0.0.4"
  }
}
`, template, data.RandomInteger)
}
