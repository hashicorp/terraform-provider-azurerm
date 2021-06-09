package helper

import (
	"bytes"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/gofrs/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func JobScheduleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeSet,
		Optional:   true,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"schedule_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.ScheduleName(),
				},

				"parameters": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					ValidateFunc: validate.ParameterNames,
				},

				"run_on": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"job_schedule_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
		Set: resourceAutomationJobScheduleHash,
	}
}

func ExpandAutomationJobSchedule(input []interface{}, runBookName string) (*map[uuid.UUID]automation.JobScheduleCreateParameters, error) {
	res := make(map[uuid.UUID]automation.JobScheduleCreateParameters)
	if len(input) == 0 || input[0] == nil {
		return &res, nil
	}

	for _, v := range input {
		js := v.(map[string]interface{})
		jobScheduleCreateParameters := automation.JobScheduleCreateParameters{
			JobScheduleCreateProperties: &automation.JobScheduleCreateProperties{
				Schedule: &automation.ScheduleAssociationProperty{
					Name: utils.String(js["schedule_name"].(string)),
				},
				Runbook: &automation.RunbookAssociationProperty{
					Name: utils.String(runBookName),
				},
			},
		}

		if v, ok := js["parameters"]; ok {
			jsParameters := make(map[string]*string)
			for k, v := range v.(map[string]interface{}) {
				value := v.(string)
				jsParameters[k] = &value
			}
			jobScheduleCreateParameters.JobScheduleCreateProperties.Parameters = jsParameters
		}

		if v, ok := js["run_on"]; ok && v.(string) != "" {
			value := v.(string)
			jobScheduleCreateParameters.JobScheduleCreateProperties.RunOn = &value
		}
		jobScheduleUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		res[jobScheduleUUID] = jobScheduleCreateParameters
	}

	return &res, nil
}

func FlattenAutomationJobSchedule(jsMap map[uuid.UUID]automation.JobScheduleProperties) *pluginsdk.Set {
	res := &pluginsdk.Set{
		F: resourceAutomationJobScheduleHash,
	}
	for jsId, js := range jsMap {
		var scheduleName, runOn string
		if js.Schedule.Name != nil {
			scheduleName = *js.Schedule.Name
		}

		if js.RunOn != nil {
			runOn = *js.RunOn
		}

		res.Add(map[string]interface{}{
			"schedule_name":   scheduleName,
			"parameters":      utils.FlattenMapStringPtrString(js.Parameters),
			"run_on":          runOn,
			"job_schedule_id": jsId.String(),
		})
	}

	return res
}

func resourceAutomationJobScheduleHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(automation.JobScheduleProperties); ok {
		var scheduleName, runOn string
		if m.Schedule.Name != nil {
			scheduleName = *m.Schedule.Name
		}

		if m.RunOn != nil {
			runOn = *m.RunOn
		}

		buf.WriteString(fmt.Sprintf("%s-%s-%s-%s", scheduleName, utils.FlattenMapStringPtrString(m.Parameters), runOn, *m.JobScheduleID))
	}

	return pluginsdk.HashString(buf.String())
}
