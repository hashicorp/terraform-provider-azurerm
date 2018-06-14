package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLoadBalancerNatPool_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatPools/%s",
		subscriptionID, ri, ri, natPoolName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_pool.test", "id", natPoolId),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatPool_removal(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolNotExists(natPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_update(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)
	natPool2Name := fmt.Sprintf("NatPool-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePools(ri, natPoolName, natPool2Name, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPool2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_nat_pool.test2", "backend_port", "3390"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(ri, natPoolName, natPool2Name, testLocation()),
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

func TestAccAzureRMLoadBalancerNatPool_reapply(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	deleteNatPoolState := func(s *terraform.State) error {
		return s.Remove("azurerm_lb_nat_pool.test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					deleteNatPoolState,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, testLocation()),
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
		_, _, exists := findLoadBalancerNatPoolByName(lb, natPoolName)
		if !exists {
			return fmt.Errorf("A NAT Pool with name %q cannot be found.", natPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatPoolNotExists(natPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerNatPoolByName(lb, natPoolName)
		if exists {
			return fmt.Errorf("A NAT Pool with name %q has been found.", natPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatPoolDisappears(natPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).loadBalancerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, i, exists := findLoadBalancerNatPoolByName(lb, natPoolName)
		if !exists {
			return fmt.Errorf("A Nat Pool with name %q cannot be found.", natPoolName)
		}

		currentPools := *lb.LoadBalancerPropertiesFormat.InboundNatPools
		pools := append(currentPools[:i], currentPools[i+1:]...)
		lb.LoadBalancerPropertiesFormat.InboundNatPools = &pools

		id, err := parseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating LoadBalancer %+v", err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for the completion of LoadBalancer %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerNatPool_basic(rInt int, natPoolName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "test-ip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_pool" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natPoolName, rInt)
}

func testAccAzureRMLoadBalancerNatPool_removal(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "test-ip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancerNatPool_multiplePools(rInt int, natPoolName, natPool2Name string, location string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "test-ip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}

resource "azurerm_lb_nat_pool" "test" {
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id = "${azurerm_lb.test.id}"
  name = "%s"
  protocol = "Tcp"
  frontend_port_start = 80
  frontend_port_end = 81
  backend_port = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id = "${azurerm_lb.test.id}"
  name = "%s"
  protocol = "Tcp"
  frontend_port_start = 82
  frontend_port_end = 83
  backend_port = 3390
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natPoolName, rInt, natPool2Name, rInt)
}

func testAccAzureRMLoadBalancerNatPool_multiplePoolsUpdate(rInt int, natPoolName, natPool2Name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "test-ip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_pool" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port_start            = 82
  frontend_port_end              = 83
  backend_port                   = 3391
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natPoolName, rInt, natPool2Name, rInt)
}
