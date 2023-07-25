// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func encryptionSettingsSchema() *pluginsdk.Schema {
	if !features.FourPointOhBeta() {
		return &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,

						// Azure can change enabled from false to true, but not the other way around, so
						//   to keep idempotency, we'll conservatively set this to ForceNew=true
						ForceNew:   true,
						Deprecated: "Deprecated, Azure Disk Encryption is now configured directly by `disk_encryption_key` and `key_encryption_key`. To disable Azure Disk Encryption, please remove `encryption_settings` block. To enabled, specify a `encryption_settings` block`",
					},
					"disk_encryption_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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

func expandSnapshotDiskEncryptionSettings(settingsList []interface{}) *snapshots.EncryptionSettingsCollection {
	if len(settingsList) == 0 {
		return &snapshots.EncryptionSettingsCollection{}
	}
	settings := settingsList[0].(map[string]interface{})

	config := &snapshots.EncryptionSettingsCollection{
		Enabled: true,
	}

	if !features.FourPointOhBeta() {
		config.Enabled = settings["enabled"].(bool)
	}

	var diskEncryptionKey *snapshots.KeyVaultAndSecretReference
	if v := settings["disk_encryption_key"].([]interface{}); len(v) > 0 {
		dek := v[0].(map[string]interface{})

		secretURL := dek["secret_url"].(string)
		sourceVaultId := dek["source_vault_id"].(string)
		diskEncryptionKey = &snapshots.KeyVaultAndSecretReference{
			SecretUrl: secretURL,
			SourceVault: snapshots.SourceVault{
				Id: utils.String(sourceVaultId),
			},
		}
	}

	var keyEncryptionKey *snapshots.KeyVaultAndKeyReference
	if v := settings["key_encryption_key"].([]interface{}); len(v) > 0 {
		kek := v[0].(map[string]interface{})

		secretURL := kek["key_url"].(string)
		sourceVaultId := kek["source_vault_id"].(string)
		keyEncryptionKey = &snapshots.KeyVaultAndKeyReference{
			KeyUrl: secretURL,
			SourceVault: snapshots.SourceVault{
				Id: utils.String(sourceVaultId),
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

func flattenSnapshotDiskEncryptionSettings(encryptionSettings *snapshots.EncryptionSettingsCollection) []interface{} {
	if encryptionSettings == nil {
		return []interface{}{}
	}

	diskEncryptionKeys := make([]interface{}, 0)
	keyEncryptionKeys := make([]interface{}, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretUrl != "" {
				secretUrl = key.SecretUrl
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			diskEncryptionKeys = append(diskEncryptionKeys, map[string]interface{}{
				"secret_url":      secretUrl,
				"source_vault_id": sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyUrl != "" {
				keyUrl = key.KeyUrl
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			keyEncryptionKeys = append(keyEncryptionKeys, map[string]interface{}{
				"key_url":         keyUrl,
				"source_vault_id": sourceVaultId,
			})
		}
	}

	if len(diskEncryptionKeys) > 0 {
		if !features.FourPointOhBeta() {
			return []interface{}{
				map[string]interface{}{
					"enabled":             true,
					"disk_encryption_key": diskEncryptionKeys,
					"key_encryption_key":  keyEncryptionKeys,
				},
			}
		}

		return []interface{}{
			map[string]interface{}{
				"disk_encryption_key": diskEncryptionKeys,
				"key_encryption_key":  keyEncryptionKeys,
			},
		}
	} else {
		return []interface{}{}
	}
}

func expandManagedDiskEncryptionSettings(settingsList []interface{}) *disks.EncryptionSettingsCollection {
	if len(settingsList) == 0 {
		return &disks.EncryptionSettingsCollection{}
	}
	settings := settingsList[0].(map[string]interface{})

	config := &disks.EncryptionSettingsCollection{
		Enabled: true,
	}

	if !features.FourPointOhBeta() {
		config.Enabled = settings["enabled"].(bool)
	}

	var diskEncryptionKey *disks.KeyVaultAndSecretReference
	if v := settings["disk_encryption_key"].([]interface{}); len(v) > 0 {
		dek := v[0].(map[string]interface{})

		secretURL := dek["secret_url"].(string)
		sourceVaultId := dek["source_vault_id"].(string)
		diskEncryptionKey = &disks.KeyVaultAndSecretReference{
			SecretUrl: secretURL,
			SourceVault: disks.SourceVault{
				Id: utils.String(sourceVaultId),
			},
		}
	}

	var keyEncryptionKey *disks.KeyVaultAndKeyReference
	if v := settings["key_encryption_key"].([]interface{}); len(v) > 0 {
		kek := v[0].(map[string]interface{})

		secretURL := kek["key_url"].(string)
		sourceVaultId := kek["source_vault_id"].(string)
		keyEncryptionKey = &disks.KeyVaultAndKeyReference{
			KeyUrl: secretURL,
			SourceVault: disks.SourceVault{
				Id: utils.String(sourceVaultId),
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

func flattenManagedDiskEncryptionSettings(encryptionSettings *disks.EncryptionSettingsCollection) []interface{} {
	if encryptionSettings == nil {
		return []interface{}{}
	}

	diskEncryptionKeys := make([]interface{}, 0)
	keyEncryptionKeys := make([]interface{}, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretUrl != "" {
				secretUrl = key.SecretUrl
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			diskEncryptionKeys = append(diskEncryptionKeys, map[string]interface{}{
				"secret_url":      secretUrl,
				"source_vault_id": sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyUrl != "" {
				keyUrl = key.KeyUrl
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			keyEncryptionKeys = append(keyEncryptionKeys, map[string]interface{}{
				"key_url":         keyUrl,
				"source_vault_id": sourceVaultId,
			})
		}
	}

	if len(diskEncryptionKeys) > 0 {
		if !features.FourPointOhBeta() {
			return []interface{}{
				map[string]interface{}{
					"enabled":             true,
					"disk_encryption_key": diskEncryptionKeys,
					"key_encryption_key":  keyEncryptionKeys,
				},
			}
		}

		return []interface{}{
			map[string]interface{}{
				"disk_encryption_key": diskEncryptionKeys,
				"key_encryption_key":  keyEncryptionKeys,
			},
		}
	} else {
		return []interface{}{}
	}
}
