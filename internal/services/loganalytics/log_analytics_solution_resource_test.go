package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsSolutionResource struct {
}

func TestAccLogAnalyticsSolution_containerMonitoring(t *testing.T) {
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

func TestAccLogAnalyticsSolution_security(t *testing.T) {
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

func TestAccLogAnalyticsSolution_vmInsights(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vmInsights(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsSolution_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_solution", "test")
	r := LogAnalyticsSolutionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
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

func (r LogAnalyticsSolutionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsSolutionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.SolutionsClient.Get(ctx, id.ResourceGroup, id.SolutionName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LogAnalyticsSolutionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogAnalyticsSolutionResource) containerMonitoring(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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
`, r.template(data))
}

func (r LogAnalyticsSolutionResource) security(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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
`, r.template(data))
}

func (r LogAnalyticsSolutionResource) vmInsights(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "VMInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/VMInsights"
  }
}
`, r.template(data))
}

func (r LogAnalyticsSolutionResource) custom(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "acctest-Custom-%[2]s"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/VMInsights"
  }
}
`, r.template(data), data.RandomString)
}

func (r LogAnalyticsSolutionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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
