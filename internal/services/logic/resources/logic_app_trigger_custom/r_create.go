// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_custom

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerCustomCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := workflowtriggers.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	bodyRaw := d.Get("body").(string)

	var body map[string]any
	if err := json.Unmarshal([]byte(bodyRaw), &body); err != nil {
		return fmt.Errorf("unmarshalling JSON for %s: %+v", id.ID(), err)
	}

	log.Printf("[DEBUG] logic_custom_trigger initial body is: %s", body)

	if err := helpers.ResourceLogicAppTriggerUpdate(d, meta, *workflowId, id, body, "azurerm_logic_app_trigger_custom"); err != nil {
		return err
	}

	return resourceLogicAppTriggerCustomRead(d, meta)
}
