package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func TestAccAzureRMLogAnalyticsSavedSearch_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_saved_search", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSavedSearch_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSavedSearchExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsSavedSearch_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_saved_search", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSavedSearch_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSavedSearchExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsSavedSearch_withTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_saved_search", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSavedSearch_withTag(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSavedSearchExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsSavedSearch_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_saved_search", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSavedSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSavedSearch_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSavedSearchExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogAnalyticsSavedSearch_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_saved_search"),
			},
		},
	})
}

func testCheckAzureRMLogAnalyticsSavedSearchDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_saved_search" {
			continue
		}

		id, err := parse.LogAnalyticsSavedSearchID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics Saved Search still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsSavedSearchExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.SavedSearchesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsSavedSearchID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Saved Search Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("bad: Log Analytics Saved Search %q (Workspace: %q / Resource Group: %q) does not exist", id.Name, id.WorkspaceName, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsSavedSearch_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_saved_search" "test" {
  name                       = "acctestLASS-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  category     = "Saved Search Test Category"
  display_name = "Create or Update Saved Search Test"
  query        = "Heartbeat | summarize Count() by Computer | take a"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsSavedSearch_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsSavedSearch_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_saved_search" "import" {
  name                       = azurerm_log_analytics_saved_search.test.name
  log_analytics_workspace_id = azurerm_log_analytics_saved_search.test.log_analytics_workspace_id

  category     = azurerm_log_analytics_saved_search.test.category
  display_name = azurerm_log_analytics_saved_search.test.display_name
  query        = azurerm_log_analytics_saved_search.test.query
}
`, template)
}

func testAccAzureRMLogAnalyticsSavedSearch_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_saved_search" "test" {
  name                       = "acctestLASS-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  category     = "Saved Search Test Category"
  display_name = "Create or Update Saved Search Test"
  query        = "Heartbeat | summarize Count() by Computer | take a"

  function_alias      = "heartbeat_func"
  function_parameters = ["a:int=1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsSavedSearch_withTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_saved_search" "test" {
  name                       = "acctestLASS-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  category     = "Saved Search Test Category"
  display_name = "Create or Update Saved Search Test"
  query        = "Heartbeat | summarize Count() by Computer | take a"
  tags = {
    "Environment" = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
