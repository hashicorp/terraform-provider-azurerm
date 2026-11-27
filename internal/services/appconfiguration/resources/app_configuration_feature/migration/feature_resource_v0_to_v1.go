// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	FeatureKeyPrefix = ".appconfig.featureflag"
)

var _ pluginsdk.StateUpgrade = FeatureResourceV0ToV1{}

type FeatureResourceV0ToV1 struct{}

func (FeatureResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		// old:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationFeature/key:name/test/Label/test:label/name
		// new:
		// 	https://appConf1.azconfig.io/kv/.appconfig.featureflag%2Fkey:name%2Ftest?label=test%3Alabel%2Fname
		oldId := rawState["id"].(string)
		fixedId := oldId

		if before, ok := strings.CutSuffix(fixedId, "/Label/\000"); ok {
			fixedId = before + "/Label/%00"
		}

		if before, ok := strings.CutSuffix(fixedId, "/Label/"); ok {
			fixedId = before + "/Label/%00"
		}

		parsedOldId, err := parse.FeatureId(fixedId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Key Resource %q: %+v", fixedId, err)
		}

		configurationStoreId, err := configurationstores.ParseConfigurationStoreIDInsensitively(parsedOldId.ConfigurationStoreId)
		if err != nil {
			return rawState, fmt.Errorf("parsing Configuration Store ID %q: %+v", configurationStoreId, err)
		}

		domainSuffix, ok := account.From(meta).Environment.AppConfiguration.DomainSuffix()
		if !ok {
			return rawState, errors.New("app configuration is not supported in this Environment")
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"percentage_filter_value": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"targeting_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"default_rollout_percentage": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"groups": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"rollout_percentage": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					}},
				},
				"users": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			}},
		},
		"timewindow_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"end": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"start": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			}},
		},
	}
}
