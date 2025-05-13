// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupCostManagementView struct{}

func TestAccResourceGroupCostManagementView_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_view", "test")
	r := ResourceGroupCostManagementView{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupCostManagementView_table(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_view", "test")
	r := ResourceGroupCostManagementView{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.table(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupCostManagementView_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_view", "test")
	r := ResourceGroupCostManagementView{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupCostManagementView_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_view", "test")
	r := ResourceGroupCostManagementView{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_resource_group_cost_management_view"),
		},
	})
}

func (t ResourceGroupCostManagementView) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := views.ParseScopedViewID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.CostManagement.ViewsClient.GetByScope(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ResourceGroupCostManagementView) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_resource_group_cost_management_view" "test" {
  name              = "testcostview%s"
  resource_group_id = azurerm_resource_group.test.id
  chart_type        = "StackedColumn"
  display_name      = "Test View %s"

  accumulated = "false"
  report_type = "Usage"
  timeframe   = "MonthToDate"

  dataset {
    granularity = "Monthly"
    sorting {
      direction = "Ascending"
      name      = "BillingMonth"
    }
    grouping {
      name = "ResourceGroupName"
      type = "Dimension"
    }
    aggregation {
      name        = "totalCost"
      column_name = "Cost"
    }
    aggregation {
      name        = "totalCostUSD"
      column_name = "CostUSD"
    }
  }

  kpi {
    type = "Forecast"
  }
  pivot {
    type = "Dimension"
    name = "ServiceName"
  }
  pivot {
    type = "Dimension"
    name = "ResourceLocation"
  }
  pivot {
    type = "Dimension"
    name = "ResourceGroupName"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (ResourceGroupCostManagementView) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_resource_group_cost_management_view" "test" {
  name              = "testcostview%s"
  resource_group_id = azurerm_resource_group.test.id
  chart_type        = "Line"
  display_name      = "Test View 2 %s"

  accumulated = "false"
  report_type = "Usage"
  timeframe   = "YearToDate"

  dataset {
    granularity = "Daily"
    aggregation {
      name        = "totalCost"
      column_name = "Cost"
    }
    aggregation {
      name        = "totalCostUSD"
      column_name = "CostUSD"
    }
  }

  kpi {
    type = "Forecast"
  }
  pivot {
    type = "Dimension"
    name = "ResourceLocation"
  }
  pivot {
    type = "Dimension"
    name = "ResourceGroupName"
  }
  pivot {
    type = "Dimension"
    name = "ServiceName"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (ResourceGroupCostManagementView) requiresImport(data acceptance.TestData) string {
	template := ResourceGroupCostManagementView{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_cost_management_view" "import" {
  name              = azurerm_resource_group_cost_management_view.test.name
  resource_group_id = azurerm_resource_group_cost_management_view.test.resource_group_id
  chart_type        = azurerm_resource_group_cost_management_view.test.chart_type
  display_name      = azurerm_resource_group_cost_management_view.test.display_name

  accumulated = azurerm_resource_group_cost_management_view.test.accumulated

  report_type = "Usage"
  timeframe   = "MonthToDate"

  dataset {
    granularity = "Monthly"
    sorting {
      direction = "Ascending"
      name      = "BillingMonth"
    }
    grouping {
      name = "ResourceGroupName"
      type = "Dimension"
    }
    aggregation {
      name        = "totalCost"
      column_name = "Cost"
    }
    aggregation {
      name        = "totalCostUSD"
      column_name = "CostUSD"
    }
  }

  kpi {
    type = "Forecast"
  }
  pivot {
    type = "Dimension"
    name = "ServiceName"
  }
  pivot {
    type = "Dimension"
    name = "ResourceLocation"
  }
  pivot {
    type = "Dimension"
    name = "ResourceGroupName"
  }
}
`, template)
}

func (ResourceGroupCostManagementView) table(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_resource_group_cost_management_view" "test" {
  name              = "testcostview%s"
  resource_group_id = azurerm_resource_group.test.id
  chart_type        = "Table"
  display_name      = "Test View %s"

  accumulated = "false"
  report_type = "Usage"
  timeframe   = "MonthToDate"

  dataset {
    granularity = "Monthly"
    sorting {
      direction = "Ascending"
      name      = "BillingMonth"
    }
    grouping {
      name = "ResourceGroupName"
      type = "Dimension"
    }
    aggregation {
      name        = "totalCost"
      column_name = "Cost"
    }
    aggregation {
      name        = "totalCostUSD"
      column_name = "CostUSD"
    }
  }

  kpi {
    type = "Forecast"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
