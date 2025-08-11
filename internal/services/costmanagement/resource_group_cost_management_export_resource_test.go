// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupCostManagementExport struct{}

func TestAccResourceGroupCostManagementExport_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_export", "test")
	r := ResourceGroupCostManagementExport{}

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

func TestAccResourceGroupCostManagementExport_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_export", "test")
	r := ResourceGroupCostManagementExport{}

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

func TestAccResourceGroupCostManagementExport_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_cost_management_export", "test")
	r := ResourceGroupCostManagementExport{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_resource_group_cost_management_export"),
		},
	})
}

func (t ResourceGroupCostManagementExport) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := exports.ParseScopedExportID(state.ID)
	if err != nil {
		return nil, err
	}

	var opts exports.GetOperationOptions
	resp, err := clients.CostManagement.ExportClient.Get(ctx, *id, opts)
	if err != nil {
		return nil, fmt.Errorf("retrieving (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ResourceGroupCostManagementExport) basic(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctestcontainer%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_resource_group_cost_management_export" "test" {
  name                         = "accrg%d"
  resource_group_id            = azurerm_resource_group.test.id
  recurrence_type              = "Monthly"
  recurrence_period_start_date = "%sT00:00:00Z"
  recurrence_period_end_date   = "%sT00:00:00Z"

  export_data_storage_location {
    container_id     = azurerm_storage_container.test.resource_manager_id
    root_folder_path = "/root"
  }
  export_data_options {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, start, end)
}

func (ResourceGroupCostManagementExport) update(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	end := time.Now().AddDate(0, 4, 0).Format("2006-01-02")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctestcontainer%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_resource_group_cost_management_export" "test" {
  name                         = "accrg%d"
  resource_group_id            = azurerm_resource_group.test.id
  recurrence_type              = "Monthly"
  recurrence_period_start_date = "%sT00:00:00Z"
  recurrence_period_end_date   = "%sT00:00:00Z"
  file_format                  = "Csv"

  export_data_storage_location {
    container_id     = azurerm_storage_container.test.resource_manager_id
    root_folder_path = "/root"
  }
  export_data_options {
    type       = "Usage"
    time_frame = "WeekToDate"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, start, end)
}

func (ResourceGroupCostManagementExport) requiresImport(data acceptance.TestData) string {
	template := ResourceGroupCostManagementExport{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_cost_management_export" "import" {
  name                         = azurerm_resource_group_cost_management_export.test.name
  resource_group_id            = azurerm_resource_group.test.id
  recurrence_type              = azurerm_resource_group_cost_management_export.test.recurrence_type
  recurrence_period_start_date = azurerm_resource_group_cost_management_export.test.recurrence_period_start_date
  recurrence_period_end_date   = azurerm_resource_group_cost_management_export.test.recurrence_period_start_date

  export_data_storage_location {
    container_id     = azurerm_storage_container.test.resource_manager_id
    root_folder_path = "/root"
  }

  export_data_options {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, template)
}
