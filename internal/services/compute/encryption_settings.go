// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func encryptionSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"disk_encryption_key": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_url": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
				"key_encryption_key": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"key_url": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func expandSnapshotDiskEncryptionSettings(settingsList []any) *snapshots.EncryptionSettingsCollection {
	if len(settingsList) == 0 {
		return &snapshots.EncryptionSettingsCollection{}
	}
	settings := settingsList[0].(map[string]any)

	config := &snapshots.EncryptionSettingsCollection{
		Enabled: true,
	}

	var diskEncryptionKey *snapshots.KeyVaultAndSecretReference
	if v := settings["disk_encryption_key"].([]any); len(v) > 0 {
		dek := v[0].(map[string]any)

		secretURL := dek["secret_url"].(string)
		sourceVaultId := dek["source_vault_id"].(string)
		diskEncryptionKey = &snapshots.KeyVaultAndSecretReference{
			SecretURL: secretURL,
			SourceVault: snapshots.SourceVault{
				Id: pointer.To(sourceVaultId),
			},
		}
	}

	var keyEncryptionKey *snapshots.KeyVaultAndKeyReference
	if v := settings["key_encryption_key"].([]any); len(v) > 0 {
		kek := v[0].(map[string]any)

		secretURL := kek["key_url"].(string)
		sourceVaultId := kek["source_vault_id"].(string)
		keyEncryptionKey = &snapshots.KeyVaultAndKeyReference{
			KeyURL: secretURL,
			SourceVault: snapshots.SourceVault{
				Id: pointer.To(sourceVaultId),
			},
		}
	}

	// at this time we only support a single element
	config.EncryptionSettings = &[]snapshots.EncryptionSettingsElement{
		{
			DiskEncryptionKey: diskEncryptionKey,
			KeyEncryptionKey:  keyEncryptionKey,
		},
	}
	return config
}

func flattenSnapshotDiskEncryptionSettings(encryptionSettings *snapshots.EncryptionSettingsCollection) []any {
	if encryptionSettings == nil {
		return []any{}
	}

	diskEncryptionKeys := make([]any, 0)
	keyEncryptionKeys := make([]any, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretURL != "" {
				secretUrl = key.SecretURL
			}

			sourceVaultId := pointer.From(key.SourceVault.Id)

			diskEncryptionKeys = append(diskEncryptionKeys, map[string]any{
				"secret_url":      secretUrl,
				"source_vault_id": sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyURL != "" {
				keyUrl = key.KeyURL
			}

			sourceVaultId := pointer.From(key.SourceVault.Id)

			keyEncryptionKeys = append(keyEncryptionKeys, map[string]any{
				"key_url":         keyUrl,
				"source_vault_id": sourceVaultId,
			})
		}
	}

	if len(diskEncryptionKeys) > 0 {
		return []any{
			map[string]any{
				"disk_encryption_key": diskEncryptionKeys,
				"key_encryption_key":  keyEncryptionKeys,
			},
		}
	} else {
		return []any{}
	}
}

func expandManagedDiskEncryptionSettings(settingsList []any) *disks.EncryptionSettingsCollection {
	if len(settingsList) == 0 {
		return &disks.EncryptionSettingsCollection{}
	}
	settings := settingsList[0].(map[string]any)

	config := &disks.EncryptionSettingsCollection{
		Enabled: true,
	}

	var diskEncryptionKey *disks.KeyVaultAndSecretReference
	if v := settings["disk_encryption_key"].([]any); len(v) > 0 {
		dek := v[0].(map[string]any)

		secretURL := dek["secret_url"].(string)
		sourceVaultId := dek["source_vault_id"].(string)
		diskEncryptionKey = &disks.KeyVaultAndSecretReference{
			SecretURL: secretURL,
			SourceVault: disks.SourceVault{
				Id: pointer.To(sourceVaultId),
			},
		}
	}

	var keyEncryptionKey *disks.KeyVaultAndKeyReference
	if v := settings["key_encryption_key"].([]any); len(v) > 0 {
		kek := v[0].(map[string]any)

		secretURL := kek["key_url"].(string)
		sourceVaultId := kek["source_vault_id"].(string)
		keyEncryptionKey = &disks.KeyVaultAndKeyReference{
			KeyURL: secretURL,
			SourceVault: disks.SourceVault{
				Id: pointer.To(sourceVaultId),
			},
		}
	}

	// at this time we only support a single element
	config.EncryptionSettings = &[]disks.EncryptionSettingsElement{
		{
			DiskEncryptionKey: diskEncryptionKey,
			KeyEncryptionKey:  keyEncryptionKey,
		},
	}
	return config
}

func flattenManagedDiskEncryptionSettings(encryptionSettings *disks.EncryptionSettingsCollection) []any {
	if encryptionSettings == nil {
		return []any{}
	}

	diskEncryptionKeys := make([]any, 0)
	keyEncryptionKeys := make([]any, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretURL != "" {
				secretUrl = key.SecretURL
			}

			sourceVaultId := pointer.From(key.SourceVault.Id)

			diskEncryptionKeys = append(diskEncryptionKeys, map[string]any{
				"secret_url":      secretUrl,
				"source_vault_id": sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyURL != "" {
				keyUrl = key.KeyURL
			}

			sourceVaultId := pointer.From(key.SourceVault.Id)

			keyEncryptionKeys = append(keyEncryptionKeys, map[string]any{
				"key_url":         keyUrl,
				"source_vault_id": sourceVaultId,
			})
		}
	}

	if len(diskEncryptionKeys) > 0 {
		return []any{
			map[string]any{
				"disk_encryption_key": diskEncryptionKeys,
				"key_encryption_key":  keyEncryptionKeys,
			},
		}
	} else {
		return []any{}
	}
}
