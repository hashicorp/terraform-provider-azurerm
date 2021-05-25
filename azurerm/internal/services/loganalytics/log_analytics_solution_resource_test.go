package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsSolutionResource struct {
}

func TestAccLogAnalyticsSolution_basicContainerMonitoring(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerMonitoring(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsSolution_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerMonitoring(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.security(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsSolutionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["solutions"]

	resp, err := clients.LogAnalytics.SolutionsClient.Get(ctx, resGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Log Analytics Solution (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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

func (r LogAnalyticsSolutionResource) requiresImport(data acceptance.TestData) string {
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
`, r.containerMonitoring(data))
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
