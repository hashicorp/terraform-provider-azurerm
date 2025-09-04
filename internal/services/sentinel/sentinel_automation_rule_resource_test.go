// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelAutomationRuleResource struct {
	uuid string
}

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

func TestAccSentinelAutomationRule_trigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_automation_rule", "test")
	r := SentinelAutomationRuleResource{uuid: uuid.New().String()}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This defaults to incident created
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.triggerIncidentUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.triggerAlertCreated(data),
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

	id, err := automationrules.ParseAutomationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
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
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
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
	expDate := time.Now().AddDate(0, 1, 0).UTC().Format(time.RFC3339)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_sentinel_automation_rule" "test" {
  name                       = "%s"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "acctest-SentinelAutoRule-%d-update"
  order                      = 2
  enabled                    = false
  expiration                 = "%s"

  condition_json = jsonencode(
    [
      {
        conditionProperties = {
          operator     = "Contains"
          propertyName = "IncidentTitle"
          propertyValues = [
            "a",
            "b",
          ]
        }
        conditionType = "Property"
      },
      {
        conditionProperties = {
          operator     = "Contains"
          propertyName = "IncidentTitle"
          propertyValues = [
            "c",
            "d",
          ]
        }
        conditionType = "Property"
      },
    ]
  )

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

}
`, template, r.uuid, data.RandomInteger, expDate)
}

func (r SentinelAutomationRuleResource) triggerIncidentUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_sentinel_automation_rule" "test" {
  name                       = "%s"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "acctest-SentinelAutoRule-%d-update"
  order                      = 1
  condition_json = jsonencode([
    {
      conditionType = "PropertyChanged"
      conditionProperties = {
        propertyName   = "IncidentStatus"
        changeType     = "ChangedTo"
        operator       = "Equals"
        propertyValues = ["New"]
      }
    }
  ])

  triggers_when = "Updated"

  action_incident {
    order    = 1
    owner_id = data.azurerm_client_config.current.object_id
  }
}
`, template, r.uuid, data.RandomInteger)
}

func (r SentinelAutomationRuleResource) triggerAlertCreated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_managed_api" "test" {
  name     = "azuresentinel"
  location = azurerm_resource_group.test.location
}

resource "azuread_application" "test" {
  display_name = "acctest-sar-%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.object_id
}


resource "azurerm_api_connection" "test" {
  managed_api_id      = data.azurerm_managed_api.test.id
  name                = "azuresentinel-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  parameter_values = {
    "token:TenantId"     = azuread_service_principal.test.application_tenant_id
    "token:clientId"     = azuread_service_principal.test.client_id
    "token:clientSecret" = azuread_service_principal_password.test.value
    "token:grantType"    = "client_credentials"
  }
  lifecycle {
    ignore_changes = [parameter_values]
  }
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestLogicApp-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  parameters = {
    "$connections" = jsonencode({
      azuresentinel = {
        connectionId         = azurerm_api_connection.test.id
        connectionName       = azurerm_api_connection.test.name
        connectionProperties = {}
        id                   = data.azurerm_managed_api.test.id
      }
    })
  }

  workflow_parameters = {
    "$connections" = jsonencode({
      defaultValue = {}
      type         = "Object"
    })
  }
}

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "Microsoft_Sentinel_alert"
  logic_app_id = azurerm_logic_app_workflow.test.id
  body         = <<BODY
{
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
        "path": "/subscribe"
    }
}
BODY
}


data "azurerm_role_definition" "sentinel" {
  name  = "Microsoft Sentinel Automation Contributor"
  scope = azurerm_resource_group.test.id
}

data "azuread_service_principal" "sentinel" {
  application_id = "98785600-1bb7-4fb9-b9fa-19afe2c8a360"
}

resource "azurerm_role_assignment" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.sentinel.id
  principal_id       = data.azuread_service_principal.sentinel.object_id
}

resource "azurerm_sentinel_automation_rule" "test" {
  name                       = "%[3]s"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "acctest-SentinelAutoRule-%[2]d-update"
  order                      = 1
  triggers_on                = "Alerts"
  action_playbook {
    logic_app_id = azurerm_logic_app_workflow.test.id
    order        = 1
  }

  depends_on = [azurerm_logic_app_trigger_custom.test, azurerm_role_assignment.test]
}




`, template, data.RandomInteger, r.uuid)
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
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
