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

func TestAccEventGridDomainTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventGridDomainTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventGridDomainTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventGridDomainTopicExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccEventGridDomainTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventGridDomainTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventGridDomainTopicExists(data.ResourceName),
				),
			},
			{
				Config:      testAccEventGridDomainTopic_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_domain_topic"),
			},
		},
	})
}

func testCheckEventGridDomainTopicDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainTopicsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_domain_topic" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		domainName := rs.Primary.Attributes["domain_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, domainName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("EventGrid Domain Topic still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckEventGridDomainTopicExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainTopicsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		domainName := rs.Primary.Attributes["domain_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for EventGrid Domain Topic: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, domainName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventGrid Domain Topic %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on EventGrid.DomainTopicsClient: %s", err)
		}

		return nil
	}
}

func testAccEventGridDomainTopic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_eventgrid_domain" "test" {
  name                = "acctestegdomain-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_eventgrid_domain_topic" "test" {
  name                = "acctestegtopic-%d"
  domain_name         = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccEventGridDomainTopic_requiresImport(data acceptance.TestData) string {
	template := testAccEventGridDomainTopic_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_eventgrid_domain_topic" "import" {
  name                = azurerm_eventgrid_domain_topic.test.name
  domain_name         = azurerm_eventgrid_domain_topic.test.domain_name
  resource_group_name = azurerm_eventgrid_domain_topic.test.resource_group_name
}
`, template)
}
