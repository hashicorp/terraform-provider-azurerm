package extendedlocation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomLocationResource struct{}

type CustomLocationResourceModel struct {
	Name                string      `tfschema:"name"`
	ResourceGroupName   string      `tfschema:"resource_group_name"`
	Location            string      `tfschema:"location"`
	Authentication      []AuthModel `tfschema:"authentication"`
	ClusterExtensionIds []string    `tfschema:"cluster_extension_ids"`
	DisplayName         string      `tfschema:"display_name"`
	HostResourceId      string      `tfschema:"host_resource_id"`
	HostType            string      `tfschema:"host_type"`
	Namespace           string      `tfschema:"namespace"`
}

type AuthModel struct {
	Type  string `tfschema:"type"`
	Value string `tfschema:"value"`
}

func (r CustomLocationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"namespace": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"host_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(customlocations.HostTypeKubernetes),
			}, false),
		},
	}
}

func (r CustomLocationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomLocationResource) ModelObject() interface{} {
	return &CustomLocationResourceModel{}
}

func (r CustomLocationResource) ResourceType() string {
	return "azurerm_extended_custom_location"
}

func (r CustomLocationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CustomLocationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ExtendedLocation.CustomLocationsClient

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customLocationProps := customlocations.CustomLocationProperties{}

			if model.Authentication != nil {
				auth := model.Authentication[0]
				customLocationProps.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
					Type:  &auth.Type,
					Value: &auth.Value,
				}
			}

			if model.ClusterExtensionIds != nil {
				customLocationProps.ClusterExtensionIds = &model.ClusterExtensionIds
			}

			if model.DisplayName != "" {
				customLocationProps.DisplayName = &model.DisplayName
			}

			if model.HostResourceId != "" {
				customLocationProps.HostResourceId = &model.HostResourceId
			}

			if model.HostType != "" {
				hostType := customlocations.HostType(model.HostType)
				customLocationProps.HostType = &hostType
			}

			if model.Namespace != "" {
				customLocationProps.Namespace = &model.Namespace
			}

			props := customlocations.CustomLocation{
				Id:         utils.String(id.ID()),
				Location:   model.Location,
				Name:       utils.String(model.Name),
				Properties: &customLocationProps,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomLocationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state := CustomLocationResourceModel{
					Name:              id.CustomLocationName,
					ResourceGroupName: id.ResourceGroupName,
					Location:          model.Location,
				}

				if props != nil && props.Authentication != nil {
					authType := pointer.From(props.Authentication.Type)
					authValue := pointer.From(props.Authentication.Value)
					state.Authentication = []AuthModel{
						{
							Type:  authType,
							Value: authValue,
						},
					}
				}

				if props.ClusterExtensionIds != nil {
					state.ClusterExtensionIds = *props.ClusterExtensionIds
				}

				if props.DisplayName != nil {
					state.DisplayName = *props.DisplayName
				}

				if props.HostResourceId != nil {
					state.HostResourceId = *props.HostResourceId
				}

				if props.HostType != nil {
					state.HostType = string(*props.HostType)
				}

				if props.Namespace != nil {
					state.Namespace = *props.Namespace
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r CustomLocationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r CustomLocationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state CustomLocationResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			customLocationProps := customlocations.CustomLocationProperties{}
			d := metadata.ResourceData

			if d.HasChanges("authentication") {
				if state.Authentication != nil {
					auth := state.Authentication[0]
					customLocationProps.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
						Type:  &auth.Type,
						Value: &auth.Value,
					}
				}
			}

			if d.HasChange("cluster_extension_ids") {
				customLocationProps.ClusterExtensionIds = &state.ClusterExtensionIds
			}

			if d.HasChange("display_name") {
				customLocationProps.DisplayName = &state.DisplayName
			}

			if d.HasChange("host_resource_id") {
				customLocationProps.HostResourceId = &state.HostResourceId
			}

			if d.HasChange("host_type") {
				hostType := customlocations.HostType(state.HostType)
				customLocationProps.HostType = &hostType
			}

			if d.HasChange("namespace") {
				customLocationProps.Namespace = &state.Namespace
			}

			props := customlocations.PatchableCustomLocations{
				Properties: &customLocationProps,
			}

			if _, err := client.Update(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r CustomLocationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return customlocations.ValidateCustomLocationID
}
