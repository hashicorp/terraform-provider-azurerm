package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
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
		newID := parse.NewAutomationJobScheduleID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName,
			runbookName, scheduleName)

		log.Printf("[DEBUG] Upgrade automation job schedule resource id from %s to %s", oldId, newID.ID())
		rawState["id"] = newID.ID()
		rawState["resource_manager_id"] = oldId
		return rawState, nil
	}
}
