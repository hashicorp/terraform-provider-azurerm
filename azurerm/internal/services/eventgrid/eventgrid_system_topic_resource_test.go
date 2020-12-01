package eventgrid_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventGridSystemTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridSystemTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridSystemTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridSystemTopicExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "source_arm_resource_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "topic_type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metric_arm_resource_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridSystemTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridSystemTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridSystemTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridSystemTopicExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMEventGridSystemTopic_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_system_topic"),
			},
		},
	})
}

func TestAccAzureRMEventGridSystemTopic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridSystemTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridSystemTopic_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridSystemTopicExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Foo", "Bar"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "source_arm_resource_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "topic_type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metric_arm_resource_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMEventGridSystemTopicDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.SystemTopicsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_system_topic" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Event Grid System Topic still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMEventGridSystemTopicExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.SystemTopicsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Grid System Topic: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Grid System Topic %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on EventGridSystemTopicsClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMEventGridSystemTopic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctestEGST%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12), data.RandomIntOfLength(10))
}

func testAccAzureRMEventGridSystemTopic_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMEventGridSystemTopic_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_system_topic" "import" {
  name                   = azurerm_eventgrid_system_topic.test.name
  location               = azurerm_eventgrid_system_topic.test.location
  resource_group_name    = azurerm_eventgrid_system_topic.test.resource_group_name
  source_arm_resource_id = azurerm_eventgrid_system_topic.test.source_arm_resource_id
  topic_type             = azurerm_eventgrid_system_topic.test.topic_type
}
`, template)
}

func testAccAzureRMEventGridSystemTopic_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctestEGST%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"

  tags = {
    "Foo" = "Bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12), data.RandomIntOfLength(10))
}
