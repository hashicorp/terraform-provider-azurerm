package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateEndpoint_basic(t *testing.T) {
	resourceName := "azurerm_private_link_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.name", fmt.Sprintf("acctestconnection-%d", ri)),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_service_connection.0.private_link_service_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.request_message", "Please approve my connection"),
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

func TestAccAzureRMPrivateEndpoint_complete(t *testing.T) {
	resourceName := "azurerm_private_link_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.name", fmt.Sprintf("acctestconnection-%d", ri)),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_service_connection.0.private_link_service_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.request_message", "plz approve my request"),
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

func TestAccAzureRMPrivateEndpoint_update(t *testing.T) {
	resourceName := "azurerm_private_link_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.name", fmt.Sprintf("acctestconnection-%d", ri)),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_service_connection.0.private_link_service_id"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.request_message", "Please approve my connection"),
				),
			},
			{
				Config: testAccAzureRMPrivateEndpoint_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.request_message", "plz approve my request"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccAzureRMPrivateEndpoint_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "private_link_service_connection.0.request_message", "Please approve my connection"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func testCheckAzureRMPrivateEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Endpoint not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.PrivateEndpointClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Private Endpoint %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on PrivateEndpointClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.PrivateEndpointClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_link_endpoint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on PrivateEndpointClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateEndpointTemplate_standardResources(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                   = "acctestsnet-%d"
  resource_group_name    = azurerm_resource_group.test.name
  virtual_network_name   = azurerm_virtual_network.test.name
  address_prefix         = "10.5.1.0/24"

  disable_private_link_service_network_policies = true
  disable_private_endpoint_network_policies     = true
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestpls-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
  }
  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateEndpoint_basic(rInt int, location string) string {
	standardResources := testAccAzureRMPrivateEndpointTemplate_standardResources(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_private_link_endpoint" "test" {
  name                = "acctestendpoint-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_link_service_connection {
    name                    = "acctestconnection-%d"
    private_link_service_id = azurerm_private_link_service.test.id
  }
}
`, standardResources, rInt, rInt)
}

func testAccAzureRMPrivateEndpoint_complete(rInt int, location string) string {
	standardResources := testAccAzureRMPrivateEndpointTemplate_standardResources(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_private_link_endpoint" "test" {
  name                = "acctestendpoint-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_link_service_connection {
    name                    = "acctestconnection-%d"
    private_link_service_id = azurerm_private_link_service.test.id
    group_ids               = []
    request_message         = "plz approve my request"
  }

  tags = {
    env = "test"
  }
}
`, standardResources, rInt, rInt)
}
