// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	FeatureKeyPrefix = ".appconfig.featureflag"
)

var _ pluginsdk.StateUpgrade = FeatureResourceV0ToV1{}

type FeatureResourceV0ToV1 struct{}

func (FeatureResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/key:name/test/Label/test:label/name
		// new:
		// 	https://appConf1.azconfig.io/kv/.appconfig.featureflag%2Fkey:name%2Ftest?label=test%3Alabel%2Fname
		oldId := rawState["id"].(string)
		fixedId := oldId

		if strings.HasSuffix(fixedId, "/Label/\000") {
			fixedId = strings.TrimSuffix(fixedId, "/Label/\000") + "/Label/%00"
		}

		if strings.HasSuffix(fixedId, "/Label/") {
			fixedId = strings.TrimSuffix(fixedId, "/Label/") + "/Label/%00"
		}

		parsedOldId, err := parse.FeatureId(fixedId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Key Resource %q: %+v", fixedId, err)
		}

		configurationStoreId, err := configurationstores.ParseConfigurationStoreIDInsensitively(parsedOldId.ConfigurationStoreId)
		if err != nil {
			return rawState, fmt.Errorf("parseing Configuration Store ID %q: %+v", configurationStoreId, err)
		}

		domainSuffix, ok := meta.(*clients.Client).Account.Environment.AppConfiguration.DomainSuffix()
		if !ok {
			return rawState, fmt.Errorf("App Configuration is not supported in this Environment")
		}

		configurationStoreEndpoint := fmt.Sprintf("https://%s.%s", configurationStoreId.ConfigurationStoreName, *domainSuffix)
		featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, parsedOldId.Name)
		nestedItemId, err := parse.NewNestedItemID(configurationStoreEndpoint, featureKey, parsedOldId.Label)
		if err != nil {
			return rawState, err
		}

		newId := nestedItemId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func (FeatureResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"description": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"etag": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"label": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"locked": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"percentage_filter_value": {
			Optional: true,
			Type:     pluginsdk.TypeInt,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
		"targeting_filter": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"default_rollout_percentage": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"groups": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"rollout_percentage": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"users": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"timewindow_filter": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"end": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"start": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
	}
}
