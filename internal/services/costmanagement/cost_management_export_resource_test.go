package costmanagement_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CostManagementExportResource struct {
}

func TestAccCostManagementExport_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export", "test")
	r := CostManagementExportResource{}

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

func TestAccCostManagementExport_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export", "test")
	r := CostManagementExportResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}
func TestAccCostManagementExport_subscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export", "test")
	r := CostManagementExportResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCostManagementExport_managementgroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export", "test")
	r := CostManagementExportResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managementgroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCostManagementExport_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export", "test")
	r := CostManagementExportResource{}

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
			Config: r.subscription(data),
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

func (t CostManagementExportResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CostManagementExportID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.CostManagement.ExportClient.Get(ctx, id.ResourceId, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Cost Management Export ResourceGroup %q (resource group: %q) does not exist", id.Name, id.ResourceId)
	}

	return utils.Bool(resp.ExportProperties != nil), nil
}

func (r CostManagementExportResource) basic(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	return fmt.Sprintf(`
%s

resource "azurerm_cost_management_export" "test" {
  name                    = "accrg%d"
  scope                   = azurerm_resource_group.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  delivery_info {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root"
  }

  query {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, r.template(data), data.RandomInteger, start, end)
}

func (r CostManagementExportResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cost_management_export" "import" {
  name                    = azurerm_cost_management_export.test.name
  scope                   = azurerm_cost_management_export.test.scope
  recurrence_type         = azurerm_cost_management_export.test.recurrence_type
  recurrence_period_start = azurerm_cost_management_export.test.recurrence_period_start
  recurrence_period_end   = azurerm_cost_management_export.test.recurrence_period_end

  delivery_info {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root"
  }

  query {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, r.basic(data))
}

func (r CostManagementExportResource) subscription(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	return fmt.Sprintf(`
%s

data "azurerm_subscription" "current" {
}

resource "azurerm_cost_management_export" "test" {
  name                    = "accrg%d"
  scope                   = data.azurerm_subscription.current.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  delivery_info {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root"
  }

  query {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, r.template(data), data.RandomInteger, start, end)
}

func (r CostManagementExportResource) managementgroup(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	return fmt.Sprintf(`
%s

resource "azurerm_management_group" "test" {
  display_name = "testGroup"
  name         = "testGroup"
}

resource "azurerm_cost_management_export" "test" {
  name                    = "accrg%d"
  scope                   = azurerm_management_group.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  delivery_info {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root"
  }

  query {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, r.template(data), data.RandomInteger, start, end)
}

func (r CostManagementExportResource) update(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	end := time.Now().AddDate(0, 4, 0).Format("2006-01-02")

	return fmt.Sprintf(`
%s

resource "azurerm_cost_management_export" "test" {
  name                    = "accrg%d"
  scope                   = azurerm_resource_group.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  delivery_info {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root/updated"
  }

  query {
    type       = "Usage"
    time_frame = "WeekToDate"
  }
}
`, r.template(data), data.RandomInteger, start, end)
}

func (CostManagementExportResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}