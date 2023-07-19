// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"bytes"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/jobschedule"
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

func ExpandAutomationJobSchedule(input []interface{}, runBookName string) (*map[uuid.UUID]jobschedule.JobScheduleCreateParameters, error) {
	res := make(map[uuid.UUID]jobschedule.JobScheduleCreateParameters)
	if len(input) == 0 || input[0] == nil {
		return &res, nil
	}

	for _, v := range input {
		js := v.(map[string]interface{})
		jobScheduleCreateParameters := jobschedule.JobScheduleCreateParameters{
			Properties: jobschedule.JobScheduleCreateProperties{
				Schedule: jobschedule.ScheduleAssociationProperty{
					Name: utils.String(js["schedule_name"].(string)),
				},
				Runbook: jobschedule.RunbookAssociationProperty{
					Name: utils.String(runBookName),
				},
			},
		}

		if v, ok := js["parameters"]; ok {
			jsParameters := make(map[string]string)
			for k, v := range v.(map[string]interface{}) {
				value := v.(string)
				jsParameters[k] = value
			}
			jobScheduleCreateParameters.Properties.Parameters = &jsParameters
		}

		if v, ok := js["run_on"]; ok && v.(string) != "" {
			value := v.(string)
			jobScheduleCreateParameters.Properties.RunOn = &value
		}
		jobScheduleUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		res[jobScheduleUUID] = jobScheduleCreateParameters
	}

	return &res, nil
}

func FlattenAutomationJobSchedule(jsMap map[uuid.UUID]jobschedule.JobScheduleProperties) *pluginsdk.Set {
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
			"parameters":      js.Parameters,
			"run_on":          runOn,
			"job_schedule_id": jsId.String(),
		})
	}

	return res
}

func resourceAutomationJobScheduleHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(jobschedule.JobScheduleProperties); ok {
		var scheduleName, runOn string
		if m.Schedule.Name != nil {
			scheduleName = *m.Schedule.Name
		}

		if m.RunOn != nil {
			runOn = *m.RunOn
		}

		buf.WriteString(fmt.Sprintf("%s-%s-%s-%s", scheduleName, m.Parameters, runOn, *m.JobScheduleId))
	}

	return pluginsdk.HashString(buf.String())
}
