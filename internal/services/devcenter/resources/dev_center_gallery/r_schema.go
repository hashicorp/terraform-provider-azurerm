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
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"shared_gallery_id": {
			Description: "The ID of the Shared Gallery which should be connected to the Dev Center Gallery.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"name": {
			Description: "The name of this Dev Center Gallery.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
