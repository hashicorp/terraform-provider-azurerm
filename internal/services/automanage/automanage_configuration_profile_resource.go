package automanage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type AutoManageConfigurationProfileModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Configuration     string            `tfschema:"configuration_json"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type AutoManageConfigurationProfileResource struct{}

var _ sdk.ResourceWithUpdate = AutoManageConfigurationProfileResource{}

func (r AutoManageConfigurationProfileResource) ResourceType() string {
	return "azurerm_automanage_configuration_profile"
}

func (r AutoManageConfigurationProfileResource) ModelObject() interface{} {
	return &AutoManageConfigurationProfileModel{}
}

func (r AutoManageConfigurationProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AutomanageConfigurationProfileID
}

func (r AutoManageConfigurationProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"configuration_json": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AutoManageConfigurationProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutoManageConfigurationProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AutoManageConfigurationProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Automanage.ConfigurationProfileClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewAutomanageConfigurationProfileID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := automanage.ConfigurationProfile{
				Location:   utils.String(location.Normalize(model.Location)),
				Properties: &automanage.ConfigurationProfileProperties{},
				Tags:       tags.FromTypedObject(model.Tags),
			}

			if model.Configuration != "" {
				var configurationValue interface{}
				err = json.Unmarshal([]byte(model.Configuration), &configurationValue)
				if err != nil {
					return err
				}
				properties.Properties.Configuration = &configurationValue
			}

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.ResourceGroup, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutoManageConfigurationProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileClient

			id, err := parse.AutomanageConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AutoManageConfigurationProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("configuration_json") {
				var configurationValue interface{}
				err := json.Unmarshal([]byte(model.Configuration), &configurationValue)
				if err != nil {
					return err
				}

				resp.Properties.Configuration = &configurationValue
			}

			resp.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				resp.Tags = tags.FromTypedObject(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.ResourceGroup, resp); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AutoManageConfigurationProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileClient

			id, err := parse.AutomanageConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AutoManageConfigurationProfileModel{
				Name:              id.ConfigurationProfileName,
				ResourceGroupName: id.ResourceGroup,
				Location:          location.NormalizeNilable(resp.Location),
			}

			if properties := resp.Properties; properties != nil {
				if properties.Configuration != nil {
					configurationValue, err := json.Marshal(properties.Configuration)
					if err != nil {
						return err
					}

					state.Configuration = string(configurationValue)
				}
			}
			if resp.Tags != nil {
				state.Tags = tags.ToTypedObject(resp.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutoManageConfigurationProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileClient

			id, err := parse.AutomanageConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
