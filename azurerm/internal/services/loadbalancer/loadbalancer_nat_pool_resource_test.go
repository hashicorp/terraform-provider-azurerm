package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLoadBalancerNatPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLoadBalancerNatPool_requiresImport),
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatPool_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolIsMissing("azurerm_lb.test", fmt.Sprintf("NatPool-%d", data.RandomInteger)),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePools(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolExists(data.ResourceName),
					testCheckAzureRMLoadBalancerNatPoolExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3390"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatPoolExists(data.ResourceName),
					testCheckAzureRMLoadBalancerNatPoolExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3391"),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatPoolIsMissing(loadBalancerName string, natPoolName string) resource.TestCheckFunc {
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
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Nat Pool removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Nat Pool removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.InboundNatPools == nil {
			return fmt.Errorf("Nat Pool %q not found in Load Balancer %q (resource group %q)", natPoolName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.InboundNatPools {
			if v.Name != nil && *v.Name == natPoolName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Nat Pool %q not removed from Load Balancer %q (resource group %q)", natPoolName, id.Name, id.ResourceGroup)
		}
		return nil
	}
}

func testCheckAzureRMLoadBalancerNatPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		id, err := parse.LoadBalancerInboundNatPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found for Nat Pool %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatPoolName)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Nat Pool %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatPoolName)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.InboundNatPools == nil || len(*props.InboundNatPools) == 0 {
			return fmt.Errorf("Nat Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatPoolName, id.LoadBalancerName, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.InboundNatPools {
			if v.Name != nil && *v.Name == id.InboundNatPoolName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Nat Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatPoolName, id.LoadBalancerName, id.ResourceGroup)
		}
		return nil

	}
}

func testAccAzureRMLoadBalancerNatPool_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerNatPool_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLoadBalancerNatPool_basic(data)
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

func testAccAzureRMLoadBalancerNatPool_removal(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerNatPool_multiplePools(data, data2 acceptance.TestData) string {
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

func testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(data, data2 acceptance.TestData) string {
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
