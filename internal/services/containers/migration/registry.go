// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ pluginsdk.StateUpgrade = RegistryV0ToV1{}
	_ pluginsdk.StateUpgrade = RegistryV1ToV2{}
)

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
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, time.Minute*5)
			defer cancel()

			raw := v.(*pluginsdk.Set).List()
			rawVals := raw[0].(map[string]interface{})
			storageAccountName := rawVals["name"].(string)

			account, err := meta.(*clients.Client).Storage.FindAccount(ctx, subscriptionId, storageAccountName)
			if err != nil {
				return nil, fmt.Errorf("finding Storage Account %q: %+v", storageAccountName, err)
			}

			storageAccountId = account.StorageAccountId.ID()
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

		//lintignore:S018
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
