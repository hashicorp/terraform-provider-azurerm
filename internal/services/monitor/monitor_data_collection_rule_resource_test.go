package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/sdk/2021-04-01/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorDataCollectionRuleResource struct{}

func (r MonitorDataCollectionRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := datacollectionrules.ParseDataCollectionRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Monitor.DataCollectionRulesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

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

func TestAccApiManagementNotificationRecipientEmail_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule", "test")
	r := MonitorDataCollectionRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
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
		{
			Config: r.update(data),
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

func (r MonitorDataCollectionRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  destinations {
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }
  data_flows {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_log_analytics_workspace" "test1" {
  name                = "acctestlaw1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  data_sources {
    syslog {
      facility_names = ["*"]
      log_levels     = ["*"]
      name           = "test-datasource-syslog"
    }
    performance_counters {
      streams                       = ["Microsoft-Perf", "Microsoft-InsightsMetrics"]
      sampling_frequency_in_seconds = 10
      counter_specifiers            = ["Processor(*)\\%% Processor Time"]
      name                          = "test-datasource-perfcounter"
    }
  }
  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test1.id
      name                  = "test-destination-log1"
    }
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }

  }
  data_flows {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
  data_flows {
    streams      = ["Microsoft-InsightsMetrics", "Microsoft-Syslog", "Microsoft-Perf"]
    destinations = ["test-destination-log1"]
  }
  kind        = "Linux"
  description = "acc test monitor_data_collection_rule"
  tags = {
    ENV = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_log_analytics_workspace" "test1" {
  name                = "acctestlaw1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctestlaw2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  data_sources {
    syslog {
      facility_names = [
        "auth",
        "authpriv",
        "cron",
        "daemon",
        "kern",
        "lpr",
        "mail",
        "mark",
        "news",
        "syslog",
        "user",
        "uucp",
        "local0",
        "local1",
        "local2",
        "local3",
        "local4",
        "local5",
        "local6",
        "local7",
      ]
      log_levels = [
        "Debug",
        "Info",
        "Notice",
        "Warning",
        "Error",
        "Critical",
        "Alert",
        "Emergency",
      ]
      name = "test-datasource-syslog"
    }
    performance_counters {
      streams                       = ["Microsoft-Perf", "Microsoft-InsightsMetrics"]
      sampling_frequency_in_seconds = 10
      counter_specifiers = [
        "Processor(*)\\%% Processor Time",
        "Processor(*)\\%% Idle Time",
        "Processor(*)\\%% User Time",
        "Processor(*)\\%% Nice Time",
        "Processor(*)\\%% Privileged Time",
        "Processor(*)\\%% IO Wait Time",
        "Processor(*)\\%% Interrupt Time",
        "Processor(*)\\%% DPC Time",
      ]
      name = "test-datasource-perfcounter"
    }
    performance_counters {
      streams                       = ["Microsoft-Perf"]
      sampling_frequency_in_seconds = 20
      counter_specifiers = [
        "Memory(*)\\Available MBytes Memory",
        "Memory(*)\\%% Available Memory",
        "Memory(*)\\Used Memory MBytes",
        "Memory(*)\\%% Used Memory",
        "Memory(*)\\Pages/sec",
        "Memory(*)\\Page Reads/sec",
        "Memory(*)\\Page Writes/sec",
        "Memory(*)\\Available MBytes Swap",
        "Memory(*)\\%% Available Swap Space",
        "Memory(*)\\Used MBytes Swap Space",
        "Memory(*)\\%% Used Swap Space",
        "Logical Disk(*)\\%% Free Inodes",
        "Logical Disk(*)\\%% Used Inodes",
        "Logical Disk(*)\\Free Megabytes",
        "Logical Disk(*)\\%% Free Space",
        "Logical Disk(*)\\%% Used Space",
        "Logical Disk(*)\\Logical Disk Bytes/sec",
        "Logical Disk(*)\\Disk Read Bytes/sec",
        "Logical Disk(*)\\Disk Write Bytes/sec",
        "Logical Disk(*)\\Disk Transfers/sec",
        "Logical Disk(*)\\Disk Reads/sec",
        "Logical Disk(*)\\Disk Writes/sec",
        "Network(*)\\Total Bytes Transmitted",
        "Network(*)\\Total Bytes Received",
        "Network(*)\\Total Bytes",
        "Network(*)\\Total Packets Transmitted",
        "Network(*)\\Total Packets Received",
        "Network(*)\\Total Rx Errors",
        "Network(*)\\Total Tx Errors",
        "Network(*)\\Total Collisions"
      ]
      name = "test-datasource-perfcounter2"
    }
    windows_event_logs {
      streams        = ["Microsoft-WindowsEvent"]
      x_path_queries = ["*[System/Level=1]"]
      name           = "test-datasource-wineventlog"
    }
    extensions {
      streams            = ["Microsoft-WindowsEvent"]
      input_data_sources = ["test-datasource-wineventlog"]
      extension_name     = "test-extension-name"
      extension_settings = jsonencode({
        a = 1
        b = "hello"
      })
      name = "test-datasource-extension"
    }
  }
  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test1.id
      name                  = "test-destination-log1"
    }
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test2.id
      name                  = "test-destination-log2"
    }
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }
  data_flows {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
  data_flows {
    streams      = ["Microsoft-InsightsMetrics", "Microsoft-Syslog", "Microsoft-Perf"]
    destinations = ["test-destination-log1"]
  }
  data_flows {
    streams      = ["Microsoft-Event", "Microsoft-WindowsEvent"]
    destinations = ["test-destination-log1", "test-destination-log2"]
  }
  description = "acc test monitor_data_collection_rule complete"
  tags = {
    ENV  = "test"
    ENV2 = "test2"
  }
}

`, r.template(data), data.RandomInteger)
}

func (r MonitorDataCollectionRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_data_collection_rule" "import" {
  name                = azurerm_monitor_data_collection_rule.test.name
  resource_group_name = azurerm_monitor_data_collection_rule.test.resource_group_name
  location            = azurerm_monitor_data_collection_rule.test.location
  destinations        = azurerm_monitor_data_collection_rule.test.destinations
  data_flows          = azurerm_monitor_data_collection_rule.test.data_flows
}
`, r.template(data))
}

func (r MonitorDataCollectionRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-DataCollectionRule-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}
