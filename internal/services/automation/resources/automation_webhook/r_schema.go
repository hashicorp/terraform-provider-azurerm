// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"expiry_time": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"runbook_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"run_on_worker_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			// NOTE: O+C this is generated if omitted and cannot be updated so this should remain
			Computed:     true,
			Sensitive:    true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}
