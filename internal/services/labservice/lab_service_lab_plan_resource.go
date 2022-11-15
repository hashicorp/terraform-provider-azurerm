package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LabServiceLabPlanModel struct {
	Name                     string                     `tfschema:"name"`
	ResourceGroupName        string                     `tfschema:"resource_group_name"`
	Location                 string                     `tfschema:"location"`
	AllowedRegions           []string                   `tfschema:"allowed_regions"`
	DefaultConnectionProfile []DefaultConnectionProfile `tfschema:"default_connection_profile"`
	DefaultNetworkProfile    []DefaultNetworkProfile    `tfschema:"default_network_profile"`
	SharedGalleryId          string                     `tfschema:"shared_gallery_id"`
	Tags                     map[string]string          `tfschema:"tags"`
}

type DefaultConnectionProfile struct {
	ClientRdpAccess labplan.ConnectionType `tfschema:"client_rdp_access"`
	ClientSshAccess labplan.ConnectionType `tfschema:"client_ssh_access"`
	WebRdpAccess    labplan.ConnectionType `tfschema:"web_rdp_access"`
	WebSshAccess    labplan.ConnectionType `tfschema:"web_ssh_access"`
}

type DefaultNetworkProfile struct {
	SubnetId string `tfschema:"subnet_id"`
}

type LabServiceLabPlanResource struct{}

var _ sdk.ResourceWithUpdate = LabServiceLabPlanResource{}

func (r LabServiceLabPlanResource) ResourceType() string {
	return "azurerm_lab_service_lab_plan"
}

func (r LabServiceLabPlanResource) ModelObject() interface{} {
	return &LabServiceLabPlanModel{}
}

func (r LabServiceLabPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return labplan.ValidateLabPlanID
}

