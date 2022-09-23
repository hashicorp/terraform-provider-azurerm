package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func encryptionSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,

					// Azure can change enabled from false to true, but not the other way around, so
					//   to keep idempotency, we'll conservatively set this to ForceNew=true
					ForceNew: true,
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

func expandSnapshotDiskEncryptionSettings(settings map[string]interface{}) *compute.EncryptionSettingsCollection {
	enabled := settings["enabled"].(bool)
	config := &compute.EncryptionSettingsCollection{
		Enabled: utils.Bool(enabled),
	}

	var diskEncryptionKey *compute.KeyVaultAndSecretReference
	if v := settings["disk_encryption_key"].([]interface{}); len(v) > 0 {
		dek := v[0].(map[string]interface{})

		secretURL := dek["secret_url"].(string)
		sourceVaultId := dek["source_vault_id"].(string)
		diskEncryptionKey = &compute.KeyVaultAndSecretReference{
			SecretURL: utils.String(secretURL),
			SourceVault: &compute.SourceVault{
				ID: utils.String(sourceVaultId),
			},
		}
	}

	var keyEncryptionKey *compute.KeyVaultAndKeyReference
	if v := settings["key_encryption_key"].([]interface{}); len(v) > 0 {
		kek := v[0].(map[string]interface{})

		secretURL := kek["key_url"].(string)
		sourceVaultId := kek["source_vault_id"].(string)
		keyEncryptionKey = &compute.KeyVaultAndKeyReference{
			KeyURL: utils.String(secretURL),
			SourceVault: &compute.SourceVault{
				ID: utils.String(sourceVaultId),
			},
		}
	}

	// at this time we only support a single element
	config.EncryptionSettings = &[]compute.EncryptionSettingsElement{
		{
			DiskEncryptionKey: diskEncryptionKey,
			KeyEncryptionKey:  keyEncryptionKey,
		},
	}
	return config
}

func flattenSnapshotDiskEncryptionSettings(encryptionSettings *compute.EncryptionSettingsCollection) []interface{} {
	if encryptionSettings == nil {
		return []interface{}{}
	}

	enabled := false
	if encryptionSettings.Enabled != nil {
		enabled = *encryptionSettings.Enabled
	}

	diskEncryptionKeys := make([]interface{}, 0)
	keyEncryptionKeys := make([]interface{}, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretURL != nil {
				secretUrl = *key.SecretURL
			}

			sourceVaultId := ""
			if key.SourceVault != nil && key.SourceVault.ID != nil {
				sourceVaultId = *key.SourceVault.ID
			}

			diskEncryptionKeys = append(diskEncryptionKeys, map[string]interface{}{
				"secret_url":      secretUrl,
				"source_vault_id": sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyURL != nil {
				keyUrl = *key.KeyURL
			}

			sourceVaultId := ""
			if key.SourceVault != nil && key.SourceVault.ID != nil {
				sourceVaultId = *key.SourceVault.ID
			}

			keyEncryptionKeys = append(keyEncryptionKeys, map[string]interface{}{
				"key_url":         keyUrl,
				"source_vault_id": sourceVaultId,
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":             enabled,
			"disk_encryption_key": diskEncryptionKeys,
			"key_encryption_key":  keyEncryptionKeys,
		},
	}
}

func expandManagedDiskEncryptionSettings(settings map[string]interface{}) *disks.EncryptionSettingsCollection {
	enabled := settings["enabled"].(bool)
	config := &disks.EncryptionSettingsCollection{
		Enabled: enabled,
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

	enabled := false
	if encryptionSettings.Enabled {
		enabled = encryptionSettings.Enabled
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

	return []interface{}{
		map[string]interface{}{
			"enabled":             enabled,
			"disk_encryption_key": diskEncryptionKeys,
			"key_encryption_key":  keyEncryptionKeys,
		},
	}
}
