package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsDataExportRuleResource struct {
}

func TestAccLogAnalyticsDataExportRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")
	r := LogAnalyticsDataExportRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.basic(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataExportRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")
	r := LogAnalyticsDataExportRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.basicLower(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.requiresImport(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			ExpectError:        acceptance.RequiresImportError("azurerm_log_analytics_data_export_rule"),
		},
	})
}

func TestAccLogAnalyticsDataExportRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")
	r := LogAnalyticsDataExportRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.basic(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:             r.update(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataExportRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_data_export_rule", "test")
	r := LogAnalyticsDataExportRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.complete(data),
			ExpectNonEmptyPlan: true, // Due to API changing case of attributes you need to ignore a non-empty plan for this resource
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsDataExportRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsDataExportID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.DataExportClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.DataexportName)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Data Export (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LogAnalyticsDataExportRuleResource) template(data acceptance.TestData) string {
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

func (r LogAnalyticsDataExportRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-DER-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, r.template(data), data.RandomInteger)
}

// I have to make this a lower case to get the requiresImport test to pass since the RP lowercases everything when it sends the data back to you
func (r LogAnalyticsDataExportRuleResource) basicLower(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-der-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataExportRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "import" {
  name                    = azurerm_log_analytics_data_export_rule.test.name
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat"]
}
`, r.basicLower(data))
}

func (r LogAnalyticsDataExportRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_data_export_rule" "test" {
  name                    = "acctest-DER-%d"
  resource_group_name     = azurerm_resource_group.test.name
  workspace_resource_id   = azurerm_log_analytics_workspace.test.id
  destination_resource_id = azurerm_storage_account.test.id
  table_names             = ["Heartbeat", "Event"]
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataExportRuleResource) complete(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}
