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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHubNamespace_free(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(resourceName),
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

func TestAccAzureRMNotificationHubNamespace_withTags(t *testing.T) {
	resourceName := "azurerm_notification_hub_namespace.test"

	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHubNamespace_withTags(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(resourceName),
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
				Config: testAccAzureRMNotificationHubNamespace_withTagsUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(resourceName),
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

func testAccAzureRMNotificationHubNamespace_free(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

func testAccAzureRMNotificationHubNamespace_withTags(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  tags {
    environment = "test"
    company     = "hashicorp"
  }
}
`, ri, location, ri)
}

func testAccAzureRMNotificationHubNamespace_withTagsUpdated(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  tags {
    environment = "production"
  }
}
`, ri, location, ri)
}
