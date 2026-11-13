// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorMetricAlertSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"scopes": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
			Set: pluginsdk.HashString,
		},

		"target_resource_type": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Description: `The resource type (e.g. Microsoft.Compute/virtualMachines) of the target pluginsdk. Required when using subscription, resource group scope or multiple scopes.`,
		},

		"target_resource_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			Description:      `The location of the target pluginsdk. Required when using subscription, resource group scope or multiple scopes.`,
		},

		// static criteria
		"criteria": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_namespace": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"metric_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"aggregation": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Average",
							"Count",
							"Minimum",
							"Maximum",
							"Total",
						}, false),
					},
					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
										"StartsWith",
									}, false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"operator": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForOperator(), false),
					},
					"threshold": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
					},
					"skip_metric_validation": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		// lintignore: S018
		"dynamic_criteria": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			// Currently, it allows to define only one dynamic criteria in one metric alert.
			MaxItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_namespace": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"metric_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"aggregation": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Average",
							"Count",
							"Minimum",
							"Maximum",
							"Total",
						}, false),
					},
					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
										"StartsWith",
									}, false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"operator": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForDynamicThresholdOperator(), false),
					},
					"alert_sensitivity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForDynamicThresholdSensitivity(), false),
					},

					"evaluation_total_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      4,
					},

					"evaluation_failure_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      4,
					},

					"ignore_data_before": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"skip_metric_validation": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"application_insights_web_test_location_availability_criteria": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			MaxItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"web_test_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: webtestsapis.ValidateWebTestID,
					},
					"component_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: componentsapis.ValidateComponentID,
					},
					"failed_location_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"action": {
			Type:     pluginsdk.TypeSet,
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
			Set: resourceMonitorMetricAlertActionHash,
		},

		"auto_mitigate": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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

		"frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "PT1M",
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT15M",
				"PT30M",
				"PT1H",
			}, false),
		},

		"severity": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"window_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "PT5M",
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT15M",
				"PT30M",
				"PT1H",
				"PT6H",
				"PT12H",
				"P1D",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
