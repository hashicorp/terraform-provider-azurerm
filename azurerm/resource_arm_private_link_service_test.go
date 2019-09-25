package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateLinkService_basic(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "testFqdns"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configurations.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_complete(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "testFqdns2"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_update(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "testFqdns"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configurations.#", "1"),
				),
			},
			{
				Config: testAccAzureRMPrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "testFqdns2"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fqdns.0", "testFqdns"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configurations.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMPrivateLinkServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Link Service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Private Link Service %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateLinkServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_link_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateLinkService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-%d"
  resource_group_name                   = "${azurerm_resource_group.test.name}"
  virtual_network_name                  = "${azurerm_virtual_network.test.name}"
  address_prefix                        = "10.5.1.0/24"
  private_link_service_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "${azurerm_public_ip.test.name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestpls-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  fqdns               = ["testFqdns"]

  ip_configurations {
    name                         = "${azurerm_public_ip.test.name}"
    subnet_id                    = "${azurerm_subnet.test.id}"
    private_ip_address           = "10.5.1.17"
    private_ip_address_version   = "IPv4"
    private_ip_allocation_method = "Static"
  }

  load_balancer_frontend_ip_configurations {
    id = "${azurerm_lb.test.frontend_ip_configuration.0.id}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-%d"
  resource_group_name                   = "${azurerm_resource_group.test.name}"
  virtual_network_name                  = "${azurerm_virtual_network.test.name}"
  address_prefix                        = "10.5.1.0/24"
  private_link_service_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "${azurerm_public_ip.test.name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestpls-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  fqdns               = ["testFqdns2"]

  ip_configurations {
    name                         = "${azurerm_public_ip.test.name}"
    subnet_id                    = "${azurerm_subnet.test.id}"
    private_ip_address           = "10.5.1.17"
    private_ip_address_version   = "IPv4"
    private_ip_allocation_method = "Static"
  }

  load_balancer_frontend_ip_configurations {
    id = "${azurerm_lb.test.frontend_ip_configuration.0.id}"
  }

  tags = {
    env = "test"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}
