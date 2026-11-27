// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_contact

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterContactSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Security Center Contact.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email": {
			Description:  "The email of the Security Center Contact.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"phone": {
			Description:  "The phone number of the Security Center Contact.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"alert_notifications": {
			Description: "Whether to send security alerts notifications to the security contact.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"alerts_to_admins": {
			Description: "Whether to send security alerts notifications to subscription admins.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},
	}
}
