package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateEndpoint_basic(t *testing.T) {
	resourceName := "azurerm_private_endpoint.test"
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

func TestAccAzureRMPrivateEndpoint_requestMessage(t *testing.T) {
	resourceName := "azurerm_private_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_requestMessage(ri, location, "CATS: ALL YOUR BASE ARE BELONG TO US."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "private_service_connection.0.request_message", "CATS: ALL YOUR BASE ARE BELONG TO US."),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateEndpoint_requestMessage(ri, location, "CAPTAIN: WHAT YOU SAY!!"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "private_service_connection.0.request_message", "CAPTAIN: WHAT YOU SAY!!"),
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

// The update and complete test cases had to be totally removed since there is a bug with tags and the support for
// tags has been removed, all other attributes are ForceNew.
// API Issue "Unable to remove Tags from Private Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467

func testCheckAzureRMPrivateEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Endpoint not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Network.PrivateEndpointClient
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
	client := testAccProvider.Meta().(*ArmClient).Network.PrivateEndpointClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_endpoint" {
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

func testAccAzureRMPrivateEndpointTemplate_template(rInt int, location string, seviceCfg string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                   = "acctestsnetservice-%d"
  resource_group_name    = azurerm_resource_group.test.name
  virtual_network_name   = azurerm_virtual_network.test.name
  address_prefix         = "10.5.1.0/24"

  enforce_private_link_service_network_policies  = true
}

resource "azurerm_subnet" "endpoint" {
  name                   = "acctestsnetendpoint-%d"
  resource_group_name    = azurerm_resource_group.test.name
  virtual_network_name   = azurerm_virtual_network.test.name
  address_prefix         = "10.5.2.0/24"

  enforce_private_link_endpoint_network_policies = true
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

%s
`, rInt, location, rInt, rInt, rInt, rInt, rInt, seviceCfg)
}

func testAccAzureRMPrivateEndpoint_serviceAutoApprove(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    primary                      = true
    subnet_id                    = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
  azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, rInt, rInt)
}

func testAccAzureRMPrivateEndpoint_serviceManualApprove(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    primary                      = true
    subnet_id                    = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
  azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, rInt, rInt)
}

func testAccAzureRMPrivateEndpoint_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }
}
`, testAccAzureRMPrivateEndpointTemplate_template(rInt, location, testAccAzureRMPrivateEndpoint_serviceAutoApprove(rInt)), rInt)
}

func testAccAzureRMPrivateEndpoint_requestMessage(rInt int, location string, msg string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = true
    private_connection_resource_id = azurerm_private_link_service.test.id
    request_message                = %q
  }
}
`, testAccAzureRMPrivateEndpointTemplate_template(rInt, location, testAccAzureRMPrivateEndpoint_serviceManualApprove(rInt)), rInt, msg)
}
