// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_custom_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogicAppActionCustomResource struct{}

func TestAccLogicAppActionCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

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

func TestAccLogicAppActionCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_custom"),
		},
	})
}

func (LogicAppActionCustomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ActionID(state.ID)
	if err != nil {
		return nil, err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)
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

	raw, ok := definition["actions"]
	if !ok || raw == nil {
		return pointer.To(false), nil
	}

	actions, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type for `actions` in definition of %s: %T", workflowId, raw)
	}

	for k := range actions {
		if strings.EqualFold(k, id.Name) {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (r LogicAppActionCustomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY

}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppActionCustomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "import" {
  name         = azurerm_logic_app_action_custom.test.name
  logic_app_id = azurerm_logic_app_action_custom.test.logic_app_id
  body         = azurerm_logic_app_action_custom.test.body
}
`, r.basic(data))
}

func (LogicAppActionCustomResource) template(data acceptance.TestData) string {
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