func (r LabServiceLabPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LabPlanName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"allowed_regions": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 28,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     location.EnhancedValidate,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"default_connection_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_rdp_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypeNone),
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"client_ssh_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypeNone),
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"web_rdp_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypeNone),
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"web_ssh_access": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypeNone),
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},
				},
			},
		},

		"default_network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: networkValidate.SubnetID,
					},
				},
			},
		},

		"shared_gallery_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: computeValidate.SharedImageGalleryID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r LabServiceLabPlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceLabPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServiceLabPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabService.LabPlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := labplan.NewLabPlanID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := &labplan.LabPlan{
				Location: location.Normalize(model.Location),
				Properties: labplan.LabPlanProperties{
					AllowedRegions: &model.AllowedRegions,
				},
				Tags: &model.Tags,
			}

			defaultConnectionProfile, err := expandDefaultConnectionProfile(model.DefaultConnectionProfile)
			if err != nil {
				return err
			}
			props.Properties.DefaultConnectionProfile = defaultConnectionProfile

			defaultNetworkProfile, err := expandLabPlanDefaultNetworkProfile(model.DefaultNetworkProfile)
			if err != nil {
				return err
			}
			props.Properties.DefaultNetworkProfile = defaultNetworkProfile

			if model.SharedGalleryId != "" {
				props.Properties.SharedGalleryId = &model.SharedGalleryId
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServiceLabPlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabPlanClient

			id, err := labplan.ParseLabPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServiceLabPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("allowed_regions") {
				properties.Properties.AllowedRegions = &model.AllowedRegions
			}

			if metadata.ResourceData.HasChange("default_connection_profile") {
				defaultConnectionProfile, err := expandDefaultConnectionProfile(model.DefaultConnectionProfile)
				if err != nil {
					return err
				}
				properties.Properties.DefaultConnectionProfile = defaultConnectionProfile
			}

			if metadata.ResourceData.HasChange("default_network_profile") {
				defaultNetworkProfile, err := expandLabPlanDefaultNetworkProfile(model.DefaultNetworkProfile)
				if err != nil {
					return err
				}
				properties.Properties.DefaultNetworkProfile = defaultNetworkProfile
			}

			if metadata.ResourceData.HasChange("shared_gallery_id") {
				if model.SharedGalleryId != "" {
					properties.Properties.SharedGalleryId = &model.SharedGalleryId
				} else {
					properties.Properties.SharedGalleryId = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LabServiceLabPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabPlanClient

			id, err := labplan.ParseLabPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LabServiceLabPlanModel{
				Name:              id.LabPlanName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties
			if properties.AllowedRegions != nil {
				state.AllowedRegions = *properties.AllowedRegions
			}

			defaultConnectionProfile, err := flattenDefaultConnectionProfile(properties.DefaultConnectionProfile)
			if err != nil {
				return err
			}

			state.DefaultConnectionProfile = defaultConnectionProfile

			defaultNetworkProfile, err := flattenLabPlanDefaultNetworkProfile(properties.DefaultNetworkProfile)
			if err != nil {
				return err
			}
			state.DefaultNetworkProfile = defaultNetworkProfile

			if properties.SharedGalleryId != nil {
				state.SharedGalleryId = *properties.SharedGalleryId
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceLabPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabPlanClient

			id, err := labplan.ParseLabPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDefaultConnectionProfile(input []DefaultConnectionProfile) (*labplan.ConnectionProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	defaultConnectionProfile := &input[0]
	output := labplan.ConnectionProfile{}

	if defaultConnectionProfile.ClientRdpAccess != "" {
		output.ClientRdpAccess = &defaultConnectionProfile.ClientRdpAccess
	}

	if defaultConnectionProfile.ClientSshAccess != "" {
		output.ClientSshAccess = &defaultConnectionProfile.ClientSshAccess
	}

	if defaultConnectionProfile.WebRdpAccess != "" {
		output.WebRdpAccess = &defaultConnectionProfile.WebRdpAccess
	}

	if defaultConnectionProfile.WebSshAccess != "" {
		output.WebSshAccess = &defaultConnectionProfile.WebSshAccess
	}

	return &output, nil
}

func flattenDefaultConnectionProfile(input *labplan.ConnectionProfile) ([]DefaultConnectionProfile, error) {
	var defaultConnectionProfiles []DefaultConnectionProfile
	if input == nil {
		return defaultConnectionProfiles, nil
	}

	defaultConnectionProfile := DefaultConnectionProfile{}

	if input.ClientRdpAccess != nil {
		defaultConnectionProfile.ClientRdpAccess = *input.ClientRdpAccess
	}

	if input.ClientSshAccess != nil {
		defaultConnectionProfile.ClientSshAccess = *input.ClientSshAccess
	}

	if input.WebRdpAccess != nil {
		defaultConnectionProfile.WebRdpAccess = *input.WebRdpAccess
	}

	if input.WebSshAccess != nil {
		defaultConnectionProfile.WebSshAccess = *input.WebSshAccess
	}

	return append(defaultConnectionProfiles, defaultConnectionProfile), nil
}

func expandLabPlanDefaultNetworkProfile(input []DefaultNetworkProfile) (*labplan.LabPlanNetworkProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	defaultNetworkProfile := &input[0]
	output := labplan.LabPlanNetworkProfile{}

	if defaultNetworkProfile.SubnetId != "" {
		output.SubnetId = &defaultNetworkProfile.SubnetId
	}

	return &output, nil
}

func flattenLabPlanDefaultNetworkProfile(input *labplan.LabPlanNetworkProfile) ([]DefaultNetworkProfile, error) {
	var defaultNetworkProfiles []DefaultNetworkProfile
	if input == nil {
		return defaultNetworkProfiles, nil
	}

	defaultNetworkProfile := DefaultNetworkProfile{}

	if input.SubnetId != nil {
		defaultNetworkProfile.SubnetId = *input.SubnetId
	}

	return append(defaultNetworkProfiles, defaultNetworkProfile), nil
}
