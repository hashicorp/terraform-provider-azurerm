// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_http_request

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerHttpRequestCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	schemaRaw := d.Get("schema").(string)
	var schema map[string]any
	if err := json.Unmarshal([]byte(schemaRaw), &schema); err != nil {
		return fmt.Errorf("unmarshalling JSON from Schema: %+v", err)
	}

	inputs := map[string]any{
		"schema": schema,
	}

	if v, ok := d.GetOk("method"); ok {
		inputs["method"] = v.(string)
	}

	if v, ok := d.GetOk("relative_path"); ok {
		inputs["relativePath"] = v.(string)
	}

	trigger := map[string]any{
		"inputs": inputs,
		"kind":   "Http",
		"type":   "Request",
	}

	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := workflowtriggers.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	if err := helpers.ResourceLogicAppTriggerUpdate(d, meta, *workflowId, id, trigger, "azurerm_logic_app_trigger_http_request"); err != nil {
		return err
	}

	return resourceLogicAppTriggerHttpRequestRead(d, meta)
}
