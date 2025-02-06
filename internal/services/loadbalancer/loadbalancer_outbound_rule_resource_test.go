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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LoadBalancerOutboundRule struct{}

func TestAccAzureRMLoadBalancerOutboundRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

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

func TestAccAzureRMLoadBalancerOutboundRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

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

func TestAccAzureRMLoadBalancerOutboundRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test2")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRules(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleRulesUpdate(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPublicIPPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_allocatedOutboundPorts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allocatedOutboundPortsDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.allocatedOutboundPortsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerOutboundRule) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseOutboundRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", plbId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	props := lb.Model
	if props == nil || props.Properties == nil || len(*props.Properties.OutboundRules) == 0 {
		return nil, fmt.Errorf("Outbound Rule %q not found in Load Balancer %q (resource group %q)", id.OutboundRuleName, id.LoadBalancerName, id.ResourceGroupName)
	}

	found := false
	for _, v := range *props.Properties.OutboundRules {
		if v.Name != nil && *v.Name == id.OutboundRuleName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Outbound Rule %q not found in Load Balancer %q (resource group %q)", id.OutboundRuleName, id.LoadBalancerName, id.ResourceGroupName)
	}
	return pointer.To(found), nil
}

func (r LoadBalancerOutboundRule) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseOutboundRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
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
	if lb.Model.Properties.OutboundRules == nil {
		return nil, fmt.Errorf("`properties.OutboundRules` was nil")
	}

	outboundRules := make([]loadbalancers.OutboundRule, 0)
	for _, outboundRule := range *lb.Model.Properties.OutboundRules {
		if outboundRule.Name == nil || *outboundRule.Name == id.OutboundRuleName {
			continue
		}

		outboundRules = append(outboundRules, outboundRule)
	}
	lb.Model.Properties.OutboundRules = &outboundRules

	err = client.LoadBalancers.LoadBalancersClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r LoadBalancerOutboundRule) basic(data acceptance.TestData) string {
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

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerOutboundRule) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_outbound_rule" "import" {
  name                    = azurerm_lb_outbound_rule.test.name
  loadbalancer_id         = azurerm_lb_outbound_rule.test.loadbalancer_id
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = azurerm_lb_outbound_rule.test.frontend_ip_configuration[0].name
  }
}
`, template)
}

func (r LoadBalancerOutboundRule) multipleRules(data, data2 acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
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
    name                 = "fe1-%d"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "Udp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerOutboundRule) multipleRulesUpdate(data, data2 acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
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
    name                 = "fe1-%d"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerOutboundRule) withPublicIPPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 31
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                = "one-%d"
    public_ip_prefix_id = azurerm_public_ip_prefix.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerOutboundRule) allocatedOutboundPortsDefault(data acceptance.TestData) string {
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

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerOutboundRule) allocatedOutboundPortsUpdated(data acceptance.TestData) string {
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

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  allocated_outbound_ports = 0

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
