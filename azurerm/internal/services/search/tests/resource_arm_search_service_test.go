package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSearchService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSearchService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_search_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSearchService_requiresImport),
		},
	})
}

func TestAccAzureRMSearchService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "replica_count", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "query_keys.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSearchService_tagUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_withCustomTagValue(data, "staging"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			{
				Config: testAccAzureRMSearchService_withCustomTagValue(data, "production"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func testCheckAzureRMSearchServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Search.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.SearchServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Search Service %q (resource group %q) was not found: %+v", id.Name, id.ResourceGroup, err)
			}

			return fmt.Errorf("Bad: GetSearchService: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSearchServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Search.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_search_service" {
			continue
		}

		id, err := parse.SearchServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Bad: Search Service %q (resource group %q) still exists: %+v", id.Name, id.ResourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSearchService_withCustomTagValue(data acceptance.TestData, tagValue string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"

  tags = {
    environment = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, tagValue)
}

func testAccAzureRMSearchService_basic(data acceptance.TestData) string {
	return testAccAzureRMSearchService_withCustomTagValue(data, "staging")
}

func testAccAzureRMSearchService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSearchService_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_search_service" "import" {
  name                = "${azurerm_search_service.test.name}"
  resource_group_name = "${azurerm_search_service.test.resource_group_name}"
  location            = "${azurerm_search_service.test.location}"
  sku                 = "${azurerm_search_service.test.sku}"

  tags = {
    environment = "staging"
  }
}
`, template)
}

func testAccAzureRMSearchService_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"
  replica_count       = 2

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
