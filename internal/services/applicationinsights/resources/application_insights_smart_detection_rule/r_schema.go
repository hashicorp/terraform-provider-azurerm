// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_smart_detection_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationInsightsSmartDetectionRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Slow page load time",
				"Slow server response time",
				"Long dependency duration",
				"Degradation in server response time",
				"Degradation in dependency duration",
				// The below rules are currently preview and may change in future
				"Degradation in trace severity ratio",
				"Abnormal rise in exception volume",
				"Potential memory leak detected",
				"Potential security issue detected",
				"Abnormal rise in daily data volume",
			}, false),
			DiffSuppressFunc: smartDetectionRuleNameDiff,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"send_emails_to_subscription_owners": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"additional_email_recipients": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}
