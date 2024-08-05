// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ExpandAutomationJobSchedule(input []interface{}, runBookName string) (*map[string]jobschedule.JobScheduleCreateParameters, error) {
	res := make(map[string]jobschedule.JobScheduleCreateParameters)
	if len(input) == 0 || input[0] == nil {
		return &res, nil
	}

	for _, v := range input {
		js := v.(map[string]interface{})
		// skip SDK v2 bug: https://github.com/hashicorp/terraform-plugin-sdk/issues/1248
		if js["schedule_name"] == "" {
			continue
		}
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
		res[ResourceAutomationJobScheduleDigest(jobScheduleCreateParameters.Properties)] = jobScheduleCreateParameters
	}

	return &res, nil
}

func FlattenAutomationJobSchedule(jsMap map[uuid.UUID]jobschedule.JobScheduleProperties) *pluginsdk.Set {
	res := &pluginsdk.Set{
		F: ResourceAutomationJobScheduleHash,
	}
	for jsId, js := range jsMap {
		var scheduleName, runOn string
		if js.Schedule != nil && js.Schedule.Name != nil {
			scheduleName = *js.Schedule.Name
		}

		if js.RunOn != nil {
			runOn = *js.RunOn
		}

		// for API casing issue: https://github.com/Azure/azure-sdk-for-go/issues/4780
		parameters := map[string]string{}
		if js.Parameters != nil {
			for key, value := range *js.Parameters {
				parameters[strings.ToLower(key)] = value
			}
		}

		res.Add(map[string]interface{}{
			"schedule_name":   scheduleName,
			"parameters":      parameters,
			"run_on":          runOn,
			"job_schedule_id": jsId.String(),
		})
	}

	return res
}

func ResourceAutomationJobScheduleDigest(v interface{}) string {
	var buf bytes.Buffer
	var paramString map[string]string
	var scheduleName, runOn string
	switch job := v.(type) {
	case map[string]interface{}:
		scheduleName = job["schedule_name"].(string)
		runOn = job["run_on"].(string)
		switch param := job["parameters"].(type) {
		case map[string]string:
			paramString = param
		case map[string]interface{}:
			paramString = map[string]string{}
			for k, v := range param {
				paramString[k] = fmt.Sprintf("%v", v)
			}
		}
	case jobschedule.JobScheduleCreateProperties:
		scheduleName = pointer.From(job.Schedule.Name)
		runOn = pointer.From(job.Runbook.Name)
		paramString = pointer.From(job.Parameters)
	case *jobschedule.JobScheduleProperties:
		scheduleName = pointer.From(pointer.From(job.Schedule).Name)
		runOn = pointer.From(pointer.From(job.Runbook).Name)
		paramString = pointer.From(job.Parameters)
	}
	buf.WriteString(fmt.Sprintf("%s-%s-", scheduleName, runOn))

	var keys []string
	for k := range paramString {
		// params key will be returned as title-cased even created with lower-case
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s:%v;", strings.ToLower(k), paramString[k]))
	}
	return buf.String()
}

func ResourceAutomationJobScheduleHash(v interface{}) int {
	return pluginsdk.HashString(ResourceAutomationJobScheduleDigest(v))
}
