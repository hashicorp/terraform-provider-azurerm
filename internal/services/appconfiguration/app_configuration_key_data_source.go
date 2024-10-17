// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

type KeyDataSource struct{}

var _ sdk.DataSource = KeyDataSource{}

func (k KeyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: configurationstores.ValidateConfigurationStoreID,
		},
		"key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "",
		},
	}
}

func (k KeyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"content_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"value": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"vault_key_reference": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"tags": tags.SchemaDataSource(),
	}
}

func (k KeyDataSource) ModelObject() interface{} {
	return &KeyResourceModel{}
}

func (k KeyDataSource) ResourceType() string {
	return "azurerm_app_configuration_key"
}

func (k KeyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			configurationStoreEndpoint, err := metadata.Client.AppConfiguration.EndpointForConfigurationStore(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("retrieving Endpoint for feature %q in %q: %s", model.Key, *configurationStoreId, err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(*configurationStoreEndpoint)
			if err != nil {
				return err
			}

			nestedItemId, err := parse.NewNestedItemID(client.Endpoint, model.Key, model.Label)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValue(ctx, model.Key, model.Label, "", "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("key %s was not found", model.Key)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
			}

			if contentType := utils.NormalizeNilableString(kv.ContentType); contentType != VaultKeyContentType {
				model.Type = KeyTypeKV
				model.ContentType = contentType
				model.Value = utils.NormalizeNilableString(kv.Value)
			} else {
				var ref VaultKeyReference
				refBytes := []byte(utils.NormalizeNilableString(kv.Value))
				err := json.Unmarshal(refBytes, &ref)
				if err != nil {
					return fmt.Errorf("while unmarshalling vault reference: %+v", err)
				}

				model.Type = KeyTypeVault
				model.VaultKeyReference = ref.URI
				model.ContentType = VaultKeyContentType
				model.Value = ref.URI
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}
			model.Etag = utils.NormalizeNilableString(kv.Etag)

			metadata.SetID(nestedItemId)
			return metadata.Encode(&model)
		},
	}
}
