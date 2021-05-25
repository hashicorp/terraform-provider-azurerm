package migration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = RegistryV0ToV1{}
var _ pluginsdk.StateUpgrade = RegistryV1ToV2{}

type RegistryV0ToV1 struct{}

func (RegistryV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return registrySchemaForV0AndV1()
}

func (RegistryV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawState["sku"] = "Basic"
		return rawState, nil
	}
}

type RegistryV1ToV2 struct{}

func (RegistryV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return registrySchemaForV0AndV1()
}

func (RegistryV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Basic's been renamed Classic to allow for "ManagedBasic" ¯\_(ツ)_/¯
		rawState["sku"] = "Classic"

		storageAccountId := ""
		if v, ok := rawState["storage_account"]; ok {
			raw := v.(*pluginsdk.Set).List()
			rawVals := raw[0].(map[string]interface{})
			storageAccountName := rawVals["name"].(string)

			client := meta.(*clients.Client).Storage.AccountsClient
			ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, time.Minute*5)
			defer cancel()

			accounts, err := client.ListComplete(ctx)
			if err != nil {
				return rawState, fmt.Errorf("listing storage accounts")
			}

			for accounts.NotDone() {
				account := accounts.Value()
				if strings.EqualFold(*account.Name, storageAccountName) {
					storageAccountId = *account.ID
					break
				}

				if err := accounts.NextWithContext(ctx); err != nil {
					return rawState, fmt.Errorf("retrieving accounts: %+v", err)
				}
			}
		}

		if storageAccountId == "" {
			return rawState, fmt.Errorf("unable to determine storage account ID")
		}

		return rawState, nil
	}
}

func registrySchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// lintignore:S018
		"storage_account": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"access_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_password": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
