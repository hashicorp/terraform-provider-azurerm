package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
			t.Fatalf("Expected the Azure RM Load Balancer private_ip_address_allocation to trigger a validation error")
		}
	}
}

func TestAccAzureRMLoadBalancer_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLoadBalancer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancer_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_lb"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancer_standard(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAzureRMLoadBalancer_frontEndConfigPublicIPPrefix(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	location := testLocation()
	resourceName := "azurerm_lb.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_frontEndConfigPublicIPPrefix(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.#", "1"),
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
func TestAccAzureRMLoadBalancer_frontEndConfig(t *testing.T) {
	var lb network.LoadBalancer
	resourceName := "azurerm_lb.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMLoadBalancer_emptyPrivateIP(t *testing.T) {
	resourceName := "azurerm_lb.test"
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancer_emptyIPAddress(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttrSet(resourceName, "frontend_ip_configuration.0.private_ip_address"),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerExists(resourceName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		loadBalancerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for loadbalancer: %s", loadBalancerName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.LoadBalancersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, loadBalancerName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: Load Balancer %q (resource group: %q) does not exist", loadBalancerName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on loadBalancerClient: %+v", err)
		}

		*lb = resp

		return nil
	}
}

func testCheckAzureRMLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.LoadBalancersClient
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLoadBalancer_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb" "import" {
  name                = "${azurerm_lb.test.name}"
  location            = "${azurerm_lb.test.location}"
  resource_group_name = "${azurerm_lb.test.resource_group_name}"

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, template)
}

func testAccAzureRMLoadBalancer_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_updatedTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    Purpose = "AcceptanceTests"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfig(rInt int, location string) string {
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

resource "azurerm_public_ip" "test1" {
  name                = "another-test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    name                 = "two-%d"
    public_ip_address_id = "${azurerm_public_ip.test1.id}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfigRemovalWithIP(rInt int, location string) string {
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

resource "azurerm_public_ip" "test1" {
  name                = "another-test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfigPublicIPPrefix(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "test-ip-prefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  prefix_length       = 31
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                = "prefix-%d"
    public_ip_prefix_id = "${azurerm_public_ip_prefix.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_frontEndConfigRemoval(rInt int, location string) string {
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
  name                = "acctest-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancer_emptyIPAddress(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Basic"

  frontend_ip_configuration {
    name                          = "Internal"
    private_ip_address_allocation = "Dynamic"
    private_ip_address            = ""
    subnet_id                     = "${azurerm_subnet.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}
