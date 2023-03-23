package appconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

type FeatureResource struct{}

var _ sdk.ResourceWithUpdate = FeatureResource{}

type FeatureResourceModel struct {
	ConfigurationStoreId string                       `tfschema:"configuration_store_id"`
	Description          string                       `tfschema:"description"`
	Enabled              bool                         `tfschema:"enabled"`
	Name                 string                       `tfschema:"name"`
	Label                string                       `tfschema:"label"`
	Locked               bool                         `tfschema:"locked"`
	Tags                 map[string]interface{}       `tfschema:"tags"`
	PercentageFilter     int                          `tfschema:"percentage_filter_value"`
	TimewindowFilters    []TimewindowFilterParameters `tfschema:"timewindow_filter"`
	TargetingFilters     []TargetingFilterAudience    `tfschema:"targeting_filter"`
}

func (k FeatureResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationstores.ValidateConfigurationStoreID,
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
			ValidateFunc: validate.AppConfigurationFeatureName,
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
		"percentage_filter_value": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},
		"targeting_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"default_rollout_percentage": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"rollout_percentage": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 100),
								},
							},
						},
					},
					"users": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"timewindow_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"start": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"end": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
		"tags": tags.Schema(),
	}
}

func (k FeatureResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k FeatureResource) ModelObject() interface{} {
	return &FeatureResourceModel{}
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
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", model.ConfigurationStoreId)
			}

			appCfgFeatureResourceID := parse.AppConfigurationFeatureId{
				ConfigurationStoreId: model.ConfigurationStoreId,
				Name:                 model.Name,
				Label:                model.Label,
			}

			featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, model.Name)

			// from https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration
			// allow up to 15 min for role permission to be done propagated
			metadata.Logger.Infof("[DEBUG] Waiting for App Configuration Key %q read permission to be done propagated", featureKey)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:      []string{"Forbidden"},
				Target:       []string{"Error", "Exists"},
				Refresh:      appConfigurationGetKeyRefreshFunc(ctx, client, featureKey, model.Label),
				PollInterval: 20 * time.Second,
				Timeout:      15 * time.Minute,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q read permission to be propagated: %+v", featureKey, err)
			}

			kv, err := client.GetKeyValue(ctx, featureKey, model.Label, "", "", "", []string{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if !utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("got http status code %d while checking for key's %q existence: %+v", v.Response.StatusCode, featureKey, v.Error())
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", featureKey, err)
				}
			} else if kv.Response.StatusCode == 200 {
				return tf.ImportAsExistsError(k.ResourceType(), appCfgFeatureResourceID.ID())
			}

			err = createOrUpdateFeature(ctx, client, model)
			if err != nil {
				return fmt.Errorf("while creating feature: %+v", err)
			}
			if appCfgFeatureResourceID.Label == "" {
				// We set an empty label as %00 in the resource ID
				// Otherwise it breaks the ID parsing logic
				appCfgFeatureResourceID.Label = "%00"
			}
			metadata.SetID(appCfgFeatureResourceID)
			return nil
		},
		Timeout: 45 * time.Minute,
	}
}

func (k FeatureResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.FeatureId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}
			featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, resourceID.Name)

			// We set an empty label as %00 in the ID to make the ID validator happy
			// but in reality the label is just an empty string
			if resourceID.Label == "%00" {
				resourceID.Label = ""
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}
			if client == nil {
				// if the AppConfiguration is gone then all the data inside it is too
				return metadata.MarkAsGone(resourceID)
			}

			kv, err := client.GetKeyValue(ctx, featureKey, resourceID.Label, "", "", "", []string{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return metadata.MarkAsGone(resourceID)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", featureKey, err)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", featureKey, err)
			}

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

func (k FeatureResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.FeatureId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}
			featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, resourceID.Name)

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", resourceID.ConfigurationStoreId)
			}

			var model FeatureResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") || metadata.ResourceData.HasChange("enabled") || metadata.ResourceData.HasChange("locked") || metadata.ResourceData.HasChange("description") {
				// Remove the lock, if any. We will put it back again if the model says so.
				if _, err = client.DeleteLock(ctx, featureKey, resourceID.Label, "", ""); err != nil {
					return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", resourceID.Name, resourceID.Label, err)
				}
				err = createOrUpdateFeature(ctx, client, model)
				if err != nil {
					return fmt.Errorf("while updating feature: %+v", err)
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
			featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, resourceID.Name)

			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", resourceID.ConfigurationStoreId)
			}
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValues(ctx, featureKey, resourceID.Label, "", "", []string{})
			if err != nil {
				return fmt.Errorf("while checking for feature's %q existence: %+v", resourceID.Name, err)
			}
			keysFound := kv.Values()
			if len(keysFound) == 0 {
				return nil
			}

			if _, err = client.DeleteLock(ctx, featureKey, resourceID.Label, "", ""); err != nil {
				return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", resourceID.Name, resourceID.Label, err)
			}

			_, err = client.DeleteKeyValue(ctx, featureKey, resourceID.Label, "")
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
	featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, model.Name)
	entity := appconfiguration.KeyValue{
		Key:         utils.String(featureKey),
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

	value.Conditions.ClientFilters.Filters = make([]interface{}, 0)
	if model.PercentageFilter > 0 {
		value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, PercentageFeatureFilter{
			Name:       PercentageFilterName,
			Parameters: PercentageFilterParameters{Value: model.PercentageFilter},
		})
	}

	if len(model.TargetingFilters) > 0 {
		for _, tgtf := range model.TargetingFilters {
			value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, TargetingFeatureFilter{
				Name:       TargetingFilterName,
				Parameters: TargetingFilterParameters{Audience: tgtf},
			})
		}
	}

	if len(model.TimewindowFilters) > 0 {
		for _, twf := range model.TimewindowFilters {
			value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, TimewindowFeatureFilter{
				Name:       TimewindowFilterName,
				Parameters: twf,
			})
		}
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("while marshalling FeatureValue struct: %+v", err)
	}
	entity.Value = utils.String(string(valueBytes))
	if _, err = client.PutKeyValue(ctx, featureKey, model.Label, &entity, "", ""); err != nil {
		return err
	}

	if model.Locked {
		if _, err = client.PutLock(ctx, featureKey, model.Label, "", ""); err != nil {
			return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Name, model.Label, err)
		}
	} else {
		if _, err = client.DeleteLock(ctx, featureKey, model.Label, "", ""); err != nil {
			return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Name, model.Label, err)
		}
	}

	return nil
}
