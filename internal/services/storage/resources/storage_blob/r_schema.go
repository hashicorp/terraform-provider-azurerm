// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func storageBlobSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_container_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Append",
				"Block",
				"Page",
			}, false),
		},

		"size": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      0,
			ValidateFunc: validation.IntDivisibleBy(512),
		},

		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(blobs.Archive),
				string(blobs.Cool),
				string(blobs.Hot),
			}, false),
		},

		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "application/octet-stream",
		},

		"cache_control": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"encryption_scope": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageEncryptionScopeName,
		},

		"source": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"source_uri", "source_content"},
		},

		"source_content": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"source", "source_uri"},
		},

		"source_uri": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"source", "source_content"},
		},

		"content_md5": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"source_uri"},
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parallelism": {
			// TODO: @tombuildsstuff - a note this only works for Page blobs
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      8,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"metadata": helpers.MetaDataComputedSchema(),
	}
}
