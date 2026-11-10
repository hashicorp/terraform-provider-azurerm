// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppActionHTTPCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewActionID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	headersRaw := d.Get("headers").(map[string]any)
	headers, err := expandLogicAppActionHttpHeaders(headersRaw)
	if err != nil {
		return err
	}

	queriesRaw := d.Get("queries").(map[string]any)
	queries, err := expandLogicAppActionHttpQueries(queriesRaw)
	if err != nil {
		return err
	}

	inputs := map[string]any{
		"method":  d.Get("method").(string),
		"uri":     d.Get("uri").(string),
		"headers": headers,
		"queries": queries,
	}

	// if it's json object then storing action's body in json object to keep consistent with azure portal
	// if starts with dynamic function (starts with "@") then store it as string
	if bodyRaw, ok := d.GetOk("body"); ok {
		if json.Valid([]byte(bodyRaw.(string))) {
			var body map[string]any
			if err := json.Unmarshal([]byte(bodyRaw.(string)), &body); err != nil {
				return fmt.Errorf("unmarshalling JSON for Action %q: %+v", id.Name, err)
			}
			inputs["body"] = body
		} else {
			inputs["body"] = bodyRaw.(string)
		}
	}

	action := map[string]any{
		"inputs": inputs,
		"type":   "http",
	}

	if v, ok := d.GetOk("run_after"); ok {
		action["runAfter"] = expandLogicAppActionRunAfter(v.(*pluginsdk.Set).List())
	}

	if err = helpers.ResourceLogicAppActionUpdate(d, meta, *workflowId, id, action, "azurerm_logic_app_action_http"); err != nil {
		return err
	}

	return resourceLogicAppActionHTTPRead(d, meta)
}
