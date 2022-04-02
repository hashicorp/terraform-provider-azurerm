package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = KeyVaultKeyV0ToV1{}

type KeyVaultKeyV0ToV1 struct{}

func (KeyVaultKeyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_vault_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
		},

		"key_opts": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"not_before_date": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"expiration_date": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versionless_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"n": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"e": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"x": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"y": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_key_pem": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_key_openssh": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (KeyVaultKeyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Old:
		//  https://example-keyvault.vault.azure.net/keys/example/fdf067c93bbb4b22bff4d8b7a9a56217
		// New:
		//  /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1/keys/key1
		oldId := rawState["id"].(string)
		id, err := parse.ParseNestedItemID(oldId)
		if err != nil {
			return nil, err
		}
		client := meta.(*clients.Client).KeyVault
		resourcesClient := meta.(*clients.Client).Resource
		keyVaultIdRaw, err := client.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
		if err != nil {
			return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
		}
		if keyVaultIdRaw == nil {
			return nil, fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
		}

		keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
		if err != nil {
			return nil, err
		}

		name := rawState["name"].(string)
		version := rawState["version"].(string)
		newId := parse.NewKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, name, version).ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
