// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/resources/app_configuration_feature/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func (k Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			domainSuffix, ok := metadata.Client.Account.Environment.AppConfiguration.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine AppConfiguration domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			configurationStoreIdRaw, err := metadata.Client.AppConfiguration.ConfigurationStoreIDFromEndpoint(ctx, subscriptionId, nestedItemId.ConfigurationStoreEndpoint, *domainSuffix)
			if err != nil {
				return fmt.Errorf("while retrieving the Resource ID of Configuration Store at Endpoint: %q: %s", nestedItemId.ConfigurationStoreEndpoint, err)
			}
			if configurationStoreIdRaw == nil {
				// if the AppConfiguration is gone then all the data inside it is too
				log.Printf("[DEBUG] Unable to determine the Resource ID for Configuration Store at Endpoint %q - removing from state", nestedItemId.ConfigurationStoreEndpoint)
				return metadata.MarkAsGone(nestedItemId)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(*configurationStoreIdRaw)
			if err != nil {
				return err
			}

			exists, err := metadata.Client.AppConfiguration.Exists(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("while checking Configuration Store %q for feature %q existence: %v", *configurationStoreId, *nestedItemId, err)
			}
			if !exists {
				log.Printf("[DEBUG] Configuration Store %q for feature %q was not found - removing from state", *configurationStoreId, *nestedItemId)
				return metadata.MarkAsGone(nestedItemId)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, "", "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if response.WasNotFound(v.Response) {
						return metadata.MarkAsGone(nestedItemId)
					}
				}
				return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
			}

			var fv FeatureValue
			if err = json.Unmarshal([]byte(pointer.From(kv.Value)), &fv); err != nil {
				return fmt.Errorf("while unmarshalling underlying key's value: %+v", err)
			}

			model := FeatureResourceModel{
				ConfigurationStoreId: configurationStoreId.ID(),
				Description:          fv.Description,
				Enabled:              fv.Enabled,
				Key:                  strings.TrimPrefix(pointer.From(kv.Key), FeatureKeyPrefix+"/"),
				Name:                 fv.ID,
				Label:                pointer.From(kv.Label),
				Tags:                 tags.Flatten(kv.Tags),
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}

			if len(fv.Conditions.ClientFilters.Filters) > 0 {
				for _, f := range fv.Conditions.ClientFilters.Filters {
					switch f := f.(type) {
					case TimewindowFeatureFilter:
						twfp := f
						model.TimewindowFilters = append(model.TimewindowFilters, twfp.Parameters)
					case TargetingFeatureFilter:
						tfp := f
						model.TargetingFilters = append(model.TargetingFilters, tfp.Parameters.Audience)
					case PercentageFeatureFilter:
						pfp := f
						model.PercentageFilter = pfp.Parameters.Value
					case CustomFilter:
						cf := f
						model.CustomFilters = append(model.CustomFilters, cf)
					default:
						return fmt.Errorf("while unmarshaling feature payload: unknown filter type %+v", f)
					}
				}
			}
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

const (
	FeatureKeyContentType = "application/vnd.microsoft.appconfig.ff+json;charset=utf-8"
	FeatureKeyPrefix      = ".appconfig.featureflag"
)

func (k Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.FeatureResourceV0ToV1{},
		},
	}
}
