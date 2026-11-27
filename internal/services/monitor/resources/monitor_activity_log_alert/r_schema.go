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
			Description:  "The name of the activity log alert.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"scopes": {
			Description: "The Scope at which the Activity Log should be applied. A list of strings which could be a Resource Group, or a subscription, or a resource ID (such as a Storage Account).",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Set:         pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Description: "The category of the operation.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description:  "The Resource Manager Role-Based Access Control operation name. Supported operation should be of the form: `<resourceProvider>/<resourceType>/<operation>`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"caller": {
						Description:  "The email address or Azure Active Directory identifier of the user who performed the operation.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"level": {
						Description: "The severity level of the event.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description:   "A list of severity level of the event.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.level"},
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
					},
					"resource_provider": {
						Description:   "The name of the Resource Provider monitored by the activity log alert.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_providers"},
					},
					"resource_providers": {
						Description:   "A list of names of Resource Providers monitored by the activity log alert.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.resource_provider"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"resource_type": {
						Description:   "The resource type monitored by the activity log alert.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_types"},
					},
					"resource_types": {
						Description:   "A list of resource types monitored by the activity log alert.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.resource_type"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"resource_group": {
						Description:   "The name of Resource Group monitored by the activity log alert.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.resource_groups"},
					},
					"resource_groups": {
						Description:   "A list of names of Resource Groups monitored by the activity log alert.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.resource_group"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"resource_id": {
						Description:   "The specific resource monitored by the activity log alert.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  azure.ValidateResourceID,
						ConflictsWith: []string{"criteria.0.resource_ids"},
					},
					"resource_ids": {
						Description:   "A list of specific resources monitored by the activity log alert.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.resource_id"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"status": {
						Description:   "The status of the event. For example, `Started`, `Failed`, or `Succeeded`.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.statuses"},
					},
					"statuses": {
						Description:   "A list of status of the event. For example, `Started`, `Failed`, or `Succeeded`.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.status"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"sub_status": {
						Description:   "The sub status of the event.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.sub_statuses"},
					},
					"sub_statuses": {
						Description:   "A list of sub status of the event.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ConflictsWith: []string{"criteria.0.sub_status"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"recommendation_category": {
						Description: "The recommendation category of the event.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description: "The recommendation impact of the event.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description:   "The recommendation type of the event.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"criteria.0.recommendation_category", "criteria.0.recommendation_impact"},
					},
					// lintignore:XS003
					"resource_health": {
						Description:   "A block to define fine grain resource health settings.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						MaxItems:      1,
						ConflictsWith: []string{"criteria.0.caller", "criteria.0.service_health"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"current": {
									Description: "The current resource health statuses that will log an alert.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
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
								},
								"previous": {
									Description: "The previous resource health statuses that will log an alert.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
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
								},
								"reason": {
									Description: "The reason that will log an alert.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
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
								},
							},
						},
					},
					// lintignore:XS003
					"service_health": {
						Description:   "A block to define fine grain service health settings.",
						Type:          pluginsdk.TypeList,
						Optional:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						MaxItems:      1,
						ConflictsWith: []string{"criteria.0.caller", "criteria.0.resource_health"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"events": {
									Description: "Events this alert will monitor.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
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
								},
								"locations": {
									Description: "Locations this alert will monitor. For example, `West Europe`.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"services": {
									Description: "Services this alert will monitor. For example, `Activity Logs & Alerts`, `Action Groups`.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									Set:         pluginsdk.HashString,
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

		"action": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group_id": {
						Description:  "The ID of the Action Group can be sourced from [the `azurerm_monitor_action_group` resource](./monitor_action_group.html).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"webhook_properties": {
						Description: "The map of custom string properties to include with the post operation. These data are appended to the webhook payload.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"description": {
			Description: "The description of this activity log alert.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"enabled": {
			Description: "Whether this Activity Log Alert should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),
	}
}
