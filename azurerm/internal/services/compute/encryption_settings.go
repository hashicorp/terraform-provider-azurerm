package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-30/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func encryptionSettingsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Required: true,

					// Azure can change enabled from false to true, but not the other way around, so
					//   to keep idempotency, we'll conservatively set this to ForceNew=true
					ForceNew: true,
				},

				"disk_encryption_key": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"secret_url": {
								Type:     schema.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"key_encryption_key": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key_url": {
								Type:     schema.TypeString,
								Required: true,
							},

							"source_vault_id": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func expandManagedDiskEncryptionSettings(settings map[string]interface{}) *compute.EncryptionSettingsCollection {
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

func flattenManagedDiskEncryptionSettings(encryptionSettings *compute.EncryptionSettingsCollection) []interface{} {
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
