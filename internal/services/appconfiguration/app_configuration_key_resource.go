package appconfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyResource struct {
}

var _ sdk.Resource = KeyResource{}

type KeyResourceModel struct {
	ConfigurationStoreId string                 `tfschema:"configuration_store_id"`
	Key                  string                 `tfschema:"key"`
	ContentType          string                 `tfschema:"content_type"`
	Label                string                 `tfschema:"label"`
	Value                string                 `tfschema:"value"`
	Locked               bool                   `tfschema:"locked"`
	Tags                 map[string]interface{} `tfschema:"tags"`
}

func (k KeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"value": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tags": tags.Schema(),
	}
}

func (k KeyResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (k KeyResource) ModelObject() interface{} {
	return KeyResourceModel{}
}

func (k KeyResource) ResourceType() string {
	return "azurerm_app_configuration_key"
}

func (k KeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			appCfgKeyResourceID := parse.AppConfigurationKeyId{
				ConfigurationStoreId: model.ConfigurationStoreId,
				Key:                  model.Key,
				Label:                model.Label,
			}

			kv, err := client.GetKeyValues(ctx, model.Key, model.Label, "", "", []string{})
			if err != nil {
				return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
			}
			keysFound := kv.Values()
			if len(keysFound) > 0 {
				return tf.ImportAsExistsError(k.ResourceType(), appCfgKeyResourceID.ID())
			}

			entity := appconfiguration.KeyValue{
				Key:         utils.String(model.Key),
				Label:       utils.String(model.Label),
				ContentType: utils.String(model.ContentType),
				Value:       utils.String(model.Value),
				Tags:        tags.Expand(model.Tags),
			}

			if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
				return err
			}

			if model.Locked {
				_, err = client.PutLock(ctx, model.Key, model.Label, "", "")
				if err != nil {
					return fmt.Errorf("while locking key/label pair %q/%q: %+v", model.Key, model.Label, err)
				}
			}

			metadata.SetID(appCfgKeyResourceID)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.AppConfigurationKeyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			res, err := client.GetKeyValues(ctx, resourceID.Key, resourceID.Label, "", "", []string{})
			if err != nil {
				if !utils.ResponseWasNotFound(res.Response().Response) {
					return metadata.MarkAsGone(resourceID)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Key, err)
			}

			if len(res.Values()) > 1 {
				return fmt.Errorf("unexpected API response. More than one value returned for Key/Label pair %s/%s", resourceID.Key, resourceID.Label)
			}

			kv := res.Values()[0]
			var locked bool
			if kv.Locked != nil {
				locked = *kv.Locked
			}

			model := KeyResourceModel{
				ConfigurationStoreId: resourceID.ConfigurationStoreId,
				Key:                  utils.NormalizeNilableString(kv.Key),
				ContentType:          utils.NormalizeNilableString(kv.ContentType),
				Label:                utils.NormalizeNilableString(kv.Label),
				Value:                utils.NormalizeNilableString(kv.Value),
				Locked:               locked,
				Tags:                 tags.Flatten(kv.Tags),
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (k KeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			resourceID, err := parse.AppConfigurationKeyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("value") || metadata.ResourceData.HasChange("content_type") || metadata.ResourceData.HasChange("tags") {

				entity := appconfiguration.KeyValue{
					Key:         utils.String(model.Key),
					Label:       utils.String(model.Label),
					ContentType: utils.String(model.ContentType),
					Value:       utils.String(model.Value),
					Tags:        tags.Expand(model.Tags),
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

func (k KeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.AppConfigurationKeyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			_, err = client.DeleteKeyValue(ctx, resourceID.Key, resourceID.Label, "")
			if err != nil {
				return fmt.Errorf("while removing key %q from App Configuration Store %q: %+v", resourceID.Key, resourceID.ConfigurationStoreId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppConfigurationKeyID
}
