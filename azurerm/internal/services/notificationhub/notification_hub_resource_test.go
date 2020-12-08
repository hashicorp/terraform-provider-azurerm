package notificationhub_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMNotificationHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "gcm_credential.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNotificationHub_updateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNotificationHub_withoutTag(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNotificationHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNotificationHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "apns_credential.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "gcm_credential.#", "0"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMNotificationHub_requiresImport),
		},
	})
}

func testCheckAzureRMNotificationHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.HubsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.HubsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMNotificationHub_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNotificationHub_withoutTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGpol-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNotificationHub_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNotificationHub_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub" "import" {
  name                = azurerm_notification_hub.test.name
  namespace_name      = azurerm_notification_hub.test.namespace_name
  resource_group_name = azurerm_notification_hub.test.resource_group_name
  location            = azurerm_notification_hub.test.location
}
`, template)
}
