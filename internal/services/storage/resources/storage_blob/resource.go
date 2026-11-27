// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_blob/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageBlobCreate,
		Read:   resourceStorageBlobRead,
		Update: resourceStorageBlobUpdate,
		Delete: resourceStorageBlobDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.BlobV0ToV1{},
		}),

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := blobs.ParseBlobID(id, storageDomainSuffix)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageBlobSchema(),

		CustomizeDiff: func(ctx context.Context, diff *pluginsdk.ResourceDiff, i any) error {
			if content := diff.Get("source_content"); content != "" && diff.Get("type") == "Page" {
				if len(content.(string))%512 != 0 {
					return errors.New(`"source" must be aligned to 512-byte boundary for "type" set to "Page"`)
				}
			}
			return nil
		},
	}
}
