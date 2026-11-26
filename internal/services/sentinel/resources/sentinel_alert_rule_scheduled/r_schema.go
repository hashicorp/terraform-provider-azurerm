// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_scheduled

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sentinelAlertRuleScheduledSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"alert_rule_template_guid": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"alert_rule_template_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"event_grouping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aggregation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEventGroupingAggregationKind(), false),
					},
				},
			},
		},

		"tactics": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAttackTactic(), false),
			},
		},

		"techniques": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"incident": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C The API creates an incident if omitted but overwriting this/reverting to the default can be done without issue so this can remain
			Computed: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"create_incident_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"grouping": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
								"lookback_duration": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT5M",
									ValidateFunc: validation.ISO8601Duration,
								},
								"reopen_closed_incidents": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"entity_matching_method": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      alertrules.MatchingMethodAnyAlert,
									ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForMatchingMethod(), false),
								},
								"by_entities": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEntityMappingType(), false),
									},
								},
								"by_alert_details": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertDetail(), false),
									},
								},
								"by_custom_details": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"severity": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertSeverity(), false),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"query": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"query_frequency": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT5H",
			ValidateFunc: validation.ISO8601DurationBetween("PT5M", "P14D"),
		},

		"query_period": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT5H",
			ValidateFunc: validation.ISO8601DurationBetween("PT5M", "P14D"),
		},

		"trigger_operator": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(alertrules.TriggerOperatorGreaterThan),
			ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForTriggerOperator(), false),
		},

		"trigger_threshold": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"suppression_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"suppression_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT5H",
			ValidateFunc: validation.ISO8601DurationBetween("PT5M", "PT24H"),
		},
		"alert_details_override": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description_format": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"display_name_format": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"severity_column_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tactics_column_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dynamic_property": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertProperty(), false),
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
		},
		"custom_details": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"entity_mapping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 10,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"entity_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEntityMappingType(), false),
					},
					"field_mapping": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 3,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"identifier": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"column_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
		"sentinel_entity_mapping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 10,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"column_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
