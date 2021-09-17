package appconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	FeatureKeyContentType = "application/vnd.microsoft.appconfig.ff+json;charset=utf-8"
	FeatureKeyPrefix      = ".appconfig.featureflag"
)

type FeatureValue struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type FeatureResource struct {
}

var _ sdk.ResourceWithUpdate = FeatureResource{}

type FeatureResourceModel struct {
	ConfigurationStoreId string                 `tfschema:"configuration_store_id"`
	Description          string                 `tfschema:"description"`
	Enabled              bool                   `tfschema:"enabled"`
	Name                 string                 `tfschema:"name"`
	Label                string                 `tfschema:"label"`
	Locked               bool                   `tfschema:"locked"`
	Tags                 map[string]interface{} `tfschema:"tags"`
}

func (k FeatureResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
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
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tags": tags.Schema(),
	}
}

func (k FeatureResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k FeatureResource) ModelObject() interface{} {
	return FeatureResourceModel{}
}

func (k FeatureResource) ResourceType() string {
	return "azurerm_app_configuration_feature"
}

func (k FeatureResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model FeatureResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			appCfgFeatureResourceID := parse.AppConfigurationFeatureId{
				ConfigurationStoreId: model.ConfigurationStoreId,
				Name:                 model.Name,
				Label:                model.Label,
			}

			kv, err := client.GetKeyValues(ctx, model.Name, model.Label, "", "", []string{})
			if err != nil {
				return fmt.Errorf("while checking for feature's %q existence: %+v", model.Name, err)
			}
			keysFound := kv.Values()
			if len(keysFound) > 0 {
				return tf.ImportAsExistsError(k.ResourceType(), appCfgFeatureResourceID.ID())
			}

			err = createOrUpdateFeature(ctx, client, model)

			if appCfgFeatureResourceID.Label == "" {
				// We set an empty label as %00 in the resource ID
				// Otherwise it breaks the ID parsing logic
				appCfgFeatureResourceID.Label = "%00"
			}
			metadata.SetID(appCfgFeatureResourceID)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k FeatureResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.FeatureId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			// We set an empty label as %00 in the ID to make the ID validator happy
			// but in reality the label is just an empty string
			if resourceID.Label == "%00" {
				resourceID.Label = ""
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			res, err := client.GetKeyValues(ctx, resourceID.Name, resourceID.Label, "", "", []string{})
			if err != nil {
				if !utils.ResponseWasNotFound(res.Response().Response) {
					return metadata.MarkAsGone(resourceID)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Name, err)
			}

			if len(res.Values()) > 1 {
				return fmt.Errorf("unexpected API response. More than one value returned for Key/Label pair %s/%s", resourceID.Name, resourceID.Label)
			}

			kv := res.Values()[0]

			var fv FeatureValue
			err = json.Unmarshal([]byte(utils.NormalizeNilableString(kv.Value)), &fv)
			if err != nil {
				return fmt.Errorf("while unmarshalling underlying key's value: %+v", err)
			}

			model := FeatureResourceModel{
				ConfigurationStoreId: resourceID.ConfigurationStoreId,
				Description:          fv.Description,
				Enabled:              fv.Enabled,
				Name:                 fv.ID,
				Label:                utils.NormalizeNilableString(kv.Label),
				Tags:                 tags.Flatten(kv.Tags),
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (k FeatureResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			resourceID, err := parse.FeatureId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			var model FeatureResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") || metadata.ResourceData.HasChange("enabled") || metadata.ResourceData.HasChange("locked") || metadata.ResourceData.HasChange("description") {
				// Remove the lock, if any. We will put it back again if the model says so.
				if _, err = client.DeleteLock(ctx, resourceID.Name, resourceID.Label, "", ""); err != nil {
					return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", resourceID.Name, resourceID.Label, err)
				}
				err = createOrUpdateFeature(ctx, client, model)
				if err != nil {
					return err
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k FeatureResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.FeatureId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}

			if _, err = client.DeleteLock(ctx, resourceID.Name, resourceID.Label, "", ""); err != nil {
				return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", resourceID.Name, resourceID.Label, err)
			}

			_, err = client.DeleteKeyValue(ctx, resourceID.Name, resourceID.Label, "")
			if err != nil {
				return fmt.Errorf("while removing key %q from App Configuration Store %q: %+v", resourceID.Name, resourceID.ConfigurationStoreId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k FeatureResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppConfigurationFeatureID
}

func createOrUpdateFeature(ctx context.Context, client *appconfiguration.BaseClient, model FeatureResourceModel) error {
	featureKey := utils.String(fmt.Sprintf("%s/%s", FeatureKeyPrefix, model.Name))
	entity := appconfiguration.KeyValue{
		Key:         featureKey,
		Label:       utils.String(model.Label),
		Tags:        tags.Expand(model.Tags),
		ContentType: utils.String(FeatureKeyContentType),
		Locked:      utils.Bool(model.Locked),
	}

	value := FeatureValue{
		ID:          model.Name,
		Description: model.Description,
		Enabled:     model.Enabled,
	}
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("while marshalling FeatureValue struct: %+v", err)
	}
	entity.Value = utils.String(string(valueBytes))
	if _, err = client.PutKeyValue(ctx, model.Name, model.Label, &entity, "", ""); err != nil {
		return err
	}

	if model.Locked {
		if _, err = client.PutLock(ctx, model.Name, model.Label, "", ""); err != nil {
			return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Name, model.Label, err)
		}
	} else {
		if _, err = client.DeleteLock(ctx, model.Name, model.Label, "", ""); err != nil {
			return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Name, model.Label, err)
		}
	}

	return nil
}
