// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_smart_detector_alert_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorSmartDetectorAlertRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Monitor Smart Detector Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"detector_type": {
			Description: "The Built-In Smart Detector type that this alert rule will use.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"FailureAnomaliesDetector",
				"RequestPerformanceDegradationDetector",
				"DependencyPerformanceDegradationDetector",
				"ExceptionVolumeChangedDetector",
				"TraceSeverityDetector",
				"MemoryLeakDetector",
			}, false),
		},

		"scope_resource_ids": {
			Description: "The scopes of this Smart Detector Alert Rule.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Set:         set.HashStringIgnoreCase,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"severity": {
			Description: "The severity of this Smart Detector Alert Rule.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice(
				smartdetectoralertrules.PossibleValuesForSeverity(), false,
			),
		},

		"frequency": {
			Description:  "The frequency of this Smart Detector Alert Rule in ISO8601 format.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"action_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ids": {
						Description: "The action group IDs.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Set:         set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.AsGeneratedID(actiongroupsapis.ParseActionGroupIDInsensitively),
						},
					},

					"email_subject": {
						Description: "A custom email subject if Email Receiver is specified in Monitor Action Group resource.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"webhook_payload": {
						Description:      "A JSON String which Specifies the custom webhook payload if Webhook Receiver is specified in Monitor Action Group resource.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},
				},
			},
		},

		"description": {
			Description:  "A description for the Smart Detector Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Description: "Whether the Smart Detector Alert Rule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"throttling_duration": {
			Description:  "The duration (in ISO8601 format) to wait before notifying on the alert rule again.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"tags": commonschema.Tags(),
	}
}
