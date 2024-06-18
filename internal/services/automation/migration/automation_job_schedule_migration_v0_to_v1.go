// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/schedule"
	pluginsdk "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationJobScheduleV0ToV1 struct{}

func (s AutomationJobScheduleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"job_schedule_id": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"parameters": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"run_on": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"runbook_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"schedule_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (s AutomationJobScheduleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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
