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
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"description": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"maximum_dev_boxes_per_user": {
			Optional: true,
			Type:     pluginsdk.TypeInt,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
