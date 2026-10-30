// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_custom_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogicAppTriggerCustomResource struct{}

func TestAccLogicAppTriggerCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")
	r := LogicAppTriggerCustomResource{}

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

func TestAccLogicAppTriggerCustom_callbackUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")
	r := LogicAppTriggerCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nonEmptyCallback(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("callback_url").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppTriggerCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_trigger_custom", "test")
	r := LogicAppTriggerCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_trigger_custom"),
		},
	})
}

func (LogicAppTriggerCustomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workflowtriggers.ParseTriggerID(state.ID)
	if err != nil {
		return nil, err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroupName, id.WorkflowName)
	resp, err := clients.Logic.WorkflowClient.Get(ctx, workflowId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", workflowId, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Definition == nil {
		return pointer.To(false), nil
	}

	definition, ok := (*resp.Model.Properties.Definition).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type for `definition` of %s: %T", workflowId, *resp.Model.Properties.Definition)
	}

	raw, ok := definition["triggers"]
	if !ok || raw == nil {
		return pointer.To(false), nil
	}

	triggers, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type for `triggers` in definition of %s: %T", workflowId, raw)
	}

	for k := range triggers {
		if strings.EqualFold(k, id.TriggerName) {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (LogicAppTriggerCustomResource) basic(data acceptance.TestData) string {
	template := LogicAppTriggerCustomResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "recurrence-%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
  "recurrence": {
    "frequency": "Day",
    "interval": 1
  },
  "type": "Recurrence"
}
BODY

}
`, template, data.RandomInteger)
}

func (LogicAppTriggerCustomResource) nonEmptyCallback(data acceptance.TestData) string {
	template := LogicAppTriggerCustomResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "request-%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
  "inputs": {
    "method": "POST",
    "relativePath": "customers/{id}",
	"schema": {}
  },
  "kind": "Http",
  "type": "Request"
}
BODY

}
`, template, data.RandomInteger)
}

func (LogicAppTriggerCustomResource) requiresImport(data acceptance.TestData) string {
	template := LogicAppTriggerCustomResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_trigger_custom" "import" {
  name         = azurerm_logic_app_trigger_custom.test.name
  logic_app_id = azurerm_logic_app_trigger_custom.test.logic_app_id
  body         = azurerm_logic_app_trigger_custom.test.body
}
`, template)
}

func (LogicAppTriggerCustomResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
