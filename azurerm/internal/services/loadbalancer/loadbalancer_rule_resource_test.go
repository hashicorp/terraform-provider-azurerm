package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LoadBalancerRule struct {
}

func TestAccAzureRMLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

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

func TestAccAzureRMLoadBalancerRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

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

func TestAccAzureRMLoadBalancerRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.inconsistentRead(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleRules(data, data2),
			Check: resource.ComposeTestCheckFunc(
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
			Check: resource.ComposeTestCheckFunc(
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
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vmssBackendPool(data, lbRuleName, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vmssBackendPoolUpdate(data, lbRuleName, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.vmssBackendPoolWithoutLBRule(data, "Standard"),
		},
	})
}

func (r LoadBalancerRule) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancingRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	rule, err := client.LoadBalancers.LoadBalancingRulesClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(rule.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(rule.ID != nil), nil
}

func (r LoadBalancerRule) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancingRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	loadBalancer, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if loadBalancer.LoadBalancerPropertiesFormat == nil {
		return nil, fmt.Errorf(`properties was nil`)
	}
	if loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules == nil {
		return nil, fmt.Errorf(`properties.LoadBalancingRules was nil`)
	}
	rules := make([]network.LoadBalancingRule, 0)
	for _, v := range *loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules {
		if v.Name == nil || *v.Name == id.Name {
			continue
		}

		rules = append(rules, v)
	}
	loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = &rules

	future, err := client.LoadBalancers.LoadBalancersClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, loadBalancer)
	if err != nil {
		return nil, fmt.Errorf("updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.LoadBalancers.LoadBalancersClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
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
  resource_group_name            = azurerm_resource_group.test.name
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
  name                = "LbRule-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 10

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
  resource_group_name            = azurerm_lb_rule.test.resource_group_name
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
  name                = "%d-address-pool"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_probe" "test" {
  name                = "probe-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  protocol            = "Tcp"
  port                = 443
}

resource "azurerm_lb_rule" "test" {
  name                           = "LbRule-%s"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
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
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
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
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
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
  name                = "acctest-lb-BAP-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
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
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, lbRuleName)
}

func (r LoadBalancerRule) vmssBackendPoolUpdate(data acceptance.TestData, lbRuleName, sku string) string {
	template := r.vmssBackendPoolWithoutLBRule(data, sku)
	return fmt.Sprintf(`
%s
resource "azurerm_lb_rule" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
  disable_outbound_snat          = false
}
`, template, lbRuleName)
}
