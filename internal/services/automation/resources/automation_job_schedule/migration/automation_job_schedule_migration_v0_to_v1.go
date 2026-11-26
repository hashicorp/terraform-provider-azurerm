// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationJobScheduleV0ToV1 struct{}

func (s AutomationJobScheduleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"job_schedule_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"run_on": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"runbook_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"schedule_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (s AutomationJobScheduleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId := rawState["id"].(string)
		id, err := jobschedule.ParseJobScheduleID(oldId)
		if err != nil {
			return nil, err
		}

		scheduleName := rawState["schedule_name"].(string)
		runbookName := rawState["runbook_name"].(string)
		scheduleID := schedule.NewScheduleID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, scheduleName)
		runbookID := runbook.NewRunbookID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, runbookName)

		tfID := &commonids.CompositeResourceID[*schedule.ScheduleId, *runbook.RunbookId]{
			First:  &scheduleID,
			Second: &runbookID,
		}

		log.Printf("[DEBUG] Upgrade automation job schedule resource id from %s to %s", oldId, tfID.ID())
		rawState["id"] = tfID.ID()
		rawState["resource_manager_id"] = oldId
		return rawState, nil
	}
}
