package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsLinkedStorageAccountResource struct {
}

func TestAcclogAnalyticsLinkedStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	r := LogAnalyticsLinkedStorageAccountResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAcclogAnalyticsLinkedStorageAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	r := LogAnalyticsLinkedStorageAccountResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAcclogAnalyticsLinkedStorageAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	r := LogAnalyticsLinkedStorageAccountResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAcclogAnalyticsLinkedStorageAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_storage_account", "test")
	r := LogAnalyticsLinkedStorageAccountResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsLinkedStorageAccountResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsLinkedStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}
	dataSourceType := operationalinsights.DataSourceType(id.LinkedStorageAccountName)

	resp, err := clients.LogAnalytics.LinkedStorageAccountClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, dataSourceType)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Linked Service Storage Account (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LogAnalyticsLinkedStorageAccountResource) template(data acceptance.TestData) string {
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

func (r LogAnalyticsLinkedStorageAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_storage_account" "test" {
  data_source_type      = "customlogs"
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  storage_account_ids   = [azurerm_storage_account.test.id]
}
`, r.template(data))
}

func (r LogAnalyticsLinkedStorageAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_storage_account" "import" {
  data_source_type      = azurerm_log_analytics_linked_storage_account.test.data_source_type
  resource_group_name   = azurerm_log_analytics_linked_storage_account.test.resource_group_name
  workspace_resource_id = azurerm_log_analytics_linked_storage_account.test.workspace_resource_id
  storage_account_ids   = [azurerm_storage_account.test.id]
}
`, r.basic(data))
}

func (r LogAnalyticsLinkedStorageAccountResource) complete(data acceptance.TestData) string {
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
`, r.template(data), data.RandomString)
}
