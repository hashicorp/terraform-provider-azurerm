package loadbalancer_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
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
				check.That(data.ResourceName).ExistsInAzure()),
		},
		data.ImportStep(),
	})
}

func TestAccCrossRegionLoadBalancerBackendAddressPoolAddress_update(t *testing.T) {
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
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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

// todo finish the Exists function
func (CrossRegionLoadBalancerBackendAddressPoolAddressResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	poolId, err := parse.LoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	pool, err := client.LoadBalancers.LoadBalancerBackendAddressPoolsClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *poolId, err)
	}
	if pool.BackendAddressPoolPropertiesFormat == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
	}

	if state.Attributes[""]
	return utils.Bool(false), nil
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) basic(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_crlb_backend_address_pool_address" "crbackendaddress"{
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  backend_addresses {
    load_balancer {
      lb_name = azurerm_lb.test1.name
      lb_ip_address = azurerm_public_ip.test1.ip_address
      lb_frontend_ip_configuration_id = azurerm_lb.test1.frontend_ip_configuration[0].id
    }

    load_balancer {
      lb_name = azurerm_lb.test2.name
      lb_ip_address = azurerm_public_ip.test2.ip_address
      lb_frontend_ip_configuration_id = azurerm_lb.test2.frontend_ip_configuration[0].id
    }
  }
}
`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_crlb_backend_address_pool_address" "import" {
  backend_address_pool_id = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_address_pool_id
  backend_addresses {
    load_balancer {
      lb_name                         = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[0].lb_name
      lb_ip_address                   = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[0].lb_ip_address
      lb_frontend_ip_configuration_id = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[0].lb_frontend_ip_configuration_id
    }

    load_balancer {
      lb_name                         = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[1].lb_name
      lb_ip_address                   = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[1].lb_ip_address
      lb_frontend_ip_configuration_id = azurerm_crlb_backend_address_pool_address.crbackendpool.backend_addresses.load_balancer[1].lb_frontend_ip_configuration_id
    }
  }
}
`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) update(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_crlb_backend_address_pool_address" "crbackendaddress"{
  backend_address_pool_id = azurerm_lb_backend_address_pool.crbackendpool.id
  backend_addresses {
    load_balancer {
      lb_name = azurerm_lb.test1.name
      lb_ip_address = azurerm_public_ip.test11.ip_address
      lb_frontend_ip_configuration_id = azurerm_lb.test1.frontend_ip_configuration[1].id
    }

    load_balancer {
      lb_name = azurerm_lb.test2.name
      lb_ip_address = azurerm_public_ip.test2.ip_address
      lb_frontend_ip_configuration_id = azurerm_lb.test2.frontend_ip_configuration[0].id
    }
  }
}`, template)
}

func (t CrossRegionLoadBalancerBackendAddressPoolAddressResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvn-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvn-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.1.0/16"]
}

resource "azurerm_public_ip" "test1" {
  name                = "acctestpip-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test11" {
  name                = "acctestpip1-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "acctestpip-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "testip-cr" {
  name                = "acctestpip-%d-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test1" {
  name                = "acctestlb-%d-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "feip1"
    public_ip_address_id = azurerm_public_ip.test11.id
  }
}

resource "azurerm_lb_backend_address_pool" "test1" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test1.id
}

resource "azurerm_lb" "test2" {
  name                = "acctestlb-%d-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test2" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test2.id
}

resource "azurerm_lb" "testcr" {
  name                = "acctestlb-%d-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Global"

  frontend_ip_configuration {
    name                 = "one"
    public_ip_address_id = azurerm_public_ip.testip-cr.id
  }
}

resource "azurerm_lb_backend_address_pool" "crbackendpool" {
  loadbalancer_id = azurerm_lb.testcr.id
  name            = "myBackendPool-cr"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
