// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_capability

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"capability_type": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"chaos_studio_target_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: commonids.ValidateChaosStudioTargetID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"urn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
