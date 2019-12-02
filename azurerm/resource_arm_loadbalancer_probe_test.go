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

func TestAccAzureRMLoadBalancerProbe_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	probeId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/probes/%s",
		subscriptionID, ri, ri, probeName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(ri, probeName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_probe.test", "id", probeId),
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

func TestAccAzureRMLoadBalancerProbe_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)
	location := testLocation()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	probeId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/probes/%s",
		subscriptionID, ri, ri, probeName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(ri, probeName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_probe.test", "id", probeId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerProbe_requiresImport(ri, probeName, location),
				ExpectError: testRequiresImportError("azurerm_lb_probe"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(ri, probeName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_removal(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeNotExists(probeName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_update(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)
	probe2Name := fmt.Sprintf("probe-%d", tf.AccRandTimeInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbes(ri, probeName, probe2Name, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeExists(probe2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test2", "port", "80"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(ri, probeName, probe2Name, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeExists(probe2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test2", "port", "8080"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_updateProtocol(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolBefore(ri, probeName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Http"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolAfter(ri, probeName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Tcp"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	probeName := fmt.Sprintf("probe-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(ri, probeName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeDisappears(probeName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerProbeExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerProbeByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A Probe with name %q cannot be found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerProbeNotExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerProbeByName(lb, natRuleName)
		if exists {
			return fmt.Errorf("A Probe with name %q has been found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerProbeDisappears(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).Network.LoadBalancersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, i, exists := findLoadBalancerProbeByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A Probe with name %q cannot be found.", addressPoolName)
		}

		currentProbes := *lb.LoadBalancerPropertiesFormat.Probes
		probes := append(currentProbes[:i], currentProbes[i+1:]...)
		lb.LoadBalancerPropertiesFormat.Probes = &probes

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion for LoadBalancer: %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerProbe_basic(rInt int, probeName string, location string) string {
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

resource "azurerm_lb_probe" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  port                = 22
}
`, rInt, location, rInt, rInt, rInt, probeName)
}

func testAccAzureRMLoadBalancerProbe_requiresImport(rInt int, name string, location string) string {
	template := testAccAzureRMLoadBalancerProbe_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_probe" "import" {
  name                = "${azurerm_lb_probe.test.name}"
  loadbalancer_id     = "${azurerm_lb_probe.test.loadbalancer_id}"
  location            = "${azurerm_lb_probe.test.location}"
  resource_group_name = "${azurerm_lb_probe.test.resource_group_name}"
  port                = 22
}
`, template)
}

func testAccAzureRMLoadBalancerProbe_removal(rInt int, location string) string {
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

func testAccAzureRMLoadBalancerProbe_multipleProbes(rInt int, probeName, probe2Name string, location string) string {
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

resource "azurerm_lb_probe" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  port                = 80
}
`, rInt, location, rInt, rInt, rInt, probeName, probe2Name)
}

func testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(rInt int, probeName, probe2Name string, location string) string {
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

resource "azurerm_lb_probe" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  port                = 8080
}
`, rInt, location, rInt, rInt, rInt, probeName, probe2Name)
}

func testAccAzureRMLoadBalancerProbe_updateProtocolBefore(rInt int, probeName string, location string) string {
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

resource "azurerm_lb_probe" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  protocol            = "Http"
  request_path        = "/"
  port                = 80
}
`, rInt, location, rInt, rInt, rInt, probeName)
}

func testAccAzureRMLoadBalancerProbe_updateProtocolAfter(rInt int, probeName string, location string) string {
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

resource "azurerm_lb_probe" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "%s"
  protocol            = "Tcp"
  port                = 80
}
`, rInt, location, rInt, rInt, rInt, probeName)
}
