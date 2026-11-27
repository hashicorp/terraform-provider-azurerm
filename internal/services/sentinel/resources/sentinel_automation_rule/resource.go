// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_automation_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RegisterResource() *pluginsdk.Resource {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The UUID which should be used for this Sentinel Automation Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace where this Sentinel applies to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: automationrules.ValidateWorkspaceID,
		},

		"display_name": {
			Description:  "The display name which should be used for this Sentinel Automation Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"order": {
			Description:  "The order of this Sentinel Automation Rule.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"enabled": {
			Description: "Whether this Sentinel Automation Rule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"triggers_on": {
			Description:  "What triggers this automation rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersOnIncidents),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersOn(), false),
		},

		"triggers_when": {
			Description:  "When this automation rule will be triggered.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersWhenCreated),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersWhen(), false),
		},

		"expiration": {
			Description:      "The time in RFC3339 format of kind `UTC` that determines when this Automation Rule should expire and be disabled.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.IsRFC3339Time,
			DiffSuppressFunc: suppress.RFC3339Time,
		},

		"condition_json": {
			Description:  "A JSON array of one or more condition JSON objects as is defined [here](https://learn.microsoft.com/rest/api/securityinsights/preview/automation-rules/create-or-update?tabs=HTTP#automationruletriggeringlogic).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsJSON,
			// We can't use the pluginsdk.SuppressJsonDiff here as the "condition_json" is always an array, while that function assume its input is an object.
			// Once https://github.com/hashicorp/terraform-plugin-sdk/pull/1102 is merged, we can switch to pluginsdk.SuppressJsonDiff.
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				return pluginsdk.NormalizeJson(old) == pluginsdk.NormalizeJson(new)
			},
		},

		"action_incident_task": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Description:  "The execution order of this action.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"title": {
						Description:  "The title of the incident task.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Description:  "The description of the incident task.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"action_incident": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Description:  "The execution order of this action.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"status": {
						Description:  "The status to set to the incident.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentStatus(), false),
					},

					"classification": {
						Description: "The classification of the incident, when closing it.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(automationrules.IncidentClassificationUndetermined),
							string(automationrules.IncidentClassificationBenignPositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousButExpected),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonIncorrectAlertLogic),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonInaccurateData),
							string(automationrules.IncidentClassificationTruePositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousActivity),
						}, false),
					},

					"classification_comment": {
						Description:  "The comment why the incident is to be closed.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Description: "A list of labels to add to the incident.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"owner_id": {
						Description:  "The object ID of the entity this incident is assigned to.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"severity": {
						Description:  "The severity to add to the incident.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentSeverity(), false),
					},
				},
			},
		},

		"action_playbook": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
			Elem: &pluginsdk.Resource{
				Schema: sentinelAutomationRuleSchema(),
			},
		},
	}
	return &pluginsdk.Resource{
		Description: "Manages a Sentinel Automation Rule.",
		Create:      resourceSentinelAutomationRuleCreateOrUpdate,
		Read:        resourceSentinelAutomationRuleRead,
		Update:      resourceSentinelAutomationRuleCreateOrUpdate,
		Delete:      resourceSentinelAutomationRuleDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SentinelAutomationRuleV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := automationrules.ParseAutomationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: schema,
	}
}
