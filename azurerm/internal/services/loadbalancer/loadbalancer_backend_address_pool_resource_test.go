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

func TestAccAzureRMLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:      testAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(data),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerBackEndAddressPoolIsMissing("azurerm_lb.test", fmt.Sprintf("%d-address-pool", data.RandomInteger)),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		id, err := parse.LoadBalancerBackendAddressPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found for Backend Address Pool %q", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Backend Address Pool %q", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.BackendAddressPools == nil || len(*props.BackendAddressPools) == 0 {
			return fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.BackendAddressPools {
			if v.Name != nil && *v.Name == id.BackendAddressPoolName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroup)
		}
		return nil
	}
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolIsMissing(loadBalancerName string, backendPoolName string) resource.TestCheckFunc {
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
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Backend Address Pool removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Backend Address Pool removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.BackendAddressPools == nil {
			return fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", backendPoolName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.BackendAddressPools {
			if v.Name != nil && *v.Name == backendPoolName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Backend Pool %q not removed from Load Balancer %q (resource group %q)", backendPoolName, id.Name, id.ResourceGroup)
		}
		return nil
	}
}

func testAccAzureRMLoadBalancerBackEndAddressPool_basic(data acceptance.TestData) string {
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "%d-address-pool"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLoadBalancerBackEndAddressPool_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name                = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id     = azurerm_lb_backend_address_pool.test.loadbalancer_id
  resource_group_name = azurerm_lb_backend_address_pool.test.resource_group_name
}
`, template)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_removal(data acceptance.TestData) string {
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
