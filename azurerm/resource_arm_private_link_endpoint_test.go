package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMPrivateLinkEndpoint_basic(t *testing.T) {
	resourceName := "azurerm_private_link_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkEndpoint_basic(ri, location),
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

func TestAccAzureRMPrivateLinkEndpoint_requestMessage(t *testing.T) {
	resourceName := "azurerm_private_link_endpoint.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkEndpoint_requestMessage(ri, location, "CATS: ALL YOUR BASE ARE BELONG TO US."),
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
				Config: testAccAzureRMPrivateLinkEndpoint_requestMessage(ri, location, "CAPTAIN: WHAT YOU SAY!!"),
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
// API Issue "Unable to remove Tags from Private Link Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467

func testAccAzureRMPrivateLinkEndpoint_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_endpoint" "test" {
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

func testAccAzureRMPrivateLinkEndpoint_requestMessage(rInt int, location string, msg string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_endpoint" "test" {
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
