// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationRunbookSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"runbook_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(runbook.PossibleValuesForRunbookTypeEnum(), false),
		},

		"log_progress": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"log_verbose": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C the API returns some defaults for this if `publish_content_link` is specified
			Computed:     true,
			AtLeastOneOf: []string{"content", "publish_content_link", "draft"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"job_schedule": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
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
			Set: helper.ResourceAutomationJobScheduleHash,
		},

		"publish_content_link": contentLinkSchema(false),

		"draft": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"creation_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"content_link": contentLinkSchema(true),

					"edit_mode_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"last_modified_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"output_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mandatory": {
									Type:     pluginsdk.TypeBool,
									Default:  false,
									Optional: true,
								},

								"position": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"default_value": {
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
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"runtime_environment_name": {
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
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.Any(
						validation.IsURLWithScheme([]string{"http", "https"}),
						validation.StringIsEmpty,
					),
				},

				"version": {
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
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"value": {
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
