// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowrunactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppActionHttpSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the HTTP Action to be created within the Logic App Workflow.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"logic_app_id": {
			Description:  "The ID of the Logic App Workflow.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workflows.ValidateWorkflowID,
		},

		"method": {
			Description: "The HTTP Method which should be used for this HTTP Action.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				http.MethodDelete,
				http.MethodGet,
				http.MethodPatch,
				http.MethodPost,
				http.MethodPut,
			}, false),
		},

		"uri": {
			Description: "The URI which will be called when this HTTP Action is triggered.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"body": {
			Description: "The HTTP Body that should be sent to the `uri` when this HTTP Action is triggered.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				if json.Valid([]byte(oldValue)) && json.Valid([]byte(newValue)) {
					return pluginsdk.SuppressJsonDiff(k, oldValue, newValue, d)
				}
				return false
			},
		},

		"headers": {
			Description: "A Map of Key-Value Pairs that should be sent to the `uri` when this HTTP Action is triggered.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"queries": {
			Description: "A Map of Key-Value Pairs that should be sent to the `uri` when this HTTP Action is triggered.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"run_after": {
			Description: "The place of the HTTP Action in the Logic App Workflow. If not specified, the HTTP Action is right after the Trigger.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_name": {
						Description: "The name of the precedent HTTP Action.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"action_result": {
						Description: "The expected result of the precedent HTTP Action, only after which the current HTTP Action will be triggered.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
	}
}
