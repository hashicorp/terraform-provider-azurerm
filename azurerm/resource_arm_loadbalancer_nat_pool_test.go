package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLoadBalancerNatPool_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatPools/%s",
		subscriptionID, ri, ri, natPoolName)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
				// location is deprecated and was never actually used
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)
	location := testLocation()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natPoolId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatPools/%s",
		subscriptionID, ri, ri, natPoolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatPool_basic(ri, natPoolName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatPoolExists(natPoolName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_pool.test", "id", natPoolId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerNatPool_requiresImport(ri, natPoolName, location),
				ExpectError: testRequiresImportError("azurerm_lb_nat_pool"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatPool_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
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
	ri := tf.AccRandTimeInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)
	natPool2Name := fmt.Sprintf("NatPool-%d", tf.AccRandTimeInt())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMLoadBalancerNatPool_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natPoolName := fmt.Sprintf("NatPool-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
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
		client := testAccProvider.Meta().(*ArmClient).Network.LoadBalancersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, i, exists := findLoadBalancerNatPoolByName(lb, natPoolName)
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

func testAccAzureRMLoadBalancerNatPool_basic(rInt int, natPoolName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
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

func testAccAzureRMLoadBalancerNatPool_requiresImport(rInt int, name string, location string) string {
	template := testAccAzureRMLoadBalancerNatPool_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_pool" "import" {
  name                           = "${azurerm_lb_nat_pool.test.name}"
  loadbalancer_id                = "${azurerm_lb_nat_pool.test.loadbalancer_id}"
  location                       = "${azurerm_lb_nat_pool.test.location}"
  resource_group_name            = "${azurerm_lb_nat_pool.test.resource_group_name}"
  frontend_ip_configuration_name = "${azurerm_lb_nat_pool.test.frontend_ip_configuration_name}"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 3389
}
`, template)
}

func testAccAzureRMLoadBalancerNatPool_removal(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  allocation_method = "Static"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  protocol            = "Tcp"
  frontend_port_start = 80
  frontend_port_end   = 81
  backend_port        = 3389

  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_pool" "test2" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  protocol            = "Tcp"
  frontend_port_start = 82
  frontend_port_end   = 83
  backend_port        = 3390

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
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
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
