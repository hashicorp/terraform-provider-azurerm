package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubConsumerGroup_basic(t *testing.T) {
	resourceName := "azurerm_eventhub_consumer_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubConsumerGroup_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubConsumerGroupExists(resourceName),
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

func TestAccAzureRMEventHubConsumerGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_eventhub_consumer_group.test"
	ri := tf.AccRandTimeInt()

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubConsumerGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubConsumerGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMEventHubConsumerGroup_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_eventhub_consumer_group"),
			},
		},
	})
}

func TestAccAzureRMEventHubConsumerGroup_complete(t *testing.T) {
	resourceName := "azurerm_eventhub_consumer_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubConsumerGroup_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubConsumerGroupExists(resourceName),
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

func TestAccAzureRMEventHubConsumerGroup_userMetadataUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_consumer_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubConsumerGroup_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubConsumerGroupExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMEventHubConsumerGroup_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubConsumerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_metadata", "some-meta-data"),
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

func testCheckAzureRMEventHubConsumerGroupDestroy(s *terraform.State) error {
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

func testCheckAzureRMEventHubConsumerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.ConsumerGroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMEventHubConsumerGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMEventHubConsumerGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMEventHubConsumerGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_consumer_group" "import" {
  name                = "${azurerm_eventhub_consumer_group.test.name}"
  namespace_name      = "${azurerm_eventhub_consumer_group.test.namespace_name}"
  eventhub_name       = "${azurerm_eventhub_consumer_group.test.eventhub_name}"
  resource_group_name = "${azurerm_eventhub_consumer_group.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMEventHubConsumerGroup_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  user_metadata       = "some-meta-data"
}
`, rInt, location, rInt, rInt, rInt)
}
