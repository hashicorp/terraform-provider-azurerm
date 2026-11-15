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
	}
}
