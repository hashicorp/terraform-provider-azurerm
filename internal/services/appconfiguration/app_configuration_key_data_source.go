package appconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyDataSource struct{}

var _ sdk.DataSource = KeyDataSource{}

func (k KeyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
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

			decodedKey, err := url.QueryUnescape(model.Key)
			if err != nil {
				return fmt.Errorf("while decoding key of resource ID: %+v", err)
			}

			id := parse.AppConfigurationKeyId{
				ConfigurationStoreId: model.ConfigurationStoreId,
				Key:                  decodedKey,
				Label:                model.Label,
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValue(ctx, decodedKey, model.Label, "", "", "", []string{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("key %s was not found", decodedKey)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", decodedKey, err)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", decodedKey, err)
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
			if id.Label == "" {
				// We set an empty label as %00 in the resource ID
				// Otherwise it breaks the ID parsing logic
				id.Label = "%00"
			}
			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
