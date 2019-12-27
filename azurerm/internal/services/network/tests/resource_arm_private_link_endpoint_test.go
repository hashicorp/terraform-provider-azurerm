package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMPrivateLinkEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateLinkEndpoint_requestMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkEndpoint_requestMessage(data, "CATS: ALL YOUR BASE ARE BELONG TO US."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_service_connection.0.request_message", "CATS: ALL YOUR BASE ARE BELONG TO US."),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkEndpoint_requestMessage(data, "CAPTAIN: WHAT YOU SAY!!"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_service_connection.0.request_message", "CAPTAIN: WHAT YOU SAY!!"),
				),
			},
			data.ImportStep(),
		},
	})
}

// The update and complete test cases had to be totally removed since there is a bug with tags and the support for
// tags has been removed, all other attributes are ForceNew.
// API Issue "Unable to remove Tags from Private Link Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467

func testAccAzureRMPrivateLinkEndpoint_basic(data acceptance.TestData) string {
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
`, testAccAzureRMPrivateEndpointTemplate_template(data, testAccAzureRMPrivateEndpoint_serviceAutoApprove(data)), data.RandomInteger)
}

func testAccAzureRMPrivateLinkEndpoint_requestMessage(data acceptance.TestData, msg string) string {
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
`, testAccAzureRMPrivateEndpointTemplate_template(data, testAccAzureRMPrivateEndpoint_serviceManualApprove(data)), data.RandomInteger, msg)
}
