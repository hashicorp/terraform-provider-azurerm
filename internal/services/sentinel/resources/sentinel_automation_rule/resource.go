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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: automationrules.ValidateWorkspaceID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"order": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"triggers_on": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersOnIncidents),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersOn(), false),
		},

		"triggers_when": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersWhenCreated),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersWhen(), false),
		},

		"expiration": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"condition_json": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// We can't use the pluginsdk.SuppressJsonDiff here as the "condition_json" is always an array, while that function assume its input is an object.
			// Once https://github.com/hashicorp/terraform-plugin-sdk/pull/1102 is merged, we can switch to pluginsdk.SuppressJsonDiff.
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				return pluginsdk.NormalizeJson(old) == pluginsdk.NormalizeJson(new)
			},
			ValidateFunc: validation.StringIsJSON,
		},

		"action_incident_task": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"title": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"action_incident": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"status": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentStatus(), false),
					},

					"classification": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(automationrules.IncidentClassificationUndetermined),
							string(automationrules.IncidentClassificationBenignPositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousButExpected),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonIncorrectAlertLogic),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonInaccurateData),
							string(automationrules.IncidentClassificationTruePositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousActivity),
						}, false),
					},

					"classification_comment": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"owner_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"severity": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentSeverity(), false),
					},
				},
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},

		"action_playbook": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: sentinelAutomationRuleSchema(),
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},
	}
	return &pluginsdk.Resource{
		Create: resourceSentinelAutomationRuleCreateOrUpdate,
		Read:   resourceSentinelAutomationRuleRead,
		Update: resourceSentinelAutomationRuleCreateOrUpdate,
		Delete: resourceSentinelAutomationRuleDelete,

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
