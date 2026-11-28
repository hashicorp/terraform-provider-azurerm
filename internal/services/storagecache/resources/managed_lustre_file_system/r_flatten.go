// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
)

func flattenManagedLustreFileSystemMaintenanceWindow(input amlfilesystems.AmlFilesystemPropertiesMaintenanceWindow) []MaintenanceWindow {
	var result []MaintenanceWindow

	maintenanceWindow := MaintenanceWindow{
		DayOfWeek:      pointer.From(input.DayOfWeek),
		TimeOfDayInUTC: pointer.From(input.TimeOfDayUTC),
	}

	return append(result, maintenanceWindow)
}

func flattenRootSquashSettings(input *amlfilesystems.AmlFilesystemRootSquashSettings) []RootSquashSetting {
	result := make([]RootSquashSetting, 0)
	if input == nil || pointer.From(input.Mode) == amlfilesystems.AmlFilesystemSquashModeNone {
		return []RootSquashSetting{}
	}

	rootSquashSetting := RootSquashSetting{
		Mode:            pointer.FromEnum(input.Mode),
		NoSquashNidList: pointer.From(input.NoSquashNidLists),
		SquashGID:       pointer.From(input.SquashGID),
		SquashUID:       pointer.From(input.SquashUID),
	}
	return append(result, rootSquashSetting)
}

func flattenManagedLustreFileSystemEncryptionKey(input *amlfilesystems.AmlFilesystemEncryptionSettings) []EncryptionKey {
	result := make([]EncryptionKey, 0)
	if input == nil || input.KeyEncryptionKey == nil {
		return result
	}

	encryptionKey := EncryptionKey{
		KeyURL:        input.KeyEncryptionKey.KeyURL,
		SourceVaultId: pointer.From(input.KeyEncryptionKey.SourceVault.Id),
	}

	return append(result, encryptionKey)
}

func flattenManagedLustreFileSystemHsmSetting(input *amlfilesystems.AmlFilesystemPropertiesHsm) []HsmSetting {
	result := make([]HsmSetting, 0)
	if input == nil || input.Settings == nil {
		return result
	}

	hsmSetting := HsmSetting{}

	if v := input.Settings; v != nil {
		hsmSetting.ContainerId = v.Container
		hsmSetting.LoggingContainerId = v.LoggingContainer
		hsmSetting.ImportPrefix = pointer.From(v.ImportPrefix)
	}

	return append(result, hsmSetting)
}
