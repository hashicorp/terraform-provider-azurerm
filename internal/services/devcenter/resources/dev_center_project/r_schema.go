// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_id": {
			Description: "Resource ID of an associated DevCenter.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"location": commonschema.Location(),
		"name": {
			Description: "The name of this Dev Center Project.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"description": {
			Description: "Description of the project.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"maximum_dev_boxes_per_user": {
			Description: "When specified, limits the maximum number of Dev Boxes a single user can create across all pools in the project.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Description: "The URI of the Dev Center resource this project is associated with.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
