// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_http_request

import (
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppTriggerHttpRequestSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the HTTP Request Trigger to be created within the Logic App Workflow.",
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

		"schema": {
			Description:      "A JSON Blob defining the Schema of the incoming request. This needs to be valid JSON.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"method": {
			Description: "The HTTP Method which the request will be using.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				http.MethodDelete,
				http.MethodGet,
				http.MethodPatch,
				http.MethodPost,
				http.MethodPut,
			}, false),
		},

		"relative_path": {
			Description:  "The Relative Path used for this Request.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.TriggerHttpRequestRelativePath,
		},

		"callback_url": {
			Description: "The URL of the Trigger within the Logic App Workflow. For use with certain resources like [monitor_action_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_action_group) and [security_center_automation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/security_center_automation).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
