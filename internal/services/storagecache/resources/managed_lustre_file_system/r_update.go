// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storagecacheClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storagecacheClient.From(metadata.Client).AmlFilesystems

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := amlfilesystems.AmlFilesystemUpdate{
				Properties: &amlfilesystems.AmlFilesystemUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("maintenance_window") {
				properties.Properties.MaintenanceWindow = expandManagedLustreFileSystemMaintenanceWindowForUpdate(model.MaintenanceWindow)
			}

			if metadata.ResourceData.HasChange("encryption_key") {
				properties.Properties.EncryptionSettings = expandManagedLustreFileSystemEncryptionKey(model.EncryptionKey)
			}

			if metadata.ResourceData.HasChange("root_squash") {
				properties.Properties.RootSquashSettings = expandRootSquashSettings(model.RootSquashSettings)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
