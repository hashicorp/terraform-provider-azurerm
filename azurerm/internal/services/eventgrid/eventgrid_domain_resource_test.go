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

func TestAccAzureRMEventGridDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMEventGridDomain_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_domain"),
			},
		},
	})
}

func TestAccAzureRMEventGridDomain_mapping(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_mapping(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "input_mapping_fields.0.topic", "test"),
					resource.TestCheckResourceAttr(data.ResourceName, "input_mapping_fields.0.topic", "test"),
					resource.TestCheckResourceAttr(data.ResourceName, "input_mapping_default_values.0.data_version", "1.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "input_mapping_default_values.0.subject", "DefaultSubject"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventGridDomain_basicWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_basicWithTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.foo", "bar"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMEventGridDomainDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_domain" {
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
			return fmt.Errorf("EventGrid Domain still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMEventGridDomainExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for EventGrid Domain: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventGrid Domain %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventGridDomainsClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMEventGridDomain_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMEventGridDomain_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMEventGridDomain_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_domain" "import" {
  name                = azurerm_eventgrid_domain.test.name
  location            = azurerm_eventgrid_domain.test.location
  resource_group_name = azurerm_eventgrid_domain.test.resource_group_name
}
`, template)
}

func testAccAzureRMEventGridDomain_mapping(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  input_schema = "CustomEventSchema"

  input_mapping_fields {
    topic      = "test"
    event_type = "test"
  }

  input_mapping_default_values {
    data_version = "1.0"
    subject      = "DefaultSubject"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMEventGridDomain_basicWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
