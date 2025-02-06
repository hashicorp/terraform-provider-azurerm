// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ types.TestResourceVerifyingRemoved = BackendAddressPoolAddressResourceTests{}

type BackendAddressPoolAddressResourceTests struct{}

func TestAccBackendAddressPoolAddress_regionalLbBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
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

func TestAccBackendAddressPoolAddress_regionalLbRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
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

func TestAccBackendAddressPoolAddress_regionalLbDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolAddress_regionalLbUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func TestAccBackendAddressPoolAddress_globalLbUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test1")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionLoadBalancer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.crossRegionLoadBalancerUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolAddress_globalLbRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test1")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionLoadBalancer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.crossRegionLoadBalancerRemoval(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (BackendAddressPoolAddressResourceTests) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	pool, err := client.LoadBalancers.LoadBalancersClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if pool.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if backendAddress := pool.Model.Properties.LoadBalancerBackendAddresses; backendAddress != nil {
		for _, address := range *backendAddress {
			if address.Name == nil {
				continue
			}

			if *address.Name == id.AddressName {
				return pointer.To(true), nil
			}
		}
	}
	return pointer.To(false), nil
}

func (BackendAddressPoolAddressResourceTests) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	pool, err := client.LoadBalancers.LoadBalancersClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if pool.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
	if pool.Model.Properties.LoadBalancerBackendAddresses != nil {
		addresses = *pool.Model.Properties.LoadBalancerBackendAddresses
	}
	newAddresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
	for _, address := range addresses {
		if address.Name == nil {
			continue
		}

		if *address.Name != id.AddressName {
			newAddresses = append(newAddresses, address)
		}
	}

	pool.Model.Properties.LoadBalancerBackendAddresses = &newAddresses

	err = client.LoadBalancers.LoadBalancersClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, poolId, *pool.Model)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

// nolint unused - for future use
func (BackendAddressPoolAddressResourceTests) backendAddressPoolHasAddresses(expected int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(state.ID)
		if err != nil {
			return err
		}

		pool, err := clients.LoadBalancers.LoadBalancersClient.LoadBalancerBackendAddressPoolsGet(ctx, *id)
		if err != nil {
			return err
		}
		if pool.Model == nil {
			return fmt.Errorf("`model` is nil")
		}
		if pool.Model.Properties == nil {
			return fmt.Errorf("`properties` is nil")
		}
		if pool.Model.Properties.LoadBalancerBackendAddresses == nil {
			return fmt.Errorf("`properties.loadBalancerBackendAddresses` is nil")
		}

		actual := len(*pool.Model.Properties.LoadBalancerBackendAddresses)
		if actual != expected {
			return fmt.Errorf("expected %d but got %d addresses", expected, actual)
		}

		return nil
	}
}

func (t BackendAddressPoolAddressResourceTests) basic(data acceptance.TestData) string {
	template := t.templateRegionalLB(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.test.id
  ip_address              = "191.168.0.1"
  depends_on              = [azurerm_lb_backend_address_pool.test]
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool_address" "import" {
  name                    = azurerm_lb_backend_address_pool_address.test.name
  backend_address_pool_id = azurerm_lb_backend_address_pool_address.test.backend_address_pool_id
  virtual_network_id      = azurerm_lb_backend_address_pool_address.test.virtual_network_id
  ip_address              = azurerm_lb_backend_address_pool_address.test.ip_address
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) update(data acceptance.TestData) string {
	template := t.templateRegionalLB(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.test.id
  ip_address              = "191.168.0.2"
  depends_on              = [azurerm_lb_backend_address_pool.test]
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) crossRegionLoadBalancer(data acceptance.TestData) string {
	template := t.templateGlobalLB(data)
	return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool_address" "test1" {
  name                                = "address1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[0].id
}

resource "azurerm_lb_backend_address_pool_address" "test2" {
  name                                = "address2"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R2.frontend_ip_configuration[0].id
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) crossRegionLoadBalancerUpdate(data acceptance.TestData) string {
	template := t.templateGlobalLB(data)
	return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool_address" "test1" {
  name                                = "address1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[1].id
}

resource "azurerm_lb_backend_address_pool_address" "test2" {
  name                                = "address2"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R2.frontend_ip_configuration[0].id
}
`, template)
}

func (t BackendAddressPoolAddressResourceTests) crossRegionLoadBalancerRemoval(data acceptance.TestData) string {
	template := t.templateGlobalLB(data)
	return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool_address" "test1" {
  name                                = "address1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[1].id
}


`, template)
}

func (BackendAddressPoolAddressResourceTests) templateGlobalLB(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "backend-vn-R1" {
  name                = "acctestvn-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_virtual_network" "backend-vn-R2" {
  name                = "acctestvn-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_public_ip" "backend-ip-R1" {
  name                = "acctestpip-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "backend-ip-R1-1" {
  name                = "acctestpip-%d-R1-1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "backend-ip-R2" {
  name                = "acctestpip-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "backend-ip-cr" {
  name                = "acctestpip-%d-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  sku_tier            = "Global"
}

resource "azurerm_lb" "backend-lb-R1" {
  name                = "acctestlb-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.backend-ip-R1.id
  }
  frontend_ip_configuration {
    name                 = "feip1"
    public_ip_address_id = azurerm_public_ip.backend-ip-R1-1.id
  }
}

resource "azurerm_lb" "backend-lb-R2" {
  name                = "acctestlb-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.backend-ip-R2.id
  }
}

resource "azurerm_lb" "backend-lb-cr" {
  name                = "acctestlb-%d-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Global"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.backend-ip-cr.id
  }
}

resource "azurerm_lb_backend_address_pool" "backend-pool-R1" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.backend-lb-R1.id
}

resource "azurerm_lb_backend_address_pool" "backend-pool-R2" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.backend-lb-R2.id
}

resource "azurerm_lb_backend_address_pool" "backend-pool-cr" {
  name            = "myBackendPool-cr"
  loadbalancer_id = azurerm_lb.backend-lb-cr.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (BackendAddressPoolAddressResourceTests) templateRegionalLB(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test.id
  }
  depends_on = [azurerm_public_ip.test]
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
  depends_on      = [azurerm_lb.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
