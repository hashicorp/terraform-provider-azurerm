package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func JobScheduleSchema() *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeSet,
		Optional:   true,
		Computed:   true,
		ConfigMode: schema.SchemaConfigModeAttr,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"schedule_name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.ScheduleName(),
				},

				"parameters": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					ValidateFunc: func(v interface{}, _ string) (warnings []string, errors []error) {
						m := v.(map[string]interface{})
						for k := range m {
							if k != strings.ToLower(k) {
								errors = append(errors, fmt.Errorf("Due to a bug in the implementation of Runbooks in Azure, the parameter names need to be specified in lowercase only. See: \"https://github.com/Azure/azure-sdk-for-go/issues/4780\" for more information."))
							}
						}

						return warnings, errors
					},
				},

				"run_on": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"job_schedule_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
		Set: resourceAutomationJobScheduleHash,
	}
}

func ExpandAutomationJobSchedule(input []interface{}, runBookName string) map[uuid.UUID]automation.JobScheduleCreateParameters {
	res := make(map[uuid.UUID]automation.JobScheduleCreateParameters)
	if len(input) == 0 || input[0] == nil {
		return res
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
		jobScheduleUUID := uuid.NewV4()
		res[jobScheduleUUID] = jobScheduleCreateParameters
	}

	return res
}

func FlattenAutomationJobSchedule(jsMap map[uuid.UUID]automation.JobScheduleProperties) *schema.Set {
	res := &schema.Set{
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

	return hashcode.String(buf.String())
}
