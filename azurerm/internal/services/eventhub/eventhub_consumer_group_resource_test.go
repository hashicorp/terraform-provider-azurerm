package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccEventHubConsumerGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccEventHubConsumerGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccEventHubConsumerGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_eventhub_consumer_group"),
			},
		},
	})
}

func TestAccEventHubConsumerGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccEventHubConsumerGroup_userMetadataUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
				),
			},
			{
				Config: testAccEventHubConsumerGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "user_metadata", "some-meta-data"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckEventHubConsumerGroupDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ConsumerGroupClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_consumer_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		eventHubName := rs.Primary.Attributes["eventhub_name"]

		resp, err := conn.Get(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckEventHubConsumerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ConsumerGroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub Consumer Group: %s", name)
		}

		namespaceName := rs.Primary.Attributes["namespace_name"]
		eventHubName := rs.Primary.Attributes["eventhub_name"]

		resp, err := conn.Get(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Hub Consumer Group %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventHubConsumerGroupClient: %+v", err)
		}

		return nil
	}
}

func testAccEventHubConsumerGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccEventHubConsumerGroup_requiresImport(data acceptance.TestData) string {
	template := testAccEventHubConsumerGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_consumer_group" "import" {
  name                = azurerm_eventhub_consumer_group.test.name
  namespace_name      = azurerm_eventhub_consumer_group.test.namespace_name
  eventhub_name       = azurerm_eventhub_consumer_group.test.eventhub_name
  resource_group_name = azurerm_eventhub_consumer_group.test.resource_group_name
}
`, template)
}

func testAccEventHubConsumerGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
  user_metadata       = "some-meta-data"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
