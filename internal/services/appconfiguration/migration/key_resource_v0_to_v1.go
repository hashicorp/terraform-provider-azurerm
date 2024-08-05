// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = KeyResourceV0ToV1{}

type KeyResourceV0ToV1 struct{}

func (KeyResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/key%3Aname%2Ftest/Label/test%3Alabel%2Fname
		// new:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/key:name/test/Label/test:label/name
		oldId := rawState["id"].(string)
		oldKeyNames := regexp.MustCompile(`AppConfigurationKey\/(.+)\/Label`).FindStringSubmatch(oldId)
		if len(oldKeyNames) == 2 {
			decodedName, err := url.QueryUnescape(oldKeyNames[1])
			if err != nil {
				return rawState, err
			}
			oldId = strings.Replace(oldId, oldKeyNames[1], decodedName, 1)
		}
		oldLabelNames := regexp.MustCompile(`AppConfigurationKey\/.+\/Label\/(.+)`).FindStringSubmatch(oldId)
		if len(oldLabelNames) == 2 {
			decodedName, err := url.QueryUnescape(oldLabelNames[1])
			if err != nil {
				return rawState, err
			}
			oldId = strings.Replace(oldId, oldLabelNames[1], decodedName, 1)
		}
		parsedNewId, err := parse.KeyId(oldId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Key Resource %q: %+v", oldId, err)
		}
		rawState["id"] = parsedNewId.ID()

		return rawState, nil
	}
}

func (KeyResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return KeyResourceSchemaForV0AndV1()
}

func KeyResourceSchemaForV0AndV1() map[string]*pluginsdk.Schema {
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
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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
