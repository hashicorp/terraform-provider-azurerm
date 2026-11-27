// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationRunbookSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Runbook.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Runbook is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"runbook_type": {
			Description:  "The type of the runbook.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(runbook.PossibleValuesForRunbookTypeEnum(), false),
		},

		"log_progress": {
			Description: "Progress log option.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"log_verbose": {
			Description: "Verbose log option.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"description": {
			Description: "A description for the runbook.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"content": {
			Description: "The desired content of the runbook.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C the API returns some defaults for this if `publish_content_link` is specified
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"content", "publish_content_link", "draft"},
		},

		"job_schedule": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Set:        helpers.ResourceAutomationJobScheduleHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"schedule_name": {
						Description:  "The name of the Schedule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ScheduleName(),
					},

					"parameters": {
						Description:  "A map of key/value pairs corresponding to the arguments that can be passed to the Runbook.",
						Type:         pluginsdk.TypeMap,
						Optional:     true,
						ValidateFunc: validate.ParameterNames,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"run_on": {
						Description: "Name of a Hybrid Worker Group the Runbook will be executed on.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"job_schedule_id": {
						Description: "The UUID of Automation Runbook job schedule ID.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"publish_content_link": contentLinkSchema(false),

		"draft": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"creation_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"content_link": contentLinkSchema(true),

					"edit_mode_enabled": {
						Description: "Whether the draft in edit mode.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"last_modified_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"output_types": {
						Description: "The output types of the runbook.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"parameters": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Description:  "The name of the parameter.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"type": {
									Description:  "The type of this parameter.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mandatory": {
									Description: "Whether this parameter is mandatory.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},

								"position": {
									Description:  "The position of the parameter.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"default_value": {
									Description:  "The default value of the parameter.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"log_activity_trace_level": {
			Description:  "The activity-level tracing options of the runbook, available only for Graphical runbooks. Must turn on Verbose logging to see the tracing.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"runtime_environment_name": {
			Description:  "The runtime environment name for the runbook.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func contentLinkSchema(isDraft bool) *pluginsdk.Schema {
	ins := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"uri": {
					Description: "The URI of the runbook content.",
					Type:        pluginsdk.TypeString,
					Required:    true,
					ValidateFunc: validation.Any(
						validation.IsURLWithScheme([]string{"http", "https"}),
						validation.StringIsEmpty,
					),
				},

				"version": {
					Description:  "The version of the content.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"hash": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"algorithm": {
								Description:  "The hash algorithm used to hash the content.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"value": {
								Description:  "The expected hash value of the content.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
	if !isDraft {
		ins.AtLeastOneOf = []string{"content", "publish_content_link", "draft"}
	}
	return ins
}
