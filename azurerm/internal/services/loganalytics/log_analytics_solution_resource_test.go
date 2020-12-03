package loganalytics_test

import (
	`context`
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils`
)

type LogAnalyticsSolutionResource struct {
}

func TestAccLogAnalyticsSolution_basicContainerMonitoring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.containerMonitoring(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsSolution_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.containerMonitoring(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_solution"),
		},
	})
}

func TestAccLogAnalyticsSolution_basicSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.security(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogAnalyticsSolutionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	name := id.Path["solutions"]

	resp, err := clients.LogAnalytics.SolutionsClient.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Log Analytics Solutions %q (resource group: %q): %v", name, id.ResourceGroup)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (LogAnalyticsSolutionResource) containerMonitoring(data acceptance.TestData) string {
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

func (LogAnalyticsSolutionResource) requiresImport(data acceptance.TestData) string {
	template := LogAnalyticsSolutionResource{}.containerMonitoring(data)
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

func (LogAnalyticsSolutionResource) security(data acceptance.TestData) string {
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
