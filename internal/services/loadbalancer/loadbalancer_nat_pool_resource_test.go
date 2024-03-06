// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LoadBalancerNatPool struct{}

func TestAccAzureRMLoadBalancerNatPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerNatPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerNatPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerNatPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMLoadBalancerNatPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config: func(data acceptance.TestData) string {
				return r.basic(data, "Basic")
			},
			TestResource: r,
		}),
	})
}

func TestAccAzureRMLoadBalancerNatPool_updateMultiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test2")

	r := LoadBalancerNatPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePools(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3390"),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiplePoolsUpdate(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3391"),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerNatPool) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerInboundNatPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", plbId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	model := lb.Model
	if model == nil || model.Properties == nil || model.Properties.InboundNatPools == nil || len(*model.Properties.InboundNatPools) == 0 {
		return nil, fmt.Errorf("Nat Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatPoolName, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *model.Properties.InboundNatPools {
		if v.Name != nil && *v.Name == id.InboundNatPoolName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Nat Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatPoolName, id.LoadBalancerName, id.ResourceGroup)
	}

	return pointer.To(found), nil
}

func (r LoadBalancerNatPool) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerInboundNatPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	if lb.Model == nil {
		return nil, fmt.Errorf("`model` was nil")
	}
	if lb.Model.Properties == nil {
		return nil, fmt.Errorf("`properties` was nil")
	}
	if lb.Model.Properties.InboundNatPools == nil {
		return nil, fmt.Errorf("`properties.InboundNatPools` was nil")
	}

	inboundNatPools := make([]loadbalancers.InboundNatPool, 0)
	for _, inboundNatPool := range *lb.Model.Properties.InboundNatPools {
		if inboundNatPool.Name == nil || *inboundNatPool.Name == id.InboundNatPoolName {
			continue
		}

		inboundNatPools = append(inboundNatPools, inboundNatPool)
	}
	lb.Model.Properties.InboundNatPools = &inboundNatPools

	err = client.LoadBalancers.LoadBalancersClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
	if err != nil {
		return nil, fmt.Errorf("updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}
	return pointer.To(true), nil
}

func (r LoadBalancerNatPool) basic(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "%[3]s"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "%[3]s"

  frontend_ip_configuration {
    name                 = "one-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "NatPool-%[1]d"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (r LoadBalancerNatPool) complete(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
resource "azurerm_public_ip" "test" {
  name                = "test-ip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "%[3]s"
}
resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "%[3]s"

  frontend_ip_configuration {
    name                 = "one-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "NatPool-%[1]d"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%[1]d"
  floating_ip_enabled            = true
  tcp_reset_enabled              = true
  idle_timeout_in_minutes        = 10
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (r LoadBalancerNatPool) requiresImport(data acceptance.TestData) string {
	template := r.basic(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_pool" "import" {
  name                           = azurerm_lb_nat_pool.test.name
  loadbalancer_id                = azurerm_lb_nat_pool.test.loadbalancer_id
  resource_group_name            = azurerm_lb_nat_pool.test.resource_group_name
  frontend_ip_configuration_name = azurerm_lb_nat_pool.test.frontend_ip_configuration_name
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
}
`, template)
}

func (r LoadBalancerNatPool) multiplePools(data, data2 acceptance.TestData) string {
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

  allocation_method = "Static"
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

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "NatPool-%d"
  protocol            = "Tcp"
  frontend_port_start = 80
  frontend_port_end   = 81
  backend_port        = 3389

  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "NatPool-%d"
  protocol            = "Tcp"
  frontend_port_start = 82
  frontend_port_end   = 83
  backend_port        = 3390

  frontend_ip_configuration_name = "one-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerNatPool) multiplePoolsUpdate(data, data2 acceptance.TestData) string {
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

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "NatPool-%d"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "NatPool-%d"
  protocol                       = "Tcp"
  frontend_port_start            = 82
  frontend_port_end              = 83
  backend_port                   = 3391
  frontend_ip_configuration_name = "one-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger, data.RandomInteger)
}
