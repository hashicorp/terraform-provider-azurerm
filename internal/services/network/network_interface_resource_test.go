// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkInterfaceResource struct{}

func TestAccNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
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

func TestAccNetworkInterface_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccNetworkInterface_dnsServers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsServers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dnsServersUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_enableAcceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enabled
			Config: r.enableAcceleratedNetworking(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.enableAcceleratedNetworking(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.enableAcceleratedNetworking(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_enableIPForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enabled
			Config: r.enableIPForwarding(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.enableIPForwarding(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.enableIPForwarding(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_internalDomainNameLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internalDomainNameLabel(data, "1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.internalDomainNameLabel(data, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("ip_configuration.1.private_ip_address_version").HasValue("IPv6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_multipleIPConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleIPConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_multipleIPConfigurationsSecondaryAsPrimary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleIPConfigurationsSecondaryAsPrimary(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_publicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicIP(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicIPRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicIP(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface"),
		},
	})
}

func TestAccNetworkInterface_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
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

func TestAccNetworkInterface_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleIPConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_updateMultipleParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultipleParameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateMultipleParameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_pointToGatewayLB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pointToGatewayLB(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkInterfaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseNetworkInterfaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NetworkInterfaces.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkInterfaceResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseNetworkInterfaceID(state.ID)
	if err != nil {
		return nil, err
	}

	err = client.Network.NetworkInterfaces.DeleteThenPoll(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r NetworkInterfaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) withMultipleParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r NetworkInterfaceResource) updateMultipleParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.7"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r NetworkInterfaceResource) dnsServers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) dnsServersUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.6",
    "10.0.0.5"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) edgeZone(data acceptance.TestData) string {
	// @tombuildsstuff: WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkInterfaceResource) enableAcceleratedNetworking(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                          = "acctestni-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  enable_accelerated_networking = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r NetworkInterfaceResource) enableIPForwarding(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                 = "acctestni-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  enable_ip_forwarding = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r NetworkInterfaceResource) internalDomainNameLabel(data acceptance.TestData, suffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  internal_dns_name_label = "acctestni-%s-%s"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, suffix, data.RandomString)
}

func (r NetworkInterfaceResource) ipv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    private_ip_address_allocation = "Dynamic"
    private_ip_address_version    = "IPv6"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) multipleIPConfigurations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) multipleIPConfigurationsSecondaryAsPrimary(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
`, r.publicIPTemplate(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIPRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.publicIPTemplate(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIPTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "import" {
  name                = azurerm_network_interface.test.name
  location            = azurerm_network_interface.test.location
  resource_group_name = azurerm_network_interface.test.resource_group_name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.basic(data))
}

func (r NetworkInterfaceResource) static(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.15"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello     = "World"
    Elephants = "Five"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) pointToGatewayLB(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "gateway" {
  name                = "acctestvnet-gw-%[2]d"
  address_space       = ["11.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "gateway" {
  name                 = "acctestsubnet-gw-%[2]d"
  resource_group_name  = azurerm_virtual_network.gateway.resource_group_name
  virtual_network_name = azurerm_virtual_network.gateway.name
  address_prefixes     = ["11.0.2.0/24"]
}

resource "azurerm_lb" "gateway" {
  name                = "acctestlb-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Gateway"

  frontend_ip_configuration {
    name      = "feip"
    subnet_id = azurerm_subnet.gateway.id
  }
}

resource "azurerm_lb_backend_address_pool" "gateway" {
  name            = "acctestbap-%[2]d"
  loadbalancer_id = azurerm_lb.gateway.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                                               = "gateway"
    public_ip_address_id                               = azurerm_public_ip.test.id
    gateway_load_balancer_frontend_ip_configuration_id = azurerm_lb.gateway.frontend_ip_configuration.0.id
    private_ip_address_allocation                      = "Dynamic"
    subnet_id                                          = azurerm_subnet.test.id
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
