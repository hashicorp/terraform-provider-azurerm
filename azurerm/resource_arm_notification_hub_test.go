package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMNotificationHub_basic(t *testing.T) {
	resourceName := "azurerm_notification_hub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gcm_credential.#", "0"),
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

func TestAccAzureRMNotificationHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_notification_hub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gcm_credential.#", "0"),
				),
			},
			{
				Config:      testAccAzureRMNotificationHub_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_notification_hub"),
			},
		},
	})
}

func testCheckAzureRMNotificationHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).notificationHubs.HubsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		hubName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, hubName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationHubsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub does not exist: %s", hubName)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).notificationHubs.HubsClient
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

func testAccAzureRMNotificationHub_basic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
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

func testAccAzureRMNotificationHub_requiresImport(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "import" {
  name                = "${azurerm_notification_hub.test.name}"
  namespace_name      = "${azurerm_notification_hub.test.namespace_name}"
  resource_group_name = "${azurerm_notification_hub.test.resource_group_name}"
  location            = "${azurerm_notification_hub.test.location}"
}
`, testAccAzureRMNotificationHub_basic(ri, location))
}
