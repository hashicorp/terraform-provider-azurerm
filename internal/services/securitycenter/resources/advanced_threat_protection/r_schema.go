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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
	}
}
