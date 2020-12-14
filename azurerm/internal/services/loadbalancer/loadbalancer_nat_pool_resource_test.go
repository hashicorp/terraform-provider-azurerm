package loadbalancer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer"
)

func TestAccAzureRMLoadBalancerNatPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	var lb network.LoadBalancer
	natPoolName := fmt.Sprintf("NatPool-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatPools/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natPoolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data, natPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_pool.test", "id", natPoolId),
				),
			},
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	var lb network.LoadBalancer
	natPoolName := fmt.Sprintf("NatPool-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatPools/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natPoolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data, natPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_pool.test", "id", natPoolId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerNatPool_requiresImport(data, natPoolName),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_nat_pool"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	var lb network.LoadBalancer
	natPoolName := fmt.Sprintf("NatPool-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data, natPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatPool_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolNotExists(natPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test2")

	var lb network.LoadBalancer
	natPoolName := fmt.Sprintf("NatPool-%d", data.RandomInteger)
	natPool2Name := fmt.Sprintf("NatPool-%d", data2.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePools(data, natPoolName, natPool2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPool2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_nat_pool.test2", "backend_port", "3390"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(data, natPoolName, natPool2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPool2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_nat_pool.test2", "backend_port", "3391"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_pool", "test")

	var lb network.LoadBalancer
	natPoolName := fmt.Sprintf("NatPool-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(data, natPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					testCheckAzureRMLoadBalancerNatPoolDisappears(natPoolName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatPoolExists(natPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerNatPoolByName(lb, natPoolName)
		if !exists {
			return fmt.Errorf("A NAT Pool with name %q cannot be found.", natPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatPoolNotExists(natPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerNatPoolByName(lb, natPoolName)
		if exists {
			return fmt.Errorf("A NAT Pool with name %q has been found.", natPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatPoolDisappears(natPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := loadbalancer.FindLoadBalancerNatPoolByName(lb, natPoolName)
		if !exists {
			return fmt.Errorf("A Nat Pool with name %q cannot be found.", natPoolName)
		}

		currentPools := *lb.LoadBalancerPropertiesFormat.InboundNatPools
		pools := append(currentPools[:i], currentPools[i+1:]...)
		lb.LoadBalancerPropertiesFormat.InboundNatPools = &pools

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the completion of Load Balancer %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerNatPool_basic(data acceptance.TestData, natPoolName string) string {
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
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, natPoolName, data.RandomInteger)
}

func testAccAzureRMLoadBalancerNatPool_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerNatPool_basic(data, name)
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

func testAccAzureRMLoadBalancerNatPool_multiplePools(data acceptance.TestData, natPoolName, natPool2Name string) string {
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
  name                = "%s"
  protocol            = "Tcp"
  frontend_port_start = 80
  frontend_port_end   = 81
  backend_port        = 3389

  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "%s"
  protocol            = "Tcp"
  frontend_port_start = 82
  frontend_port_end   = 83
  backend_port        = 3390

  frontend_ip_configuration_name = "one-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, natPoolName, data.RandomInteger, natPool2Name, data.RandomInteger)
}

func testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(data acceptance.TestData, natPoolName, natPool2Name string) string {
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
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 82
  frontend_port_end              = 83
  backend_port                   = 3391
  frontend_ip_configuration_name = "one-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, natPoolName, data.RandomInteger, natPool2Name, data.RandomInteger)
}
