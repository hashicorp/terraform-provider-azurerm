package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ types.TestResourceVerifyingRemoved = BackendAddressPoolAddressResourceTests{}

type BackendAddressPoolAddressResourceTests struct{}

func TestAccBackendAddressPoolAddressBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolAddressRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackendAddressPoolAddressDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolAddressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool_address", "test")
	r := BackendAddressPoolAddressResourceTests{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (BackendAddressPoolAddressResourceTests) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	pool, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
			if address.Name == nil {
				continue
			}

			if *address.Name == id.AddressName {
				return utils.Bool(true), nil
			}
		}
	}
	return utils.Bool(false), nil
}

func (BackendAddressPoolAddressResourceTests) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.BackendAddressPoolAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	pool, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	addresses := make([]network.LoadBalancerBackendAddress, 0)
	if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
		addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
	}
	newAddresses := make([]network.LoadBalancerBackendAddress, 0)
	for _, address := range addresses {
		if address.Name == nil {
			continue
		}

		if *address.Name != id.AddressName {
			newAddresses = append(newAddresses, address)
		}
	}
	pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &newAddresses

	future, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (BackendAddressPoolAddressResourceTests) backendAddressPoolHasAddresses(expected int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		id, err := parse.LoadBalancerBackendAddressPoolID(state.ID)
		if err != nil {
			return err
		}

		client := clients.LoadBalancers.LoadBalancerBackendAddressPoolsClient
		pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			return err
		}
		if pool.BackendAddressPoolPropertiesFormat == nil {
			return fmt.Errorf("`properties` is nil")
		}
		if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses == nil {
			return fmt.Errorf("`properties.loadBalancerBackendAddresses` is nil")
		}

		actual := len(*pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses)
		if actual != expected {
			return fmt.Errorf("expected %d but got %d addresses", expected, actual)
		}

		return nil
	}
}

func (t BackendAddressPoolAddressResourceTests) basic(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.test.id
  ip_address              = "191.168.0.1"
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
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "second" {
  name                = "acctestvn-2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  virtual_network_id      = azurerm_virtual_network.second.id
  ip_address              = "10.0.0.1"
}
`, template, data.RandomInteger)
}

func (BackendAddressPoolAddressResourceTests) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
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
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
