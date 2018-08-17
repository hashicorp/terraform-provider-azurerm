package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNotificationHub_basic(t *testing.T) {
	resourceName := "azurerm_notification_hub.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gcm_credential.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMNotificationHub_requiresImport(t *testing.T) {
	resourceName := "azurerm_notification_hub.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
				),
			},
			{
				Config:      testAzureRMNotificationHub_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_notification_hub"),
			},
		},
	})
}

func testCheckAzureRMNotificationHubExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).notificationHubsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		hubName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, hubName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationHubsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).notificationHubsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_notification_hub" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		hubName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, hubName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Notification Hub still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMNotificationHub_basic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrgpol-%d"
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

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}" 
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, ri, location, ri, ri)
}

func testAzureRMNotificationHub_requiresImport(rInt int, location string) string {
	template := testAzureRMNotificationHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "import" {
  name                = "${azurerm_notification_hub.test.name}"
  namespace_name      = "${azurerm_notification_hub.test.namespace_name}" 
  resource_group_name = "${azurerm_notification_hub.test.resource_group_name}"
  location            = "${azurerm_notification_hub.test.location}"
}
`, template)
}
