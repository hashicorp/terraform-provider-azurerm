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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LoadBalancerBackendAddressPool struct {
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_basicSkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_standardSkuNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardSkuNIC(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_standardSkuIPBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardSkuIPBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_standardSkuIPUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardSkuIPBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuIPMultipleBackendAddress(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuIPChangeVnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_BasicSkuRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicSkuBasic,
			TestResource: r,
		}),
	})
}

func (r LoadBalancerBackendAddressPool) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return nil, fmt.Errorf("Load Balancer %q (resource group %q) not found for Backend Address Pool %q", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName)
		}
		return nil, fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Backend Address Pool %q", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName)
	}
	props := lb.LoadBalancerPropertiesFormat
	if props == nil || props.BackendAddressPools == nil || len(*props.BackendAddressPools) == 0 {
		return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *props.BackendAddressPools {
		if v.Name != nil && *v.Name == id.BackendAddressPoolName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroup)
	}
	return utils.Bool(true), nil
}

func (r LoadBalancerBackendAddressPool) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Load Balancer %q (Resource Group %q)", id.LoadBalancerName, id.ResourceGroup)
	}
	if lb.LoadBalancerPropertiesFormat == nil {
		return nil, fmt.Errorf("`properties` was nil")
	}
	if lb.LoadBalancerPropertiesFormat.BackendAddressPools == nil {
		return nil, fmt.Errorf("`properties.BackendAddressPools` was nil")
	}

	backendAddressPools := make([]network.BackendAddressPool, 0)
	for _, backendAddressPool := range *lb.LoadBalancerPropertiesFormat.BackendAddressPools {
		if backendAddressPool.Name == nil || *backendAddressPool.Name == id.BackendAddressPoolName {
			continue
		}

		backendAddressPools = append(backendAddressPools, backendAddressPool)
	}
	lb.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

	future, err := client.LoadBalancers.LoadBalancersClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, lb)
	if err != nil {
		return nil, fmt.Errorf("updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.LoadBalancers.LoadBalancersClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r LoadBalancerBackendAddressPool) basicSkuBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "Address-pool-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) requiresImport(data acceptance.TestData) string {
	template := r.basicSkuBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) standardSkuNIC(data acceptance.TestData) string {
	template := r.standardSkuTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "Address-pool-%d"
}
`, template, data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) standardSkuIPBasic(data acceptance.TestData) string {
	template := r.standardSkuTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "arm-test-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "Address-pool-%d"
  backend_address {
    name               = "addr1"
    virtual_network_id = azurerm_virtual_network.test.id
    ip_address         = "191.168.0.1"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) standardSkuIPMultipleBackendAddress(data acceptance.TestData) string {
	template := r.standardSkuTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "arm-test-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "Address-pool-%d"
  backend_address {
    name               = "addr1"
    virtual_network_id = azurerm_virtual_network.test.id
    ip_address         = "191.168.0.1"
  }

  backend_address {
    name               = "addr2"
    virtual_network_id = azurerm_virtual_network.test.id
    ip_address         = "191.168.0.2"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) standardSkuIPChangeVnet(data acceptance.TestData) string {
	template := r.standardSkuTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "arm-test-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_virtual_network" "test2" {
  name                = "arm-test-vnet2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "Address-pool-%d"
  backend_address {
    name               = "addr1"
    virtual_network_id = azurerm_virtual_network.test2.id
    ip_address         = "10.0.0.1"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
