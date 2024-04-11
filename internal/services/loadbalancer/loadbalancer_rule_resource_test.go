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

type LoadBalancerRule struct{}

func TestAccAzureRMLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

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

func TestAccAzureRMLoadBalancerRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

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

func TestAccAzureRMLoadBalancerRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

// https://github.com/hashicorp/terraform/issues/9424
func TestAccAzureRMLoadBalancerRule_inconsistentReads(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}
	p := LoadBalancerProbe{}
	b := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inconsistentRead(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_lb_probe.test").ExistsInAzure(p),
				check.That("azurerm_lb_backend_address_pool.test").ExistsInAzure(b),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_updateMultipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_rule", "test2")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRules(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("frontend_port").HasValue("3390"),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3390"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleRulesUpdate(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("frontend_port").HasValue("3391"),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3391"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_vmssBackendPoolUpdateRemoveLBRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	lbRuleName := fmt.Sprintf("LbRule-%s", data.RandomString)
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vmssBackendPool(data, lbRuleName, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vmssBackendPoolUpdate(data, lbRuleName, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vmssBackendPoolWithoutLBRule(data, "Standard"),
		},
	})
}

func TestAccAzureRMLoadBalancerRule_gatewayLBRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewayLBRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_gatewayLBRuleMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewayLBRuleMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerRule) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseLoadBalancingRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	rule, err := client.LoadBalancers.LoadBalancersClient.LoadBalancerLoadBalancingRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(rule.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(rule.Model != nil && rule.Model.Id != nil), nil
}

func (r LoadBalancerRule) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseLoadBalancingRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	if loadBalancer.Model == nil {
		return nil, fmt.Errorf(`model was nil`)
	}
	if loadBalancer.Model.Properties == nil {
		return nil, fmt.Errorf(`properties was nil`)
	}
	if loadBalancer.Model.Properties.LoadBalancingRules == nil {
		return nil, fmt.Errorf(`properties.LoadBalancingRules was nil`)
	}
	rules := make([]loadbalancers.LoadBalancingRule, 0)
	for _, v := range *loadBalancer.Model.Properties.LoadBalancingRules {
		if v.Name == nil || *v.Name == id.LoadBalancingRuleName {
			continue
		}

		rules = append(rules, v)
	}
	loadBalancer.Model.Properties.LoadBalancingRules = &rules

	err = client.LoadBalancers.LoadBalancersClient.CreateOrUpdateThenPoll(ctx, plbId, *loadBalancer.Model)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r LoadBalancerRule) template(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lb-%[1]d"
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
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (r LoadBalancerRule) basic(data acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  name                           = "LbRule-%s"
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) complete(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  name            = "LbRule-%s"
  loadbalancer_id = azurerm_lb.test.id

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 100
  load_distribution       = "SourceIP"

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "import" {
  name                           = azurerm_lb_rule.test.name
  loadbalancer_id                = azurerm_lb_rule.test.loadbalancer_id
  frontend_ip_configuration_name = azurerm_lb_rule.test.frontend_ip_configuration_name
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template)
}

// https://github.com/hashicorp/terraform/issues/9424
func (r LoadBalancerRule) inconsistentRead(data acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "%d-address-pool"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  name            = "probe-%d"
  loadbalancer_id = azurerm_lb.test.id
  protocol        = "Tcp"
  port            = 443
}

resource "azurerm_lb_rule" "test" {
  name                           = "LbRule-%s"
  loadbalancer_id                = azurerm_lb.test.id
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) multipleRules(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8), data2.RandomStringOfLength(8))
}

func (r LoadBalancerRule) multipleRulesUpdate(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8), data2.RandomStringOfLength(8))
}

func (r LoadBalancerRule) vmssBackendPoolWithoutLBRule(data acceptance.TestData, sku string) string {
	template := r.template(data, sku)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-lb-BAP-%[2]d"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-lb-vnet-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-lb-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctest-lb-vmss-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [
        azurerm_lb_backend_address_pool.test.id
      ]
    }
  }
}
`, template, data.RandomInteger)
}

func (r LoadBalancerRule) vmssBackendPool(data acceptance.TestData, lbRuleName, sku string) string {
	template := r.vmssBackendPoolWithoutLBRule(data, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, lbRuleName)
}

func (r LoadBalancerRule) vmssBackendPoolUpdate(data acceptance.TestData, lbRuleName, sku string) string {
	template := r.vmssBackendPoolWithoutLBRule(data, sku)
	return fmt.Sprintf(`
%s
resource "azurerm_lb_rule" "test" {
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.test.id]
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
  disable_outbound_snat          = false
}
`, template, lbRuleName)
}

func (r LoadBalancerRule) gatewayLBRuleTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]d"
  resource_group_name  = azurerm_virtual_network.test.resource_group_name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Gateway"

  frontend_ip_configuration {
    name      = "feip"
    subnet_id = azurerm_subnet.test.id
  }
}

resource "azurerm_public_ip" "c1" {
  name                = "acctestpip1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_lb" "c1" {
  name                = "acctestlb-consumer1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  frontend_ip_configuration {
    name                                               = "gateway"
    public_ip_address_id                               = azurerm_public_ip.c1.id
    gateway_load_balancer_frontend_ip_configuration_id = azurerm_lb.test.frontend_ip_configuration.0.id
  }
}

resource "azurerm_public_ip" "c2" {
  name                = "acctestpip2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_lb" "c2" {
  name                = "acctestlb-consumer2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  frontend_ip_configuration {
    name                                               = "gateway"
    public_ip_address_id                               = azurerm_public_ip.c2.id
    gateway_load_balancer_frontend_ip_configuration_id = azurerm_lb.test.frontend_ip_configuration.0.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LoadBalancerRule) gatewayLBRule(data acceptance.TestData) string {
	template := r.gatewayLBRuleTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctestbap-%[2]d"
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
}

resource "azurerm_lb_rule" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "abababa"
  protocol        = "All"
  frontend_port   = 0
  backend_port    = 0
  backend_address_pool_ids = [
    azurerm_lb_backend_address_pool.test.id,
  ]
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger)
}

func (r LoadBalancerRule) gatewayLBRuleMultiple(data acceptance.TestData) string {
	template := r.gatewayLBRuleTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_lb_backend_address_pool" "test1" {
  name            = "acctestbap1-%[2]d"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
}

resource "azurerm_lb_backend_address_pool" "test2" {
  name            = "acctestbap2-%[2]d"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 901
    type       = "External"
    protocol   = "VXLAN"
    port       = 15001
  }
}

resource "azurerm_lb_rule" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "abababa"
  protocol        = "All"
  frontend_port   = 0
  backend_port    = 0
  backend_address_pool_ids = [
    azurerm_lb_backend_address_pool.test1.id,
    azurerm_lb_backend_address_pool.test2.id,
  ]
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger)
}
