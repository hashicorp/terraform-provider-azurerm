package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMlogAnalyticsLinkedStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsLinkedStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMlogAnalyticsLinkedStorageAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsLinkedStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMlogAnalyticsLinkedStorageAccount_requiresImport),
		},
	})
}

func TestAccAzureRMlogAnalyticsLinkedStorageAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsLinkedStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMlogAnalyticsLinkedStorageAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMoperationalinsightsLinkedStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMoperationalinsightsLinkedStorageAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.LinkedStorageAccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("operationalinsights LinkedStorageAccount not found: %s", resourceName)
		}
		id, err := parse.LogAnalyticsLinkedStorageAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, operationalinsights.DataSourceType(id.Name)); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Operationalinsights LinkedStorageAccount %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Operationalinsights.LinkedStorageAccountClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMoperationalinsightsLinkedStorageAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_linked_storage_account" {
			continue
		}
		id, err := parse.LogAnalyticsLinkedStorageAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, operationalinsights.DataSourceType(id.Name)); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on LogAnalytics.LinkedStorageAccountClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMlogAnalyticsLinkedStorageAccount_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsap%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data acceptance.TestData) string {
	template := testAccAzureRMlogAnalyticsLinkedStorageAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_storage_account" "test" {
  data_source_type      = "customlogs"
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  storage_account_ids   = [azurerm_storage_account.test.id]
}
`, template)
}

func testAccAzureRMlogAnalyticsLinkedStorageAccount_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMlogAnalyticsLinkedStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_storage_account" "import" {
  data_source_type      = azurerm_log_analytics_linked_storage_account.test.data_source_type
  resource_group_name   = azurerm_log_analytics_linked_storage_account.test.resource_group_name
  workspace_resource_id = azurerm_log_analytics_linked_storage_account.test.workspace_resource_id
  storage_account_ids   = [azurerm_storage_account.test.id]
}
`, config)
}

func testAccAzureRMlogAnalyticsLinkedStorageAccount_complete(data acceptance.TestData) string {
	template := testAccAzureRMlogAnalyticsLinkedStorageAccount_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsas%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_log_analytics_linked_storage_account" "test" {
  data_source_type      = "customlogs"
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  storage_account_ids   = [azurerm_storage_account.test.id, azurerm_storage_account.test2.id]
}
`, template, data.RandomString)
}
