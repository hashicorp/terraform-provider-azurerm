// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storagecacheClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := storagecacheClient.From(metadata.Client).AmlFilesystems
			subscriptionId := account.From(metadata.Client).SubscriptionId
			id := amlfilesystems.NewAmlFilesystemID(subscriptionId, model.ResourceGroupName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return err
			}

			properties := &amlfilesystems.AmlFilesystem{
				Location: location.Normalize(model.Location),
				Identity: identity,
				Properties: &amlfilesystems.AmlFilesystemProperties{
					Hsm:                expandManagedLustreFileSystemHsmSetting(model.HsmSetting),
					EncryptionSettings: expandManagedLustreFileSystemEncryptionKey(model.EncryptionKey),
					MaintenanceWindow:  expandManagedLustreFileSystemMaintenanceWindowForCreate(model.MaintenanceWindow),
					FilesystemSubnet:   model.SubnetId,
					StorageCapacityTiB: float64(model.StorageCapacityInTb),
				},
				Sku: &amlfilesystems.SkuName{
					Name: pointer.To(model.SkuName),
				},
				Zones: pointer.To(model.Zones),
				Tags:  pointer.To(model.Tags),
			}

			if model.RootSquashSettings != nil {
				properties.Properties.RootSquashSettings = expandRootSquashSettings(model.RootSquashSettings)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandManagedLustreFileSystemMaintenanceWindowForCreate(input []MaintenanceWindow) amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow {
	maintenanceWindow := &input[0]

	return amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow{
		DayOfWeek:    pointer.To(maintenanceWindow.DayOfWeek),
		TimeOfDayUTC: pointer.To(maintenanceWindow.TimeOfDayInUTC),
	}
}

func expandManagedLustreFileSystemMaintenanceWindowForUpdate(input []MaintenanceWindow) *amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow {
	if len(input) == 0 {
		return nil
	}

	maintenanceWindow := &input[0]

	return &amlfilesystems.AmlFilesystemUpdatePropertiesMaintenanceWindow{
		DayOfWeek:    pointer.To(maintenanceWindow.DayOfWeek),
		TimeOfDayUTC: pointer.To(maintenanceWindow.TimeOfDayInUTC),
	}
}

func expandManagedLustreFileSystemEncryptionKey(input []EncryptionKey) *amlfilesystems.AmlFilesystemEncryptionSettings {
	if len(input) == 0 {
		return nil
	}

	encryptionKey := &input[0]

	result := &amlfilesystems.KeyVaultKeyReference{
		KeyURL: encryptionKey.KeyURL,
		SourceVault: amlfilesystems.KeyVaultKeyReferenceSourceVault{
			Id: pointer.To(encryptionKey.SourceVaultId),
		},
	}

	return &amlfilesystems.AmlFilesystemEncryptionSettings{
		KeyEncryptionKey: result,
	}
}

func expandRootSquashSettings(input []RootSquashSetting) *amlfilesystems.AmlFilesystemRootSquashSettings {
	if len(input) == 0 {
		return &amlfilesystems.AmlFilesystemRootSquashSettings{
			Mode: pointer.To(amlfilesystems.AmlFilesystemSquashModeNone),
		}
	}

	rootSquashSetting := &input[0]

	return &amlfilesystems.AmlFilesystemRootSquashSettings{
		Mode:             pointer.ToEnum[amlfilesystems.AmlFilesystemSquashMode](rootSquashSetting.Mode),
		NoSquashNidLists: pointer.To(rootSquashSetting.NoSquashNidList),
		SquashGID:        pointer.To(rootSquashSetting.SquashGID),
		SquashUID:        pointer.To(rootSquashSetting.SquashUID),
	}
}

func expandManagedLustreFileSystemHsmSetting(input []HsmSetting) *amlfilesystems.AmlFilesystemPropertiesHsm {
	if len(input) == 0 {
		return nil
	}

	hsmSetting := &input[0]

	result := &amlfilesystems.AmlFilesystemHsmSettings{
		Container:        hsmSetting.ContainerId,
		LoggingContainer: hsmSetting.LoggingContainerId,
	}

	if hsmSetting.ImportPrefix != "" {
		result.ImportPrefix = pointer.To(hsmSetting.ImportPrefix)
	}

	return &amlfilesystems.AmlFilesystemPropertiesHsm{
		Settings: result,
	}
}
