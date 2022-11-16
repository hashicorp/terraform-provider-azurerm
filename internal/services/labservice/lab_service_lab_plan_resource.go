package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LabServiceLabPlanModel struct {
	Name                       string                       `tfschema:"name"`
	ResourceGroupName          string                       `tfschema:"resource_group_name"`
	Location                   string                       `tfschema:"location"`
	AllowedRegions             []string                     `tfschema:"allowed_regions"`
	DefaultAutoShutdownProfile []DefaultAutoShutdownProfile `tfschema:"default_auto_shutdown_profile"`
	DefaultConnectionProfile   []DefaultConnectionProfile   `tfschema:"default_connection_profile"`
	DefaultNetworkProfile      []DefaultNetworkProfile      `tfschema:"default_network_profile"`
	SharedGalleryId            string                       `tfschema:"shared_gallery_id"`
	SupportInfo                []SupportInfo                `tfschema:"support_info"`
	Tags                       map[string]string            `tfschema:"tags"`
}

type DefaultAutoShutdownProfile struct {
	DisconnectDelay                 string                     `tfschema:"disconnect_delay"`
	IdleDelay                       string                     `tfschema:"idle_delay"`
	NoConnectDelay                  string                     `tfschema:"no_connect_delay"`
	ShutdownOnDisconnectEnabled     bool                       `tfschema:"shutdown_on_disconnect_enabled"`
	ShutdownOnIdle                  labplan.ShutdownOnIdleMode `tfschema:"shutdown_on_idle"`
	ShutdownEnabledWhenNotConnected bool                       `tfschema:"shutdown_enabled_when_not_connected"`
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

type SupportInfo struct {
	Email        string `tfschema:"email"`
	Instructions string `tfschema:"instructions"`
	Phone        string `tfschema:"phone"`
	Url          string `tfschema:"url"`
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

		"default_auto_shutdown_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disconnect_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"idle_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"no_connect_delay": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"shutdown_on_disconnect_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"shutdown_on_idle": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ShutdownOnIdleModeUserAbsence),
							string(labplan.ShutdownOnIdleModeLowUsage),
							string(labplan.ShutdownOnIdleModeNone),
						}, false),
					},

					"shutdown_enabled_when_not_connected": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
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

		"support_info": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SupportInfoEmail,
					},

					"instructions": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SupportInfoPhone,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
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

			defaultAutoShutdownProfile, err := expandDefaultAutoShutdownProfile(model.DefaultAutoShutdownProfile)
			if err != nil {
				return err
			}

			props.Properties.DefaultAutoShutdownProfile = defaultAutoShutdownProfile

			defaultConnectionProfile, err := expandDefaultConnectionProfile(model.DefaultConnectionProfile)
			if err != nil {
				return err
			}
			props.Properties.DefaultConnectionProfile = defaultConnectionProfile

			defaultNetworkProfile, err := expandDefaultNetworkProfile(model.DefaultNetworkProfile)
			if err != nil {
				return err
			}
			props.Properties.DefaultNetworkProfile = defaultNetworkProfile

			supportInfo, err := expandSupportInfo(model.SupportInfo)
			if err != nil {
				return err
			}
			props.Properties.SupportInfo = supportInfo

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

			if metadata.ResourceData.HasChange("default_auto_shutdown_profile") {
				defaultAutoShutdownProfile, err := expandDefaultAutoShutdownProfile(model.DefaultAutoShutdownProfile)
				if err != nil {
					return err
				}
				properties.Properties.DefaultAutoShutdownProfile = defaultAutoShutdownProfile
			}

			if metadata.ResourceData.HasChange("default_connection_profile") {
				defaultConnectionProfile, err := expandDefaultConnectionProfile(model.DefaultConnectionProfile)
				if err != nil {
					return err
				}
				properties.Properties.DefaultConnectionProfile = defaultConnectionProfile
			}

			if metadata.ResourceData.HasChange("default_network_profile") {
				defaultNetworkProfile, err := expandDefaultNetworkProfile(model.DefaultNetworkProfile)
				if err != nil {
					return err
				}
				properties.Properties.DefaultNetworkProfile = defaultNetworkProfile
			}

			if metadata.ResourceData.HasChange("support_info") {
				supportInfo, err := expandSupportInfo(model.SupportInfo)
				if err != nil {
					return err
				}
				properties.Properties.SupportInfo = supportInfo
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

			defaultAutoShutdownProfile, err := flattenDefaultAutoShutdownProfile(properties.DefaultAutoShutdownProfile)
			if err != nil {
				return err
			}
			state.DefaultAutoShutdownProfile = defaultAutoShutdownProfile

			defaultConnectionProfile, err := flattenDefaultConnectionProfile(properties.DefaultConnectionProfile)
			if err != nil {
				return err
			}
			state.DefaultConnectionProfile = defaultConnectionProfile

			defaultNetworkProfile, err := flattenDefaultNetworkProfile(properties.DefaultNetworkProfile)
			if err != nil {
				return err
			}
			state.DefaultNetworkProfile = defaultNetworkProfile

			supportInfo, err := flattenSupportInfo(properties.SupportInfo)
			if err != nil {
				return err
			}
			state.SupportInfo = supportInfo

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

func expandDefaultAutoShutdownProfile(input []DefaultAutoShutdownProfile) (*labplan.AutoShutdownProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	defaultAutoShutdownProfile := &input[0]
	result := labplan.AutoShutdownProfile{}

	if defaultAutoShutdownProfile.DisconnectDelay != "" {
		result.DisconnectDelay = &defaultAutoShutdownProfile.DisconnectDelay
	}

	if defaultAutoShutdownProfile.IdleDelay != "" {
		result.IdleDelay = &defaultAutoShutdownProfile.IdleDelay
	}

	if defaultAutoShutdownProfile.NoConnectDelay != "" {
		result.NoConnectDelay = &defaultAutoShutdownProfile.NoConnectDelay
	}

	shutdownOnDisconnectEnabled := labplan.EnableStateEnabled
	if !defaultAutoShutdownProfile.ShutdownOnDisconnectEnabled {
		shutdownOnDisconnectEnabled = labplan.EnableStateDisabled
	}
	result.ShutdownOnDisconnect = &shutdownOnDisconnectEnabled

	if defaultAutoShutdownProfile.ShutdownOnIdle != "" {
		result.ShutdownOnIdle = &defaultAutoShutdownProfile.ShutdownOnIdle
	}

	shutdownEnabledWhenNotConnected := labplan.EnableStateEnabled
	if !defaultAutoShutdownProfile.ShutdownEnabledWhenNotConnected {
		shutdownEnabledWhenNotConnected = labplan.EnableStateDisabled
	}
	result.ShutdownWhenNotConnected = &shutdownEnabledWhenNotConnected

	return &result, nil
}

func flattenDefaultAutoShutdownProfile(input *labplan.AutoShutdownProfile) ([]DefaultAutoShutdownProfile, error) {
	var defaultAutoShutdownProfiles []DefaultAutoShutdownProfile
	if input == nil {
		return defaultAutoShutdownProfiles, nil
	}

	defaultAutoShutdownProfile := DefaultAutoShutdownProfile{}

	if input.DisconnectDelay != nil {
		defaultAutoShutdownProfile.DisconnectDelay = *input.DisconnectDelay
	}

	if input.IdleDelay != nil {
		defaultAutoShutdownProfile.IdleDelay = *input.IdleDelay
	}

	if input.NoConnectDelay != nil {
		defaultAutoShutdownProfile.NoConnectDelay = *input.NoConnectDelay
	}

	if input.ShutdownOnDisconnect != nil {
		defaultAutoShutdownProfile.ShutdownOnDisconnectEnabled = *input.ShutdownOnDisconnect == labplan.EnableStateEnabled
	}

	if input.ShutdownOnIdle != nil {
		defaultAutoShutdownProfile.ShutdownOnIdle = *input.ShutdownOnIdle
	}

	if input.ShutdownWhenNotConnected != nil {
		defaultAutoShutdownProfile.ShutdownEnabledWhenNotConnected = *input.ShutdownWhenNotConnected == labplan.EnableStateEnabled
	}

	return append(defaultAutoShutdownProfiles, defaultAutoShutdownProfile), nil
}

func expandDefaultConnectionProfile(input []DefaultConnectionProfile) (*labplan.ConnectionProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	defaultConnectionProfile := &input[0]
	result := labplan.ConnectionProfile{}

	if defaultConnectionProfile.ClientRdpAccess != "" {
		result.ClientRdpAccess = &defaultConnectionProfile.ClientRdpAccess
	}

	if defaultConnectionProfile.ClientSshAccess != "" {
		result.ClientSshAccess = &defaultConnectionProfile.ClientSshAccess
	}

	if defaultConnectionProfile.WebRdpAccess != "" {
		result.WebRdpAccess = &defaultConnectionProfile.WebRdpAccess
	}

	if defaultConnectionProfile.WebSshAccess != "" {
		result.WebSshAccess = &defaultConnectionProfile.WebSshAccess
	}

	return &result, nil
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

func expandDefaultNetworkProfile(input []DefaultNetworkProfile) (*labplan.LabPlanNetworkProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	defaultNetworkProfile := &input[0]
	result := labplan.LabPlanNetworkProfile{}

	if defaultNetworkProfile.SubnetId != "" {
		result.SubnetId = &defaultNetworkProfile.SubnetId
	}

	return &result, nil
}

func flattenDefaultNetworkProfile(input *labplan.LabPlanNetworkProfile) ([]DefaultNetworkProfile, error) {
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

func expandSupportInfo(input []SupportInfo) (*labplan.SupportInfo, error) {
	if len(input) == 0 {
		return nil, nil
	}

	supportInfo := &input[0]
	result := labplan.SupportInfo{}

	if supportInfo.Email != "" {
		result.Email = &supportInfo.Email
	}

	if supportInfo.Instructions != "" {
		result.Instructions = &supportInfo.Instructions
	}

	if supportInfo.Phone != "" {
		result.Phone = &supportInfo.Phone
	}

	if supportInfo.Url != "" {
		result.Url = &supportInfo.Url
	}

	return &result, nil
}

func flattenSupportInfo(input *labplan.SupportInfo) ([]SupportInfo, error) {
	var supportInfos []SupportInfo
	if input == nil {
		return supportInfos, nil
	}

	supportInfo := SupportInfo{}

	if input.Email != nil {
		supportInfo.Email = *input.Email
	}

	if input.Instructions != nil {
		supportInfo.Instructions = *input.Instructions
	}

	if input.Phone != nil {
		supportInfo.Phone = *input.Phone
	}

	if input.Url != nil {
		supportInfo.Url = *input.Url
	}

	return append(supportInfos, supportInfo), nil
}
