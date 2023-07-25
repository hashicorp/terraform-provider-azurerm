// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type LoadBalancerBackendAddressPool struct{}

// Basic and Standard use different API's for reasons, so we need to test both flows

func TestAccBackendAddressPoolBasicSkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolBasicSkuDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicSkuBasic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolBasicSkuRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.basicSkuRequiresImport),
	})
}

func TestAccBackendAddressPoolStandardSkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolStandardSkuDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.standardSkuBasic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolStandardSkuRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.standardSkuRequiresImport),
	})
}

func TestAccBackendAddressPool_GatewaySkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPool_GatewaySkuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewaySkuComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerBackendAddressPool) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r LoadBalancerBackendAddressPool) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "pool"
  loadbalancer_id = azurerm_lb.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) basicSkuRequiresImport(data acceptance.TestData) string {
	template := r.basicSkuBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuBasic(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "pool"
  loadbalancer_id = azurerm_lb.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuBasicUpdate(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name               = "pool"
  loadbalancer_id    = azurerm_lb.test.id
  virtual_network_id = azurerm_virtual_network.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuRequiresImport(data acceptance.TestData) string {
	template := r.standardSkuBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func (LoadBalancerBackendAddressPool) template(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
locals {
  number   = %d
  location = %q
  sku      = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = local.sku
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = local.sku

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (r LoadBalancerBackendAddressPool) gatewaySkuBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-bap-${local.number}"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
}
`, r.templateGateway(data))
}

func (r LoadBalancerBackendAddressPool) gatewaySkuComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-bap-${local.number}"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
  tunnel_interface {
    identifier = 901
    type       = "External"
    protocol   = "VXLAN"
    port       = 15001
  }
  virtual_network_id = azurerm_virtual_network.test.id
}
`, r.templateGateway(data))
}

func (LoadBalancerBackendAddressPool) templateGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  number   = %d
  location = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-${local.number}"
  resource_group_name  = azurerm_virtual_network.test.resource_group_name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Gateway"

  frontend_ip_configuration {
    name      = "feip"
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
