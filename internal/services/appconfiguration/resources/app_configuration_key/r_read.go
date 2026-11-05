// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/resources/app_configuration_key/migration"
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
				} else {
					return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
				}
				return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
			}

			model := KeyResourceModel{
				ConfigurationStoreId: configurationStoreId.ID(),
				Key:                  pointer.From(kv.Key),
				ContentType:          pointer.From(kv.ContentType),
				Etag:                 pointer.From(kv.Etag),
				Label:                pointer.From(kv.Label),
				Tags:                 tags.Flatten(kv.Tags),
			}

			if pointer.From(kv.ContentType) != VaultKeyContentType {
				model.Type = KeyTypeKV
				model.Value = pointer.From(kv.Value)
			} else {
				var ref VaultKeyReference
				refBytes := []byte(pointer.From(kv.Value))
				if err := json.Unmarshal(refBytes, &ref); err != nil {
					return fmt.Errorf("while unmarshalling vault reference: %+v", err)
				}

				model.Type = KeyTypeVault
				model.VaultKeyReference = ref.URI
				model.ContentType = VaultKeyContentType
				model.Value = pointer.From(kv.Value)
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

const (
	KeyTypeVault        = "vault"
	KeyTypeKV           = "kv"
	VaultKeyContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
)

type VaultKeyReference struct {
	URI string `json:"uri"`
}

func (k Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 2,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.KeyResourceV0ToV1{},
			1: migration.KeyResourceV1ToV2{},
		},
	}
}
