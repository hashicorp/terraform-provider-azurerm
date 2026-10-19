// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func (k Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			metadata.Client.AppConfiguration.AddToCache(*configurationStoreId, nestedItemId.ConfigurationStoreEndpoint)

			// lintignore:R019 // deliberate subset: only the fields feeding the KeyValue update entity
			if metadata.ResourceData.HasChanges("value", "content_type", "tags", "type", "vault_key_reference") {
				entity := appconfiguration.KeyValue{
					Key:   pointer.To(model.Key),
					Label: pointer.To(model.Label),
					Tags:  tags.Expand(model.Tags),
				}

				switch model.Type {
				case KeyTypeKV:
					entity.ContentType = pointer.To(model.ContentType)
					entity.Value = pointer.To(model.Value)
				case KeyTypeVault:
					entity.ContentType = pointer.To(VaultKeyContentType)
					ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
					if err != nil {
						return fmt.Errorf("while encoding vault key reference: %+v", err)
					}
					entity.Value = pointer.To(string(ref))
				}
				if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
					return fmt.Errorf("while updating key/label pair %s/%s: %+v", model.Key, model.Label, err)
				}
			}

			if metadata.ResourceData.HasChange("locked") {
				if model.Locked {
					if _, err = client.PutLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				} else {
					if _, err = client.DeleteLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
