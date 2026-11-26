// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cost_anomaly_alert

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Cost Anomaly Alert. The name can contain only lowercase letters, numbers and hyphens.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CostAnomalyAlertName,
			// action names can contain only alphanumeric characters and hyphens.
		},

		"display_name": {
			Description:  "The display name which should be used for this Cost Anomaly Alert.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 25),
		},

		"subscription_id": {
			Description:  "The ID of the Subscription this Cost Anomaly Alert is scoped to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"notification_email": {
			Description:  "The email address of the point of contact that should get the unsubscribe requests and notification emails.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_subject": {
			Description:  "The email subject of the Cost Anomaly Alerts. Maximum length of the subject is 50.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"email_addresses": {
			Description: "A list of email addresses which the Anomaly Alerts are sent to.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"message": {
			Description:  "The message of the Cost Anomaly Alert. Maximum length of the message is 100.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 250),
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
