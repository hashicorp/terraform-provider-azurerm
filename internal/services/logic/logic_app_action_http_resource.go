// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowrunactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceLogicAppActionHTTP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppActionHTTPCreateUpdate,
		Read:   resourceLogicAppActionHTTPRead,
		Update: resourceLogicAppActionHTTPCreateUpdate,
		Delete: resourceLogicAppActionHTTPDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ActionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"logic_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workflows.ValidateWorkflowID,
			},

			"method": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodDelete,
					http.MethodGet,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
				}, false),
			},

			"uri": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"body": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if json.Valid([]byte(oldValue)) && json.Valid([]byte(newValue)) {
						return pluginsdk.SuppressJsonDiff(k, oldValue, newValue, d)
					}
					return false
				},
			},

			"headers": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"queries": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"run_after": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"action_result": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(workflowrunactions.WorkflowStatusSucceeded),
								string(workflowrunactions.WorkflowStatusFailed),
								string(workflowrunactions.WorkflowStatusSkipped),
								string(workflowrunactions.WorkflowStatusTimedOut),
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceLogicAppActionHTTPCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewActionID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	headersRaw := d.Get("headers").(map[string]interface{})
	headers, err := expandLogicAppActionHttpHeaders(headersRaw)
	if err != nil {
		return err
	}

	queriesRaw := d.Get("queries").(map[string]interface{})
	queries, err := expandLogicAppActionHttpQueries(queriesRaw)
	if err != nil {
		return err
	}

	inputs := map[string]interface{}{
		"method":  d.Get("method").(string),
		"uri":     d.Get("uri").(string),
		"headers": headers,
		"queries": queries,
	}

	// if it's json object then storing action's body in json object to keep consistent with azure portal
	// if starts with dynamic function (starts with "@") then store it as string
	if bodyRaw, ok := d.GetOk("body"); ok {
		if json.Valid([]byte(bodyRaw.(string))) {
			var body map[string]interface{}
			if err := json.Unmarshal([]byte(bodyRaw.(string)), &body); err != nil {
				return fmt.Errorf("unmarshalling JSON for Action %q: %+v", id.Name, err)
			}
			inputs["body"] = body
		} else {
			inputs["body"] = bodyRaw.(string)
		}
	}

	action := map[string]interface{}{
		"inputs": inputs,
		"type":   "http",
	}

	if v, ok := d.GetOk("run_after"); ok {
		action["runAfter"] = expandLogicAppActionRunAfter(v.(*pluginsdk.Set).List())
	}

	err = resourceLogicAppActionUpdate(d, meta, *workflowId, id, action, "azurerm_logic_app_action_http")
	if err != nil {
		return err
	}

	return resourceLogicAppActionHTTPRead(d, meta)
}

func resourceLogicAppActionHTTPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ActionID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)

	t, app, err := retrieveLogicAppAction(d, meta, workflowId, id.Name)
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
	if v == nil {
		return fmt.Errorf("`inputs` was nil for HTTP Action %s", id)
	}

	inputs, ok := v.(map[string]interface{})
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
		case map[string]interface{}:
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
		hv := headers.(map[string]interface{})
		if err := d.Set("headers", hv); err != nil {
			return fmt.Errorf("setting `headers` for HTTP Action %q: %+v", id.Name, err)
		}
	}

	if queries := inputs["queries"]; queries != nil {
		qv := queries.(map[string]interface{})
		if err := d.Set("queries", qv); err != nil {
			return fmt.Errorf("setting `queries` for HTTP Action %q: %+v", id.Name, err)
		}
	}

	v = action["runAfter"]
	if v != nil {
		runAfter, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("parsing `runAfter` for HTTP Action %s", id)
		}
		if err := d.Set("run_after", flattenLogicAppActionRunAfter(runAfter)); err != nil {
			return fmt.Errorf("setting `runAfter` for HTTP Action %q: %+v", id.Name, err)
		}
	}

	return nil
}

func resourceLogicAppActionHTTPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ActionID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)

	err = resourceLogicAppActionRemove(d, meta, workflowId, id.Name)
	if err != nil {
		return fmt.Errorf("removing Action %s: %+v", id, err)
	}

	return nil
}

func expandLogicAppActionHttpHeaders(headersRaw map[string]interface{}) (*map[string]string, error) {
	headers := make(map[string]string)

	for i, v := range headersRaw {
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}

		headers[i] = value
	}

	return &headers, nil
}

func expandLogicAppActionHttpQueries(queriesRaw map[string]interface{}) (*map[string]string, error) {
	queries := make(map[string]string)

	for i, v := range queriesRaw {
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}

		queries[i] = value
	}

	return &queries, nil
}
