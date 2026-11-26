// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_gallery

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_id": {
			Description: "The ID of the Dev Center within which this Dev Center Gallery should exist.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"shared_gallery_id": {
			Description: "The ID of the Shared Gallery which should be connected to the Dev Center Gallery.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"name": {
			Description: "The name of this Dev Center Gallery.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
