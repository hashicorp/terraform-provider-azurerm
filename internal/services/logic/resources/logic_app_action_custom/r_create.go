// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_custom

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppActionCustomCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewActionID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	bodyRaw := d.Get("body").(string)

	var body map[string]any
	if err := json.Unmarshal([]byte(bodyRaw), &body); err != nil {
		return fmt.Errorf("unmarshalling JSON for Custom Action %q: %+v", id.Name, err)
	}

	if err := helpers.ResourceLogicAppActionUpdate(d, meta, *workflowId, id, body, "azurerm_logic_app_action_custom"); err != nil {
		return err
	}

	return resourceLogicAppActionCustomRead(d, meta)
}
