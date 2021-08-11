package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelAutomationRuleResource struct{ uuid string }

func TestAccSentinelAutomationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_automation_rule", "test")
	r := SentinelAutomationRuleResource{uuid: uuid.New().String()}

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

func TestAccSentinelAutomationRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_automation_rule", "test")
	r := SentinelAutomationRuleResource{uuid: uuid.New().String()}

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

func TestAccSentinelAutomationRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_automation_rule", "test")
	r := SentinelAutomationRuleResource{uuid: uuid.New().String()}

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

func TestAccSentinelAutomationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_automation_rule", "test")
	r := SentinelAutomationRuleResource{uuid: uuid.New().String()}

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

func (r SentinelAutomationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.AutomationRulesClient

	id, err := parse.AutomationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, sentinel.OperationalInsightsResourceProvider, id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SentinelAutomationRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_automation_rule" "test" {
  name                       = "%s"
  log_analytics_workspace_id = azurerm_log_analytics_solution.sentinel.workspace_resource_id
  display_name               = "acctest-SentinelAutoRule-%d"
  order                      = 1

  action_incident {
    order  = 1
    status = "Active"
  }
}
`, template, r.uuid, data.RandomInteger)
}

func (r SentinelAutomationRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "securityinsights" {
  display_name = "Azure Security Insights"
}

resource "azurerm_role_assignment" "sentinel" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Azure Sentinel Automation Contributor"
  principal_id         = data.azuread_service_principal.securityinsights.object_id
}

resource "azurerm_template_deployment" "testconnection" {
  name                = "testconnection"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Incremental"
  template_body       = <<TEMPLATE
	{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": {
			"connections_azuresentinel_name": {
				"defaultValue": "azuresentinel",
				"type": "String"
			}
		},
		"variables": {},
		"resources": [
			{
				"type": "Microsoft.Web/connections",
				"apiVersion": "2016-06-01",
				"name": "[parameters('connections_azuresentinel_name')]",
				"location": "${azurerm_resource_group.test.location}",
				"kind": "V1",
				"properties": {
					"displayName": "test",
					"customParameterValues": {},
					"api": {
						"id": "[concat('/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Web/locations/${azurerm_resource_group.test.location}/managedApis/', parameters('connections_azuresentinel_name'))]"
					}
				}
			}
		]
	}
  TEMPLATE
}

resource "azurerm_template_deployment" "testlogicapp" {
  name                = "testlogicapp"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Incremental"
  template_body       = <<TEMPLATE
  {
	  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
	  "contentVersion": "1.0.0.0",
	  "parameters": {
		  "workflows_Test_name": {
			  "defaultValue": "Test",
			  "type": "String"
		  },
		  "connections_azuresentinel_externalid": {
			  "defaultValue": "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.Web/connections/azuresentinel",
			  "type": "String"
		  }
	  },
	  "variables": {},
	  "resources": [
		  {
			  "type": "Microsoft.Logic/workflows",
			  "apiVersion": "2017-07-01",
			  "name": "[parameters('workflows_Test_name')]",
			  "location": "${azurerm_resource_group.test.location}",
			  "properties": {
				  "state": "Enabled",
				  "definition": {
					  "$schema": "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
					  "contentVersion": "1.0.0.0",
					  "parameters": {
						  "$connections": {
							  "defaultValue": {},
							  "type": "Object"
						  }
					  },
					  "triggers": {
						  "When_Azure_Sentinel_incident_creation_rule_was_triggered": {
							  "type": "ApiConnectionWebhook",
							  "inputs": {
								  "body": {
									  "callback_url": "@{listCallbackUrl()}"
								  },
								  "host": {
									  "connection": {
										  "name": "@parameters('$connections')['azuresentinel']['connectionId']"
									  }
								  },
								  "path": "/incident-creation"
							  }
						  }
					  },
					  "outputs": {}
				  },
				  "parameters": {
					  "$connections": {
						  "value": {
							  "azuresentinel": {
								  "connectionId": "[parameters('connections_azuresentinel_externalid')]",
								  "connectionName": "azuresentinel",
								  "id": "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Web/locations/${azurerm_resource_group.test.location}/managedApis/azuresentinel"
							  }
						  }
					  }
				  }
			  }
		  }
	  ]
  }
  TEMPLATE

  depends_on = [azurerm_template_deployment.testconnection]
}

resource "azurerm_sentinel_automation_rule" "test" {
  name                       = "%s"
  log_analytics_workspace_id = azurerm_log_analytics_solution.sentinel.workspace_resource_id
  display_name               = "acctest-SentinelAutoRule-%d-update"
  order                      = 2
  enabled                    = false
  condition {
    property = "IncidentTitle"
    operator = "Contains"
    values   = ["a", "b"]
  }

  condition {
    property = "IncidentTitle"
    operator = "Contains"
    values   = ["c", "d"]
  }

  action_incident {
    order                  = 1
    status                 = "Closed"
    classification         = "BenignPositive_SuspiciousButExpected"
    classification_comment = "whatever reason"
  }

  action_incident {
    order  = 3
    labels = ["foo", "bar"]
  }

  action_incident {
    order    = 2
    severity = "High"
  }

  action_incident {
    order    = 4
    owner_id = data.azurerm_client_config.current.object_id
  }

  action_playbook {
    order        = 5
    logic_app_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.Logic/workflows/Test"
  }

  depends_on = [azurerm_role_assignment.sentinel, azurerm_template_deployment.testlogicapp]
}
`, template, r.uuid, data.RandomInteger)
}

func (r SentinelAutomationRuleResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_automation_rule" "import" {
  name                       = azurerm_sentinel_automation_rule.test.name
  log_analytics_workspace_id = azurerm_sentinel_automation_rule.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_automation_rule.test.display_name
  order                      = azurerm_sentinel_automation_rule.test.order
  action_incident {
    order = 1
  }
}
`, template)
}

func (r SentinelAutomationRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = %q
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-workspace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "pergb2018"
}

resource "azurerm_log_analytics_solution" "sentinel" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
