// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_webhook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationWebhookSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Description:  "The name of the Webhook.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Webhook is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"expiry_time": {
			Description:      "Timestamp when the webhook expires.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"runbook_name": {
			Description:  "Name of the Automation Runbook to execute by Webhook.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"enabled": {
			Description: "Controls if Webhook is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"run_on_worker_group": {
			Description: "Name of the hybrid worker group the Webhook job will run on.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
		"parameters": {
			Description: "Map of input parameters passed to runbook.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"uri": {
			Description: "URI to initiate the webhook. Can be generated using [Generate URI API](https://learn.microsoft.com/rest/api/automation/webhook/generate-uri). By default, new URI is generated on each new resource creation.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			// NOTE: O+C this is generated if omitted and cannot be updated so this should remain
			Computed:     true,
			Sensitive:    true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}
