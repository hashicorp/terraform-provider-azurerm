// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorActivityLogAlertSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scopes": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Set: pluginsdk.HashString,
		},

		"criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Administrative",
							"Autoscale",
							"Policy",
							"Recommendation",
							"ResourceHealth",
							"Security",
							"ServiceHealth",
						}, false),
					},
					"operation_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"caller": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"level": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Verbose",
							"Informational",
							"Warning",
							"Error",
							"Critical",
						}, false),
						ConflictsWith: []string{"criteria.0.levels"},
					},
					"levels": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Verbose",
								"Informational",
								"Warning",
								"Error",
								"Critical",
							}, false),
						},
						ConflictsWith: []string{"criteria.0.level"},
					},
					"resource_provider": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_providers"},
					},
					"resource_providers": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.resource_provider"},
					},
					"resource_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_types"},
					},
					"resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.resource_type"},
					},
					"resource_group": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_groups"},
					},
					"resource_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.resource_group"},
					},
					"resource_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  azure.ValidateResourceID,
						ConflictsWith: []string{"criteria.0.resource_ids"},
					},
					"resource_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.resource_id"},
					},
					"status": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.statuses"},
					},
					"statuses": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.status"},
					},
					"sub_status": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.sub_statuses"},
					},
					"sub_statuses": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						ConflictsWith: []string{"criteria.0.sub_status"},
					},
					"recommendation_category": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							[]string{
								"Cost",
								"Reliability",
								"OperationalExcellence",
								"Performance",
								"HighAvailability",
								"Security",
							},
							false,
						),
						ConflictsWith: []string{"criteria.0.recommendation_type"},
					},
					"recommendation_impact": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							[]string{
								"High",
								"Medium",
								"Low",
							},
							false,
						),
						ConflictsWith: []string{"criteria.0.recommendation_type"},
					},
					"recommendation_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.recommendation_category", "criteria.0.recommendation_impact"},
					},
					// lintignore:XS003
					"resource_health": {
						Type:     pluginsdk.TypeList,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"current": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"Available",
												"Degraded",
												"Unavailable",
												"Unknown",
											},
											false,
										),
									},
									Set: pluginsdk.HashString,
								},
								"previous": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"Available",
												"Degraded",
												"Unavailable",
												"Unknown",
											},
											false,
										),
									},
									Set: pluginsdk.HashString,
								},
								"reason": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"PlatformInitiated",
												"UserInitiated",
												"Unknown",
											},
											false,
										),
									},
									Set: pluginsdk.HashString,
								},
							},
						},
						ConflictsWith: []string{"criteria.0.caller", "criteria.0.service_health"},
					},
					// lintignore:XS003
					"service_health": {
						Type:     pluginsdk.TypeList,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"events": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"Incident",
												"Maintenance",
												"Informational",
												"ActionRequired",
												"Security",
											},
											false,
										),
									},
									Set: pluginsdk.HashString,
								},
								"locations": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									Set: pluginsdk.HashString,
								},
								"services": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									Set: pluginsdk.HashString,
								},
							},
						},
						ConflictsWith: []string{"criteria.0.caller", "criteria.0.resource_health"},
					},
				},
			},
		},

		"action": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"webhook_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}
