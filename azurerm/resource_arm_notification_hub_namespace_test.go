package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNotificationHubNamespace_free(t *testing.T) {
	resourceName := "azurerm_notification_hub_namespace.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubNamespace_free(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMNotificationHubNamespaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).notificationNamespacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationNamespacesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub Namespace does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubNamespaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).notificationNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_notification_hub_namespace" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, resourceGroup, namespaceName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Notification Hub Namespace still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMNotificationHubNamespace_free(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku {
    name = "Free"
  }
}
`, ri, location, ri)
}

func testAzureRMNotificationHubNamespace_basic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku {
    name = "Basic"
  }
}
`, ri, location, ri)
}
