// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkInterfaceBackendAddressPoolResource struct{}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface_backend_address_pool_association"),
		},
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentionally not using a DisppearsStep as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccNetworkInterfaceBackendAddressPoolAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	r := NetworkInterfaceBackendAddressPoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional as this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNIC(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkInterfaceBackendAddressPoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
	if err != nil {
		return nil, err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	read, err := clients.Network.Client.NetworkInterfaces.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if read.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}
	if read.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if read.Model.Properties.IPConfigurations == nil {
		return nil, fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}

	config := network.FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return nil, fmt.Errorf("IP Configuration %q wasn't found for %s", id.First.IpConfigurationName, networkInterfaceId)
	}

	found := false
	if config.Properties.LoadBalancerBackendAddressPools != nil {
		for _, pool := range *config.Properties.LoadBalancerBackendAddressPools {
			if *pool.Id == id.Second.ID() {
				found = true
				break
			}
		}
	}

	return pointer.To(found), nil
}

func (NetworkInterfaceBackendAddressPoolResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := commonids.ParseCompositeResourceID(state.ID, &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	read, err := client.Network.Client.NetworkInterfaces.Get(ctx2, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	config := network.FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q wasn't found for %s", id.First.IpConfigurationName, networkInterfaceId)
	}

	updatedPools := make([]networkinterfaces.BackendAddressPool, 0)
	if config.Properties.LoadBalancerBackendAddressPools != nil {
		for _, pool := range *config.Properties.LoadBalancerBackendAddressPools {
			if *pool.Id != id.Second.ID() {
				updatedPools = append(updatedPools, pool)
			}
		}
	}
	config.Properties.LoadBalancerBackendAddressPools = &updatedPools

	if err := client.Network.Client.NetworkInterfaces.CreateOrUpdateThenPoll(ctx2, networkInterfaceId, *read.Model); err != nil {
		return fmt.Errorf("removing Backend Address Pool Association for %s: %+v", networkInterfaceId, err)
	}

	return nil
}

func (r NetworkInterfaceBackendAddressPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = azurerm_network_interface.test.id
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceBackendAddressPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_backend_address_pool_association" "import" {
  network_interface_id    = azurerm_network_interface_backend_address_pool_association.test.network_interface_id
  ip_configuration_name   = azurerm_network_interface_backend_address_pool_association.test.ip_configuration_name
  backend_address_pool_id = azurerm_network_interface_backend_address_pool_association.test.backend_address_pool_id
}
`, r.basic(data))
}

func (r NetworkInterfaceBackendAddressPoolResource) updateNIC(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    private_ip_address_version    = "IPv6"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = azurerm_network_interface.test.id
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceBackendAddressPoolResource) template(data acceptance.TestData) string {
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
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "primary"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "acctestpool"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
