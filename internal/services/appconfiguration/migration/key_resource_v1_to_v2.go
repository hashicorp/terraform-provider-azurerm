// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = KeyResourceV0ToV1{}

type KeyResourceV1ToV2 struct{}

func (KeyResourceV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/key:name/test/Label/test:label/name
		// new:
		// 	https://appConf1.azconfig.io/kv/key:name%2Ftest?label=test%3Alabel%2Fname
		oldId := rawState["id"].(string)
		fixedId := oldId

		// if the ID is like below, it should be bugs for no-label, see https://github.com/hashicorp/terraform-provider-azurerm/issues/20849
		// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/appConfKey1/Label/\000/AppConfigurationKey/appConfKey1/Label/
		if index1, index2 := strings.Index(fixedId, "/AppConfigurationKey/"), strings.LastIndex(fixedId, "/AppConfigurationKey/"); index1 != index2 && fixedId[index2-1] == '\000' && fixedId[index1:index2-1] == fixedId[index2:] {
			fixedId = fixedId[:index2]
		}

		if strings.HasSuffix(fixedId, "/Label/\000") {
			fixedId = strings.TrimSuffix(fixedId, "/Label/\000") + "/Label/%00"
		}

		if strings.HasSuffix(fixedId, "/Label/") {
			fixedId = strings.TrimSuffix(fixedId, "/Label/") + "/Label/%00"
		}

		parsedOldId, err := parse.KeyId(fixedId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Key Resource %q: %+v", fixedId, err)
		}

		configurationStoreId, err := configurationstores.ParseConfigurationStoreIDInsensitively(parsedOldId.ConfigurationStoreId)
		if err != nil {
			return rawState, fmt.Errorf("parseing Configuration Store ID %q: %+v", configurationStoreId, err)
		}

		configurationStoreEndpoint := fmt.Sprintf("https://%s.azconfig.io", configurationStoreId.ConfigurationStoreName)

		nestedItemId, err := parse.NewNestedItemID(configurationStoreEndpoint, parsedOldId.Key, parsedOldId.Label)
		if err != nil {
			return rawState, err
		}

		newId := nestedItemId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func (KeyResourceV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"content_type": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"etag": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"key": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"label": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"locked": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
		"type": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"value": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"vault_key_reference": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
