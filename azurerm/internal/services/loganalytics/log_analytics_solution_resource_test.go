package loganalytics_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMLogAnalyticsSolution_basicContainerMonitoring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSolutionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSolution_containerMonitoring(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSolutionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsSolution_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSolutionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSolution_containerMonitoring(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSolutionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogAnalyticsSolution_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_solution"),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsSolution_basicSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsSolutionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsSolution_security(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsSolutionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogAnalyticsSolutionDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.SolutionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_solution" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics solution still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsSolutionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.SolutionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Workspace: %q", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Solutions Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Solutions %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsSolution_containerMonitoring(data acceptance.TestData) string {
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

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "ContainerInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsSolution_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsSolution_containerMonitoring(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_solution" "import" {
  solution_name         = azurerm_log_analytics_solution.test.solution_name
  location              = azurerm_log_analytics_solution.test.location
  resource_group_name   = azurerm_log_analytics_solution.test.resource_group_name
  workspace_resource_id = azurerm_log_analytics_solution.test.workspace_resource_id
  workspace_name        = azurerm_log_analytics_solution.test.workspace_name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }
}
`, template)
}

func testAccAzureRMLogAnalyticsSolution_security(data acceptance.TestData) string {
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

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "Security"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Security"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
