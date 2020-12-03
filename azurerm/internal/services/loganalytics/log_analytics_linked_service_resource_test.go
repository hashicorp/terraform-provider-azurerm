package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsLinkedServiceResource struct {
}

func TestAccLogAnalyticsLinkedService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestlaw-%d/Automation", data.RandomInteger)),
				check.That(data.ResourceName).Key("workspace_name").HasValue(fmt.Sprintf("acctestlaw-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("linked_service_name").HasValue("automation"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsLinkedService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestlaw-%d/Automation", data.RandomInteger)),
				check.That(data.ResourceName).Key("workspace_name").HasValue(fmt.Sprintf("acctestlaw-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("linked_service_name").HasValue("automation"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_linked_service"),
		},
	})
}

func TestAccLogAnalyticsLinkedService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("linked_service_name").HasValue("automation"),
			),
		},
		data.ImportStep(),
	})
}

func (LogAnalyticsLinkedServiceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resGroup := id.ResourceGroup
	workspaceName := id.Path["workspaces"]
	lsName := id.Path["linkedservices"]

	resp, err := clients.LogAnalytics.LinkedServicesClient.Get(ctx, resGroup, workspaceName, lsName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Log Analytics Linked Service %s (resource group: %s): %v", lsName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.LinkedServiceProperties != nil), nil
}

func (r LogAnalyticsLinkedServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  resource_id         = azurerm_automation_account.test.id
}
`, r.template(data))
}

func (r LogAnalyticsLinkedServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "import" {
  resource_group_name = azurerm_log_analytics_linked_service.test.resource_group_name
  workspace_name      = azurerm_log_analytics_linked_service.test.workspace_name
  resource_id         = azurerm_log_analytics_linked_service.test.resource_id
}
`, r.basic(data))
}

func (r LogAnalyticsLinkedServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  linked_service_name = "automation"
  resource_id         = azurerm_automation_account.test.id
}
`, r.template(data))
}

func (LogAnalyticsLinkedServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
