package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/costmanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCostManagementExportResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCostManagementExportResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCostManagementExportResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCostManagementExportResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCostManagementExportResourceGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCostManagementExportResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCostManagementExportResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCostManagementExportResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCostManagementExportResourceGroup_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCostManagementExportResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCostManagementExportResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCostManagementExportResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCostManagementExportResourceGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).CostManagement.ExportClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.CostManagementExportResourceGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceId, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on costManagementExportResourceGroupClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Cost Management Export ResourceGroup %q (resource group: %q) does not exist", id.Name, id.ResourceId)
		}

		return nil
	}
}

func testCheckAzureRMCostManagementExportResourceGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).CostManagement.ExportClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cost_management_export_resource_group" {
			continue
		}

		id, err := parse.CostManagementExportResourceGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceId, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Cost Management Export ResourceGroup still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMCostManagementExportResourceGroup_basic(data acceptance.TestData) string {
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

resource "azurerm_cost_management_export_resource_group" "test" {
  name                    = "accrg%d"
  resource_group_id       = azurerm_resource_group.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "2020-06-18T00:00:00Z"
  recurrence_period_end   = "2020-07-18T00:00:00Z"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMCostManagementExportResourceGroup_update(data acceptance.TestData) string {
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

resource "azurerm_cost_management_export_resource_group" "test" {
  name                    = "accrg%d"
  resource_group_id       = azurerm_resource_group.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "2020-08-18T00:00:00Z"
  recurrence_period_end   = "2020-09-18T00:00:00Z"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
