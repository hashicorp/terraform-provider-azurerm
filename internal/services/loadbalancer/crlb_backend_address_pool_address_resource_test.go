package loadbalancer_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CrossRegionLoadBalancerBackendAddressPoolAddressResource struct{}

func TestAccCrossRegionLoadBalancerBackendAddressPoolAddress_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_crlb_backend_address_pool_address", "test")
	r := CrossRegionLoadBalancerBackendAddressPoolAddressResource{}
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

func TestAccCrossRegionLoadBalancerBackendAddressPoolAddress_MultipleAddresses(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_crlb_backend_address_pool_address", "test1")
	r := CrossRegionLoadBalancerBackendAddressPoolAddressResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_crlb_backend_address_pool_address.test1").ExistsInAzure(r),
				check.That("azurerm_crlb_backend_address_pool_address.test2").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStepFor("azurerm_crlb_backend_address_pool_address.test2"),
	})
}

func TestAccCrossRegionLoadBalancerBackendAddressPoolAddress_updateMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_crlb_backend_address_pool_address", "test1")
	r := CrossRegionLoadBalancerBackendAddressPoolAddressResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_crlb_backend_address_pool_address.test1").ExistsInAzure(r),
				check.That("azurerm_crlb_backend_address_pool_address.test2").ExistsInAzure(r),
			),
		},
		data.ImportStepFor("azurerm_crlb_backend_address_pool_address.test1"),
		data.ImportStepFor("azurerm_crlb_backend_address_pool_address.test2"),
		{
			Config: r.updateMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_crlb_backend_address_pool_address.test1").ExistsInAzure(r),
				check.That("azurerm_crlb_backend_address_pool_address.test2").ExistsInAzure(r),
			),
		},
		data.ImportStepFor("azurerm_crlb_backend_address_pool_address.test1"),
		data.ImportStepFor("azurerm_crlb_backend_address_pool_address.test2"),
	})
}

func TestAccCrossRegionLoadBalancerBackendAddressPoolAddress_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_crlb_backend_address_pool_address", "test")
	r := CrossRegionLoadBalancerBackendAddressPoolAddressResource{}
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

func (CrossRegionLoadBalancerBackendAddressPoolAddressResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) basic(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_crlb_backend_address_pool_address" "test" {
  name                                = "myBackendPoolConfig-R1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_address          = azurerm_public_ip.backend-ip-R1.ip_address
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[0].id
}

`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) multiple(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_crlb_backend_address_pool_address" "test1" {
  name                                = "myBackendPoolConfig-R1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_address          = azurerm_public_ip.backend-ip-R1.ip_address
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[0].id
}

resource "azurerm_crlb_backend_address_pool_address" "test2" {
  name                                = "myBackendPoolConfig-R2"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_address          = azurerm_public_ip.backend-ip-R2.ip_address
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R2.frontend_ip_configuration[0].id
}
`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_crlb_backend_address_pool_address" "import" {
  name                                = azurerm_crlb_backend_address_pool_address.test.name
  backend_address_pool_id             = azurerm_crlb_backend_address_pool_address.test.backend_address_pool_id
  backend_address_ip_address          = azurerm_crlb_backend_address_pool_address.test.backend_address_ip_address
  backend_address_ip_configuration_id = azurerm_crlb_backend_address_pool_address.test.backend_address_ip_configuration_id
}

`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) updateMultiple(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s
resource "azurerm_crlb_backend_address_pool_address" "test1" {
  name                                = "myBackendPoolConfig-R1"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_address          = azurerm_public_ip.backend-ip-R1-1.ip_address
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[1].id
}

resource "azurerm_crlb_backend_address_pool_address" "test2" {
  name                                = "myBackendPoolConfig-R2"
  backend_address_pool_id             = azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_address          = azurerm_public_ip.backend-ip-R2.ip_address
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R2.frontend_ip_configuration[0].id
}
`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
  name                = "acctestpip1-%d-R1"
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
resource "azurerm_lb_backend_address_pool" "backend-pool-R1" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.backend-lb-R1.id
}
resource "azurerm_lb_backend_address_pool_address" "backend-address-R1" {
  name                    = "address1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.backend-pool-R1.id
  virtual_network_id      = azurerm_virtual_network.backend-vn-R1.id
  ip_address              = "192.168.0.0"
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
resource "azurerm_lb_backend_address_pool" "backend-pool-R2" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.backend-lb-R2.id
}
resource "azurerm_lb_backend_address_pool_address" "backend-address-R2" {
  name                    = "address1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.backend-pool-R2.id
  virtual_network_id      = azurerm_virtual_network.backend-vn-R2.id
  ip_address              = "10.0.0.0"
}
resource "azurerm_lb" "lb-cr" {
  name                = "acctestlb-%d-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Global"
  frontend_ip_configuration {
    name                 = "one"
    public_ip_address_id = azurerm_public_ip.backend-ip-cr.id
  }
}
resource "azurerm_lb_backend_address_pool" "backend-pool-cr" {
  loadbalancer_id = azurerm_lb.lb-cr.id
  name            = "myBackendPool-cr"
}`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
