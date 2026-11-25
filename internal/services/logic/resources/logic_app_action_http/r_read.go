// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppActionHTTPRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := parse.ActionID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)

	t, app, err := helpers.RetrieveLogicAppAction(d, meta, workflowId, id.Name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Action %q - removing from state", id.WorkflowName, id.ResourceGroup, id.Name)
		d.SetId("")
		return nil
	}

	action := *t

	d.Set("name", id.Name)
	d.Set("logic_app_id", app.Id)

	actionType := action["type"].(string)
	if !strings.EqualFold(actionType, "http") {
		return fmt.Errorf("expected an HTTP Action for Action %s - got %q", id, actionType)
	}

	v := action["inputs"]
	if v != nil {
		inputs, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("parsing `inputs` for HTTP Action %s", id)
		}

		if uri := inputs["uri"]; uri != nil {
			d.Set("uri", uri.(string))
		}

		if method := inputs["method"]; method != nil {
			d.Set("method", method.(string))
		}

		if body := inputs["body"]; body != nil {
			switch body.(type) {
			case map[string]any:
				// if user edit workflow in portal, the body becomes json object
				v, err := json.Marshal(body)
				if err != nil {
					return fmt.Errorf("serializing `body` for Action %q: %+v", id.Name, err)
				}
				d.Set("body", string(v))
			case string:
				d.Set("body", body)
			}
		}

		if headers := inputs["headers"]; headers != nil {
			if err := d.Set("headers", headers.(map[string]any)); err != nil {
				return fmt.Errorf("setting `headers` for HTTP Action %q: %+v", id.Name, err)
			}
		}

		if queries := inputs["queries"]; queries != nil {
			if err := d.Set("queries", queries.(map[string]any)); err != nil {
				return fmt.Errorf("setting `queries` for HTTP Action %q: %+v", id.Name, err)
			}
		}
	}

	v = action["runAfter"]
	if v != nil {
		runAfter, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("parsing `runAfter` for HTTP Action %s", id)
		}
		if err := d.Set("run_after", flattenLogicAppActionRunAfter(runAfter)); err != nil {
			return fmt.Errorf("setting `runAfter` for HTTP Action %q: %+v", id.Name, err)
		}
	}

	return nil
}
