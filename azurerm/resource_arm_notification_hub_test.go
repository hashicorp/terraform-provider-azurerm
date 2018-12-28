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

func TestAccAzureRMNotificationHub_withTags(t *testing.T) {
	resourceName := "azurerm_notification_hub.test"

	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_withTags(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gcm_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.company", "hashicorp"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNotificationHub_withTagsUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gcm_credential.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
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

func testAccAzureRMNotificationHub_preReqs(ri int, location string) string {
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
`, ri, location, ri)
}

func testAccAzureRMNotificationHub_basic(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, testAccAzureRMNotificationHub_preReqs(ri, location), ri)
}

func testAccAzureRMNotificationHub_withTags(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tags {
    environment = "test"
    company     = "hashicorp"
  }
}
`, testAccAzureRMNotificationHub_preReqs(ri, location), ri)
}

func testAccAzureRMNotificationHub_withTagsUpdated(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tags {
    environment = "production"
  }
}
`, testAccAzureRMNotificationHub_preReqs(ri, location), ri)
}
