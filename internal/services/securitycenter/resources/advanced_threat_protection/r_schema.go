// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advanced_threat_protection

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func advancedThreatProtectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_resource_id": {
			Description:  "The ID of the Azure Resource which to enable Advanced Threat Protection on.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"enabled": {
			Description: "Whether Advanced Threat Protection should be enabled on this resource.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},
	}
}
