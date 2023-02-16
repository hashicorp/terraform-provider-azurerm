package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelDataConnectorApiPollingResource struct{}

func TestAccAzureRMSentinelDataConnectorApiPolling_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_api_polling", "test")
	r := SentinelDataConnectorApiPollingResource{}

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

func TestAccAzureRMSentinelDataConnectorApiPolling_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_api_polling", "test")
	r := SentinelDataConnectorApiPollingResource{}

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

func TestAccAzureRMSentinelDataConnectorApiPolling_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_api_polling", "test")
	r := SentinelDataConnectorApiPollingResource{}

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

func (r SentinelDataConnectorApiPollingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.DataConnectorsClient

	id, err := parse.DataConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SentinelDataConnectorApiPollingResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_api_polling" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  active                    = true

  auth {
    type               = "APIKey"
    api_key_name       = "Authorization"
    api_key_identifier = "Bearer"
  }

  request {
    api_endpoint              = "https://api.slack.com/audit/v1/logs"
    http_method               = "GET"
    query_time_format         = "UnixTimestamp"
    start_time_attribute_name = "oldest"
    end_time_attribute_name   = "latest"
    query_window_in_minutes   = 5
    query_parameters = jsonencode({
      "limit" = 1000
    })
  }

  paging {
    type                      = "PageToken"
    next_page_parameter_name  = "cursor"
    next_page_token_json_path = "..response_metadata.next_cursor"
  }

  response {
    events_json_paths = ["$..entries"]
  }

  ui {
    title                    = "Slack"
    publisher                = "Slack"
    description_markdown     = "The [Slack](https://slack.com) data connector provides the capability to ingest [Slack Audit Records](https://api.slack.com/admins/audit-logs) events into Microsoft Sentinel through the REST API. Refer to [API documentation](https://api.slack.com/admins/audit-logs#the_audit_event) for more information. The connector provides ability to get events which helps to examine potential security risks, analyze your team's use of collaboration, diagnose configuration problems and more. This data connector uses Microsoft Sentinel native polling capability."
    graph_queries_table_name = "SlackAuditNativePoller_CL"
    graph_query {
      metric_name = "Total data received"
      legend      = "Slack audit events"
      base_query  = "{{graphQueriesTableName}}"
    }

    sample_query {
      description = "All Slack audit events"
      query       = "{{graphQueriesTableName}}\n| sort by TimeGenerated desc"
    }

    data_type {
      name                     = "{{graphQueriesTableName}}"
      last_data_received_query = "{{graphQueriesTableName}}\n            | summarize Time = max(TimeGenerated)\n            | where isnotempty(Time)"
    }

    connectivity_criteria {
      type = "IsConnectedQuery"
    }

    availability {
      enabled = true
      preview = true
    }

    permission {
      resource_provider {
        name                     = "Microsoft.OperationalInsights/workspaces"
        display_name = "read and write permissions are required."
        display_text = "Workspace"
        scope                    = "Workspace"
        required_permissions {
          read   = true
          write  = true
          delete = true
        }
      }
      custom {
        name        = "Slack API credentials"
        description = "**SlackAPIBearerToken** is required for REST API. [See the documentation to learn more about API](https://api.slack.com/web#authentication). Check all [requirements and follow the instructions](https://api.slack.com/web#authentication) for obtaining credentials."
      }
    }

    instruction {
      title       = "Connect Slack to Microsoft Sentinel"
      description = "Enable Slack audit Logs"
      step {
        type = "InfoMessage"
        parameters = jsonencode({
          "enable" = "true"
        })
      }
    }
  }

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorApiPollingResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_sentinel_data_connector_azure_active_directory" "test" {
  name                       = "accTestDC-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  tenant_id                  = data.azurerm_client_config.test.tenant_id
  depends_on                 = [azurerm_log_analytics_solution.test]
}
`, template, data.RandomInteger)
}

func (r SentinelDataConnectorApiPollingResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_azure_active_directory" "import" {
  name                       = azurerm_sentinel_data_connector_azure_active_directory.test.name
  log_analytics_workspace_id = azurerm_sentinel_data_connector_azure_active_directory.test.log_analytics_workspace_id
}
`, template)
}

func (r SentinelDataConnectorApiPollingResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
