package helper

import (
	"bytes"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

// parameter v should be instanced of `map[string]interface{}`
// and the hash should ignore the job_schedule_id
func resourceAutomationJobScheduleHash(v interface{}) int {
	var buf bytes.Buffer

	var scheduleName, runOn string
	var parameters map[string]*string
	if m, ok := v.(automation.JobScheduleProperties); ok {
		if m.Schedule.Name != nil {
			scheduleName = *m.Schedule.Name
		}
		if m.RunOn != nil {
			runOn = *m.RunOn
		}
		parameters = m.Parameters
	} else if m, ok := v.(map[string]interface{}); ok && m != nil {
		if v, ok := m["schedule_name"]; ok {
			scheduleName = v.(string)
		}
		if v, ok := m["run_on"]; ok {
			runOn = v.(string)
		}
		if v, ok := m["parameters"]; ok {
			parameters = utils.ExpandMapStringPtrString(v.(map[string]interface{}))
		}
	}
	buf.WriteString(fmt.Sprintf("%s-%s-%s", scheduleName, utils.FlattenMapStringPtrString(parameters), runOn))

	return pluginsdk.HashString(buf.String())
}
