// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_job_schedule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationJobScheduleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Job Schedule is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"runbook_name": {
			Description:  "The name of a Runbook to link to a Schedule. It needs to be in the same Automation Account as the Schedule and Job Schedule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"schedule_name": {
			Description:  "The name of the Schedule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ScheduleName(),
		},

		"parameters": {
			Description:  "A map of key/value pairs corresponding to the arguments that can be passed to the Runbook.",
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.ParameterNames,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"run_on": {
			Description: "Name of a Hybrid Worker Group the Runbook will be executed on.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"job_schedule_id": {
			Description: "The UUID identifying the Automation Job Schedule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C this can remain as this can change if the runbook is updated but cannot be updated by the user
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"resource_manager_id": {
			Description: "The Resource Manager ID of the Automation Job Schedule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
