// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualNetworkPeeringResource struct{}

func TestAccVirtualNetworkPeering_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(secondResourceName).Key("allow_virtual_network_access").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkPeering_withTriggers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTriggers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(data.ResourceName).Key("triggers.remote_address_space").Exists(),
				check.That(data.ResourceName).Key("triggers.remote_address_space").HasValue("10.0.2.0/24,1001:1002::/64"),
				check.That(secondResourceName).Key("allow_virtual_network_access").HasValue("true"),
			),
		},
		// triggers is an arbitrary list(string) which
		// is not known at the backend API
		data.ImportStep("triggers"),
	})
}

func TestAccVirtualNetworkPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualNetworkPeering_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccVirtualNetworkPeering_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(secondResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(data.ResourceName).Key("allow_forwarded_traffic").HasValue("false"),
				check.That(secondResourceName).Key("allow_forwarded_traffic").HasValue("false"),
			),
		},

		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(secondResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(data.ResourceName).Key("allow_forwarded_traffic").HasValue("true"),
				check.That(secondResourceName).Key("allow_forwarded_traffic").HasValue("true"),
			),
		},
	})
}

func TestAccVirtualNetworkPeering_subnetPeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subnetPeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VirtualNetworkPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Network.VirtualNetworkPeerings.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r VirtualNetworkPeeringResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.VirtualNetworkPeerings.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting on virtual network peering: %+v", err)
	}

	return utils.Bool(true), nil
}

func (r VirtualNetworkPeeringResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_virtual_network_access = true
}
`, template, data.RandomInteger)
}

func (r VirtualNetworkPeeringResource) withTriggers(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_virtual_network_access = true
  triggers = {
    remote_address_space = join(",", azurerm_virtual_network.test2.address_space)
  }
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_virtual_network_access = true
  triggers = {
    remote_address_space = join(",", azurerm_virtual_network.test1.address_space)
  }
}
`, template, data.RandomInteger)
}

func (r VirtualNetworkPeeringResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_peering" "import" {
  name                         = azurerm_virtual_network_peering.test1.name
  resource_group_name          = azurerm_virtual_network_peering.test1.resource_group_name
  virtual_network_name         = azurerm_virtual_network_peering.test1.virtual_network_name
  remote_virtual_network_id    = azurerm_virtual_network_peering.test1.remote_virtual_network_id
  allow_virtual_network_access = azurerm_virtual_network_peering.test1.allow_virtual_network_access
}
`, r.basic(data))
}

func (r VirtualNetworkPeeringResource) basicUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}
`, template, data.RandomInteger)
}

func (VirtualNetworkPeeringResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = %[2]q
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24", "1001:1001::/64"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24", "1001:1002::/64"]
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualNetworkPeeringResource) subnetPeering(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_subnet" "test1" {
  name                 = "internal1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test1.name
  address_prefixes     = ["10.0.1.0/27", "1001:1001::/64"]
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.0.2.0/27", "1001:1002::/64"]
}

resource "azurerm_virtual_network_peering" "test1" {
  name                                   = "acctestpeer-1-%[2]d"
  resource_group_name                    = azurerm_resource_group.test.name
  virtual_network_name                   = azurerm_virtual_network.test1.name
  remote_virtual_network_id              = azurerm_virtual_network.test2.id
  allow_forwarded_traffic                = true
  allow_virtual_network_access           = true
  peer_complete_virtual_networks_enabled = false
  only_ipv6_peering_enabled              = true
  local_subnet_names                     = [azurerm_subnet.test1.name]
  remote_subnet_names                    = [azurerm_subnet.test2.name]
}
`, r.template(data), data.RandomInteger)
}
