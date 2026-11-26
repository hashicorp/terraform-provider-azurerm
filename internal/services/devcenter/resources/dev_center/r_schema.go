// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			Description: "The name of this Dev Center.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"project_catalog_item_sync_enabled": {
			Description: "Whether the project catalogs associated with projects in this Dev Center are allowed to sync catalog items.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"tags":     commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Description: "The URI of the Dev Center.",
			Computed:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}
