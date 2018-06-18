package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceAzureRMLoadBalancerPrivateIpAddressAllocation_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "Static",
			ErrCount: 0,
		},
		{
			Value:    "Dynamic",
			ErrCount: 0,
		},
		{
			Value:    "STATIC",
			ErrCount: 0,
		},
		{
			Value:    "static",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateLoadBalancerPrivateIpAddressAllocation(tc.Value, "azurerm_lb")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM LoadBalancer private_ip_address_allocation to trigger a validation error")
		}
	}
}

func TestAccAzureRMLoadBalancer_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancer_standard(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_standard(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancer_frontEndConfig(t *testing.T) {
	var lb network.LoadBalancer
	resourceName := "azurerm_lb.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_frontEndConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.#", "2"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancer_frontEndConfigRemovalWithIP(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.#", "1"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancer_frontEndConfigRemoval(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancer_tags(t *testing.T) {
	var lb network.LoadBalancer
	resourceName := "azurerm_lb.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Purpose", "AcceptanceTests"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancer_updatedTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Purpose", "AcceptanceTests"),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerExists(name string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		loadBalancerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for loadbalancer: %s", loadBalancerName)
		}

		client := testAccProvider.Meta().(*ArmClient).loadBalancerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, loadBalancerName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: LoadBalancer %q (resource group: %q) does not exist", loadBalancerName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on loadBalancerClient: %+v", err)
		}

		*lb = resp

		return nil
	}
}

func testCheckAzureRMLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).loadBalancerClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_lb" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("LoadBalancer still exists:\n%#v", resp.LoadBalancerPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMLoadBalancer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    tags {
    	Environment = "production"
    	Purpose = "AcceptanceTests"
    }

}`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_lb" "test" {
    name = "acctest-loadbalancer-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku = "Standard"

    tags {
      Environment = "production"
      Purpose = "AcceptanceTests"
    }

}`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_updatedTags(rInt int, location string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    tags {
    	Purpose = "AcceptanceTests"
    }

}`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfig(rInt int, location string) string {
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

resource "azurerm_public_ip" "test1" {
    name = "another-test-ip-%d"
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

    frontend_ip_configuration {
      name = "two-%d"
      public_ip_address_id = "${azurerm_public_ip.test1.id}"
    }
}`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfig_withZone(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      subnet_id = "${azurerm_subnet.test.id}"
      zones = ["1"]
    }

    frontend_ip_configuration {
      name = "two-%d"
      subnet_id = "${azurerm_subnet.test.id}"
      zones = ["1"]
    }
}`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfigRemovalWithIP(rInt int, location string) string {
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

resource "azurerm_public_ip" "test1" {
    name = "another-test-ip-%d"
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
}`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfigRemoval(rInt int, location string) string {
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
}`, rInt, location, rInt, rInt, rInt)
}
