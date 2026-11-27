// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sentinelAutomationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"order": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"logic_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workflows.ValidateWorkflowID,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C We'll use the current tenant id if this property is absent.
			Computed:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}
