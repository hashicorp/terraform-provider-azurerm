package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
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
				// TODO: remove in 3.0
				"enabled": {
					Type:       pluginsdk.TypeBool,
					Optional:   true,
					Default:    true,
					Deprecated: "Deprecated, Azure Disk Encryption is now configured directly by `disk_encryption_key` and `key_encryption_key`. To disable Azure Disk Encryption, please remove `encryption_settings` block",
				},

				// TODO: make this Required in 3.0
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

func expandManagedDiskEncryptionSettings(settingsList []interface{}) *compute.EncryptionSettingsCollection {
	if len(settingsList) == 0 {
		return &compute.EncryptionSettingsCollection{}
	}
	settings := settingsList[0].(map[string]interface{})

	config := &compute.EncryptionSettingsCollection{
		// TODO: update to `Enabled: utils.Bool(true),` in 3.0
		Enabled: utils.Bool(settings["enabled"].(bool)),
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

func flattenManagedDiskEncryptionSettings(encryptionSettings *compute.EncryptionSettingsCollection) []interface{} {
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

	if len(diskEncryptionKeys) > 0 {
		return []interface{}{
			map[string]interface{}{
				// TODO: remove `enabled` assignment in 3.0
				"enabled":             true,
				"disk_encryption_key": diskEncryptionKeys,
				"key_encryption_key":  keyEncryptionKeys,
			},
		}
	} else {
		return []interface{}{}
	}
}
