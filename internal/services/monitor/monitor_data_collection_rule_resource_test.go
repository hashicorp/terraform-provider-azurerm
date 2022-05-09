package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorDataCollectionRuleResource struct{}

func TestAccMonitorDataCollectionRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule", "test")
	r := MonitorDataCollectionRuleResource{}

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

func TestAccMonitorDataCollectionRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule", "test")
	r := MonitorDataCollectionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_data_collection_rule"),
		},
	})
}

func TestAccMonitorDataCollectionRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule", "test")
	r := MonitorDataCollectionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDataCollectionRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule", "test")
	r := MonitorDataCollectionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func (r MonitorDataCollectionRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataCollectionRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Monitor.DataCollectionRulesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Monitor DataCollectionRule (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r MonitorDataCollectionRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dcr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "WindowsEventForwarding"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/WindowsEventForwarding"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r MonitorDataCollectionRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctest-dcr-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Windows"
  log_analytics_destination {
    name                  = "centralWorkspace"
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }

  windows_event_log_data_source {
    name          = "cloudSecurityTeamEvents"
    streams       = ["Microsoft-WindowsEvent"]
    xpath_queries = ["Security!"]
  }
  data_flows {
    streams      = ["Microsoft-WindowsEvent"]
    destinations = ["centralWorkspace"]
  }

  depends_on = [
    azurerm_log_analytics_solution.test
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_data_collection_rule" "import" {
  name                = azurerm_monitor_data_collection_rule.test.name
  location            = azurerm_monitor_data_collection_rule.test.location
  resource_group_name = azurerm_monitor_data_collection_rule.test.resource_group_name
  kind                = "Windows"

  log_analytics_destination {
    name                  = "centralWorkspace"
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }

  windows_event_log_data_source {
    name          = "cloudSecurityTeamEvents"
    streams       = ["Microsoft-WindowsEvent"]
    xpath_queries = ["Security!"]
  }
  data_flows {
    streams      = ["Microsoft-Perf", "Microsoft-Syslog", "Microsoft-WindowsEvent"]
    destinations = ["centralWorkspace"]
  }
}
`, r.basic(data))
}

func (r MonitorDataCollectionRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctest-dcr-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Windows"

  azure_monitor_metrics_destination {
    name = "amm1"
  }

  log_analytics_destination {
    name                  = "centralWorkspace"
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }

  windows_event_log_data_source {
    name          = "cloudSecurityTeamEvents"
    streams       = ["Microsoft-WindowsEvent"]
    xpath_queries = ["Security!"]
  }

  windows_event_log_data_source {
    name    = "appTeam1AppEvents"
    streams = ["Microsoft-WindowsEvent"]
    xpath_queries = ["System![System[(Level = 1 or Level = 2 or Level = 3)]]",
    "Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"]
  }

  extension_data_source {
    name               = "extension1"
    extension_name     = "mockname"
    streams            = ["Microsoft-Event"]
    input_data_sources = []
    extension_setting  = <<BODY
{
    "key1": "value1",
    "key2": "value2"
}
BODY
  }

  performance_counter_data_source {
    name    = "cloudTeamCoreCounters"
    streams = ["Microsoft-Perf"]
    specifiers = [
      "\\\\Memory\\\\Committed Bytes",
      "\\\\LogicalDisk(_Total)\\\\Free Megabytes",
      "\\\\PhysicalDisk(_Total)\\\\Avg. Disk Queue Length"
    ]
    sampling_frequency = 15
  }

  performance_counter_data_source {
    name    = "appTeamExtraCounters"
    streams = ["Microsoft-Perf"]
    specifiers = [
      "\\\\Process(_Total)\\\\Thread Count"
    ]
    sampling_frequency = 30
  }

  data_flows {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["amm1"]
  }

  data_flows {
    streams      = ["Microsoft-Perf", "Microsoft-Syslog", "Microsoft-WindowsEvent"]
    destinations = ["centralWorkspace"]
  }

  description = "this is description"

  tags = {
    Environment = "Production"
  }

  depends_on = [
    azurerm_log_analytics_solution.test
  ]
}
`, r.template(data), data.RandomInteger)
}
