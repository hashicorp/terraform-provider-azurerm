// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_target

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"target_resource_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"target_type": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
