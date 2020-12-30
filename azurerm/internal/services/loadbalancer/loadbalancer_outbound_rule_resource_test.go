package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LoadBalancerOutboundRule struct {
}

func TestAccAzureRMLoadBalancerOutboundRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

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

func TestAccAzureRMLoadBalancerOutboundRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

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

func TestAccAzureRMLoadBalancerOutboundRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removal(data),
			Check: resource.ComposeTestCheckFunc(
				r.IsMissing("azurerm_lb.test", fmt.Sprintf("OutboundRule-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test2")
	r := LoadBalancerOutboundRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleRules(data, data2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleRulesUpdate(data, data2),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withPublicIPPrefix(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerOutboundRule) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerOutboundRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return nil, fmt.Errorf("Load Balancer %q (resource group %q) not found for Outbound Rule %q", id.LoadBalancerName, id.ResourceGroup, id.OutboundRuleName)
		}
		return nil, fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Outbound Rule %q", id.LoadBalancerName, id.ResourceGroup, id.OutboundRuleName)
	}
	props := lb.LoadBalancerPropertiesFormat
	if props == nil || props.OutboundRules == nil || len(*props.OutboundRules) == 0 {
		return nil, fmt.Errorf("Outbound Rule %q not found in Load Balancer %q (resource group %q)", id.OutboundRuleName, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *props.OutboundRules {
		if v.Name != nil && *v.Name == id.OutboundRuleName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Outbound Rule %q not found in Load Balancer %q (resource group %q)", id.OutboundRuleName, id.LoadBalancerName, id.ResourceGroup)
	}
	return utils.Bool(found), nil
}

func (r LoadBalancerOutboundRule) IsMissing(loadBalancerName string, outboundRuleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[loadBalancerName]
		if !ok {
			return fmt.Errorf("not found: %q", loadBalancerName)
		}

		id, err := parse.LoadBalancerID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Outbound Rule removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Outbound Rule removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.OutboundRules == nil {
			return fmt.Errorf("Outbound Rule %q not found in Load Balancer %q (resource group %q)", outboundRuleName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.OutboundRules {
			if v.Name != nil && *v.Name == outboundRuleName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Outbound Rule %q not removed from Load Balancer %q (resource group %q)", outboundRuleName, id.Name, id.ResourceGroup)
		}
		return nil
	}
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
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
  resource_group_name     = azurerm_lb_outbound_rule.test.resource_group_name
  loadbalancer_id         = azurerm_lb_outbound_rule.test.loadbalancer_id
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = azurerm_lb_outbound_rule.test.frontend_ip_configuration[0].name
  }
}
`, template)
}

func (r LoadBalancerOutboundRule) removal(data acceptance.TestData) string {
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

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = azurerm_resource_group.test.name
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule-%d"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = azurerm_resource_group.test.name
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
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
