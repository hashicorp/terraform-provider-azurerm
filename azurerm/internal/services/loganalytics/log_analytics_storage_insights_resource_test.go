package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLogAnalyticsStorageInsights_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsStorageInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"), // key is not returned by the API
		},
	})
}

func TestAccAzureRMLogAnalyticsStorageInsights_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsStorageInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsStorageInsights_requiresImport),
		},
	})
}

func TestAccAzureRMLogAnalyticsStorageInsights_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsStorageInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"), // key is not returned by the API
		},
	})
}

func TestAccAzureRMLogAnalyticsStorageInsights_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsStorageInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"), // key is not returned by the API
		},
	})
}

func TestAccAzureRMLogAnalyticsStorageInsights_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsStorageInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
			{
				Config: testAccAzureRMLogAnalyticsStorageInsights_updateStorageAccount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsStorageInsightsExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"), // key is not returned by the API
		},
	})
}

func testCheckAzureRMLogAnalyticsStorageInsightsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.StorageInsightsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Storage Insights not found: %s", resourceName)
		}
		id, err := parse.LogAnalyticsStorageInsightsID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.StorageInsightConfigName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Log Analytics Storage Insights %q does not exist", id.StorageInsightConfigName)
			}
			return fmt.Errorf("bad: Get on LogAnalytics.StorageInsightsClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMLogAnalyticsStorageInsightsDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_storage_insights" {
			continue
		}
		id, err := parse.LogAnalyticsStorageInsightsID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.StorageInsightConfigName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on LogAnalytics.StorageInsightsClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMLogAnalyticsStorageInsights_template(data acceptance.TestData) string {
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
retention_in_days = 30
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

func testAccAzureRMLogAnalyticsStorageInsights_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsStorageInsights_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-la-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsStorageInsights_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMLogAnalyticsStorageInsights_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "import" {
  name                = azurerm_log_analytics_storage_insights.test.name
  resource_group_name = azurerm_log_analytics_storage_insights.test.resource_group_name
  workspace_id        = azurerm_log_analytics_storage_insights.test.workspace_id

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, config)
}

func testAccAzureRMLogAnalyticsStorageInsights_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsStorageInsights_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  blob_container_names = ["wad-iis-logfiles"]
  table_names          = ["WADWindowsEventLogsTable", "LinuxSyslogVer2v0"]

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsStorageInsights_updateStorageAccount(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsStorageInsights_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-la-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  blob_container_names = ["wad-iis-logfiles"]
  table_names          = ["WADWindowsEventLogsTable", "LinuxSyslogVer2v0"]

  storage_account_id  = azurerm_storage_account.test2.id
  storage_account_key = azurerm_storage_account.test2.primary_access_key
}
`, template, data.RandomStringOfLength(6), data.RandomInteger)
}
