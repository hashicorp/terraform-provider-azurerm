package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventGridTopic_basic(t *testing.T) {
	resourceName := "azurerm_eventgrid_topic.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridTopic_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridTopicExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func TestAccAzureRMEventGridTopic_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_eventgrid_topic.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridTopic_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridTopicExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMEventGridTopic_requiresImport(ri),
				ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_topic"),
			},
		},
	})
}

func TestAccAzureRMEventGridTopic_basicWithTags(t *testing.T) {
	resourceName := "azurerm_eventgrid_topic.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridTopic_basicWithTags(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridTopicExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func testCheckAzureRMEventGridTopicDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.TopicsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_topic" {
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
			return fmt.Errorf("EventGrid Topic still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMEventGridTopicExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for EventGrid Topic: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.TopicsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventGrid Topic %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventGridTopicsClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMEventGridTopic_basic(rInt int) string {
	// TODO: confirm if this is still the case
	// currently only supported in "West Central US" & "West US 2"
	location := "westus2"
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctesteg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMEventGridTopic_requiresImport(rInt int) string {
	template := testAccAzureRMEventGridTopic_basic(rInt)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_topic" "import" {
  name                = "${azurerm_eventgrid_topic.test.name}"
  location            = "${azurerm_eventgrid_topic.test.location}"
  resource_group_name = "${azurerm_eventgrid_topic.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMEventGridTopic_basicWithTags(rInt int) string {
	// currently only supported in "West Central US" & "West US 2"
	location := "westus2"
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctesteg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    "foo" = "bar"
  }
}
`, rInt, location, rInt)
}
