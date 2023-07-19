// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualNetworkDnsServersResource struct{}

func TestAccVirtualNetworkDnsServers_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_dns_servers", "test")
	r := VirtualNetworkDnsServersResource{}

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

func (t VirtualNetworkDnsServersResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkDnsServersID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VnetClient.Get(ctx, id.ResourceGroup, id.VirtualNetworkName, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	exists := resp.ID != nil && resp.VirtualNetworkPropertiesFormat != nil && resp.VirtualNetworkPropertiesFormat.DhcpOptions != nil &&
		resp.VirtualNetworkPropertiesFormat.DhcpOptions.DNSServers != nil && len(*resp.VirtualNetworkPropertiesFormat.DhcpOptions.DNSServers) > 0

	return utils.Bool(exists), nil
}

func (VirtualNetworkDnsServersResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_virtual_network_dns_servers" "test" {
  virtual_network_id = azurerm_virtual_network.test.id
  dns_servers        = ["10.7.7.2", "10.7.7.7", "10.7.7.1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
