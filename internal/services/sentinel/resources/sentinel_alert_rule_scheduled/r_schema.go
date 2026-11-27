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
			Description:  "The name which should be used for this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace this Sentinel Scheduled Alert Rule belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"display_name": {
			Description:  "The friendly name of this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"alert_rule_template_guid": {
			Description:  "The GUID of the alert rule template which is used for this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"alert_rule_template_version": {
			Description: "The version of the alert rule template which is used for this Sentinel Scheduled Alert Rule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"description": {
			Description:  "The description of this Sentinel Scheduled Alert Rule.",
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
						Description:  "The aggregation type of grouping the events.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEventGroupingAggregationKind(), false),
					},
				},
			},
		},

		"tactics": {
			Description: "A list of categories of attacks by which to classify the rule.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAttackTactic(), false),
			},
		},

		"techniques": {
			Description: "A list of techniques of attacks by which to classify the rule.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
						Description: "Whether to create an incident from alerts triggered by this Sentinel Scheduled Alert Rule.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"grouping": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Description: "Enable grouping incidents created from alerts triggered by this Sentinel Scheduled Alert Rule.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},
								"lookback_duration": {
									Description:  "Limit the group to alerts created within the lookback duration (in ISO 8601 duration format).",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT5M",
									ValidateFunc: validation.ISO8601Duration,
								},
								"reopen_closed_incidents": {
									Description: "Whether to re-open closed matching incidents.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
								"entity_matching_method": {
									Description:  "The method used to group incidents.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      alertrules.MatchingMethodAnyAlert,
									ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForMatchingMethod(), false),
								},
								"by_entities": {
									Description: "A list of entity types to group by, only when the `entity_matching_method` is `Selected`.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEntityMappingType(), false),
									},
								},
								"by_alert_details": {
									Description: "A list of alert details to group by, only when the `entity_matching_method` is `Selected`.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertDetail(), false),
									},
								},
								"by_custom_details": {
									Description: "A list of custom details keys to group by, only when the `entity_matching_method` is `Selected`. Only keys defined in the `custom_details` may be used.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
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
			Description:  "The alert severity of this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertSeverity(), false),
		},

		"enabled": {
			Description: "Whether the Sentinel Scheduled Alert Rule should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"query": {
			Description:  "The query of this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"query_frequency": {
			Description:  "The ISO 8601 timespan duration between two consecutive queries.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT5H",
			ValidateFunc: validation.ISO8601DurationBetween("PT5M", "P14D"),
		},

		"query_period": {
			Description:  "The ISO 8601 timespan duration, which determines the time period of the data covered by the query. For example, it can query the past 10 minutes of data, or the past 6 hours of data.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT5H",
			ValidateFunc: validation.ISO8601DurationBetween("PT5M", "P14D"),
		},

		"trigger_operator": {
			Description:  "The alert trigger operator, combined with `trigger_threshold`, setting alert threshold of this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(alertrules.TriggerOperatorGreaterThan),
			ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForTriggerOperator(), false),
		},

		"trigger_threshold": {
			Description:  "The baseline number of query results generated, combined with `trigger_operator`, setting alert threshold of this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"suppression_enabled": {
			Description: "Whether the Sentinel NRT Alert Rules should stop running query after an alert is generated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"suppression_duration": {
			Description:  "If `suppression_enabled` is `true`, this is ISO 8601 timespan duration, which specifies the amount of time the query should stop running after alert is generated.",
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
						Description:  "The format containing columns name(s) to override the description of this Sentinel Alert Rule.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"display_name_format": {
						Description:  "The format containing columns name(s) to override the name of this Sentinel Alert Rule.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"severity_column_name": {
						Description:  "The column name to take the alert severity from.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tactics_column_name": {
						Description:  "The column name to take the alert tactics from.",
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
									Description:  "The name of the dynamic property.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertProperty(), false),
								},
								"value": {
									Description:  "The value of the dynamic property.",
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
			Description: "A map of string key-value pairs of columns to be attached to this Sentinel Scheduled Alert Rule. The key will appear as the field name in alerts and the value is the event parameter you wish to surface in the alerts.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
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
						Description:  "The type of the entity.",
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
									Description:  "The identifier of the entity.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"column_name": {
									Description:  "The column name to be mapped to the identifier.",
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
						Description:  "The column name to be mapped to the identifier.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
