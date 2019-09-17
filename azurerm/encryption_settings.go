package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
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
	var keyEncryptionKey *compute.KeyVaultAndKeyReference

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

	value := map[string]interface{}{
		"enabled": *encryptionSettings.Enabled,
	}

	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		// at this time we only support a single element
		settings := (*encryptionSettings.EncryptionSettings)[0]
		if key := settings.DiskEncryptionKey; key != nil {
			keys := make(map[string]interface{})

			keys["secret_url"] = *key.SecretURL
			if vault := key.SourceVault; vault != nil {
				keys["source_vault_id"] = *vault.ID
			}

			value["disk_encryption_key"] = []interface{}{keys}
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keys := make(map[string]interface{})

			keys["key_url"] = *key.KeyURL

			if vault := key.SourceVault; vault != nil {
				keys["source_vault_id"] = *vault.ID
			}

			value["key_encryption_key"] = []interface{}{keys}
		}
	}

	return []interface{}{
		value,
	}
}
