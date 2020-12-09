package loganalytics_test

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

func TestAccAzureRMLogAnalyticsDataExportRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataExportRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataExportRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataExportRuleExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataExportRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataExportRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataExportRule_basicLower(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataExportRuleExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			},
			{
				Config:             testAccAzureRMLogAnalyticsDataExportRule_requiresImport(data),
				ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
				ExpectError:        acceptance.RequiresImportError("azurerm_log_analytics_data_export_rule"),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsDataExportRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataExportRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataExportRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataExportRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataExportRule_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataExportRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataExportRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataExportRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataExportRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataExportRuleExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataExportRuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataExportClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_data_export_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		workspace, err := parse.LogAnalyticsWorkspaceID(rs.Primary.Attributes["workspace_resource_id"])
		if err != nil {
			return nil
		}

		resp, err := conn.Get(ctx, resourceGroup, workspace.WorkspaceName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics Data Export Rule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsDataExportRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataExportClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Data Export Rule: %q", name)
		}

		workspace, err := parse.LogAnalyticsWorkspaceID(rs.Primary.Attributes["workspace_resource_id"])
		if err != nil {
			return fmt.Errorf("Bad: unable to access 'workspace_resource_id' for Log Analytics Data Export Rule: %q", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, workspace.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Data Export Rule Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Data Export Rule %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsDataExportRule_template(data acceptance.TestData) string {
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
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func testAccAzureRMLogAnalyticsDataExportRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataExportRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-DER-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, template, data.RandomInteger)
}

// I have to make this a lower case to get the requiresImport test to pass since the RP lowercases everything when it sends the data back to you
func testAccAzureRMLogAnalyticsDataExportRule_basicLower(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataExportRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-der-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataExportRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataExportRule_basicLower(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "import" {
  name                    = azurerm_log_analytics_data_export_rule.test.name
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, template)
}

func testAccAzureRMLogAnalyticsDataExportRule_update(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataExportRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-DER-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat", "Event"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataExportRule_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataExportRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-DER-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
  enabled                 = true
}
`, template, data.RandomInteger)
}
