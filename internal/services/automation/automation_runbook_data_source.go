// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationRunbookDataSource struct{}

type AutomationRunbookDataSourceModel struct {
	RunbookName            string                                        `tfschema:"name"`
	AutomationAccountName  string                                        `tfschema:"automation_account_name"`
	ResourceGroupName      string                                        `tfschema:"resource_group_name"`
	Location               string                                        `tfschema:"location"`
	Description            string                                        `tfschema:"description"`
	LogProgress            bool                                          `tfschema:"log_progress"`
	LogVerbose             bool                                          `tfschema:"log_verbose"`
	RunbookType            string                                        `tfschema:"runbook_type"`
	LogActivityTrace       int64                                         `tfschema:"log_activity_trace_level"`
	Content                string                                        `tfschema:"content"`
	RuntimeEnvironmentName string                                        `tfschema:"runtime_environment_name"`
	Tags                   map[string]string                             `tfschema:"tags"`
	JobSchedule            []AutomationRunbookJobScheduleDataSourceModel `tfschema:"job_schedule"`
}

type AutomationRunbookJobScheduleDataSourceModel struct {
	ScheduleName  string            `tfschema:"schedule_name"`
	Parameters    map[string]string `tfschema:"parameters"`
	RunOn         string            `tfschema:"run_on"`
	JobScheduleId string            `tfschema:"job_schedule_id"`
}

var _ sdk.DataSource = AutomationRunbookDataSource{}

func (d AutomationRunbookDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d AutomationRunbookDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"log_progress": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"log_verbose": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"runbook_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"log_activity_trace_level": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"job_schedule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"schedule_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"parameters": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"run_on": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"job_schedule_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"publish_content_link": contentLinkSchemaDataSource(),

		"draft": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"creation_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"content_link": contentLinkSchemaDataSource(),

					"edit_mode_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"last_modified_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"output_types": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"parameters": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true},

								"mandatory": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"position": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"default_value": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"runtime_environment_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d AutomationRunbookDataSource) ModelObject() interface{} {
	return &AutomationRunbookDataSourceModel{}
}

func (d AutomationRunbookDataSource) ResourceType() string {
	return "azurerm_automation_runbook"
}

func (d AutomationRunbookDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.Runbook
			jsClient := metadata.Client.Automation.JobSchedule

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AutomationRunbookDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := runbook.NewRunbookID(subscriptionId, state.ResourceGroupName, state.AutomationAccountName, state.RunbookName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			contentResp, err := client.GetContent(ctx, id)
			if err != nil {
				if !response.WasNotFound(contentResp.HttpResponse) {
					return fmt.Errorf("retrieving runbook content %s: %+v", id, err)
				}
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.Normalize(model.Location)

			if model.Properties.Description != nil {
				state.Description = pointer.From(model.Properties.Description)
			}

			if model.Properties.LogProgress != nil {
				state.LogProgress = pointer.From(model.Properties.LogProgress)
			}

			if model.Properties.RunbookType != nil {
				state.RunbookType = string(pointer.From(model.Properties.RunbookType))
			}

			if model.Properties.LogVerbose != nil {
				state.LogVerbose = pointer.From(model.Properties.LogVerbose)
			}

			if model.Properties.LogActivityTrace != nil {
				state.LogActivityTrace = pointer.From(model.Properties.LogActivityTrace)
			}

			if model.Properties.RuntimeEnvironment != nil {
				state.RuntimeEnvironmentName = pointer.From(model.Properties.RuntimeEnvironment)
			}

			jsMap := make(map[uuid.UUID]jobschedule.JobScheduleProperties)
			automationAccountId := jobschedule.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)

			filter := fmt.Sprintf("properties/runbook/name eq '%s'", id.RunbookName)
			jsIterator, err := jsClient.ListByAutomationAccount(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
			if err != nil {
				return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
			}
			for _, item := range pointer.From(jsIterator.Model) {
				if itemProps := item.Properties; itemProps != nil {
					if itemProps.JobScheduleId == nil || *itemProps.JobScheduleId == "" {
						return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
					}
					jsId, err := uuid.FromString(*itemProps.JobScheduleId)
					if err != nil {
						return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
					}
					// get job schedule from GET API, `ListByAutomationAccountComplete` lost parameters
					jobscheduleID, err := jobschedule.ParseJobScheduleID(pointer.From(item.Id))
					if err != nil {
						return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
					}
					jsResult, err := jsClient.Get(ctx, *jobscheduleID)
					if err != nil {
						return fmt.Errorf("retrieving job schedule by %s: %v", *jobscheduleID, err)
					}
					if jsResult.Model != nil && jsResult.Model.Properties != nil {
						jsMap[jsId] = *jsResult.Model.Properties
					}
				}
			}

			state.JobSchedule = FlattenAutomationRunbookJobScheduleDataSource(jsMap)

			if contentResp.Model != nil {
				state.Content = string(pointer.From(contentResp.Model))
			}

			state.Tags = pointer.From(model.Tags)

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func contentLinkSchemaDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"uri": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"hash": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"algorithm": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"value": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func FlattenAutomationRunbookJobScheduleDataSource(jsMap map[uuid.UUID]jobschedule.JobScheduleProperties) []AutomationRunbookJobScheduleDataSourceModel {
	results := make([]AutomationRunbookJobScheduleDataSourceModel, 0)

	for jsId, js := range jsMap {
		var scheduleName string
		if js.Schedule != nil {
			scheduleName = pointer.From(js.Schedule.Name)
		}

		parameters := map[string]string{}
		if js.Parameters != nil {
			for key, value := range *js.Parameters {
				parameters[strings.ToLower(key)] = value
			}
		}

		results = append(results, AutomationRunbookJobScheduleDataSourceModel{
			ScheduleName:  scheduleName,
			Parameters:    parameters,
			RunOn:         pointer.From(js.RunOn),
			JobScheduleId: jsId.String(),
		})
	}

	return results
}
