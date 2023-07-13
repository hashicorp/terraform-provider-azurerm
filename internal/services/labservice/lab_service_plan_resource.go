// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServicePlanModel struct {
	Name                   string                `tfschema:"name"`
	ResourceGroupName      string                `tfschema:"resource_group_name"`
	Location               string                `tfschema:"location"`
	AllowedRegions         []string              `tfschema:"allowed_regions"`
	DefaultAutoShutdown    []DefaultAutoShutdown `tfschema:"default_auto_shutdown"`
	DefaultConnection      []DefaultConnection   `tfschema:"default_connection"`
	DefaultNetworkSubnetId string                `tfschema:"default_network_subnet_id"`
	SharedGalleryId        string                `tfschema:"shared_gallery_id"`
	Support                []Support             `tfschema:"support"`
	Tags                   map[string]string     `tfschema:"tags"`
}

type DefaultAutoShutdown struct {
	DisconnectDelay string                     `tfschema:"disconnect_delay"`
	IdleDelay       string                     `tfschema:"idle_delay"`
	NoConnectDelay  string                     `tfschema:"no_connect_delay"`
	ShutdownOnIdle  labplan.ShutdownOnIdleMode `tfschema:"shutdown_on_idle"`
}

type DefaultConnection struct {
	ClientRdpAccess labplan.ConnectionType `tfschema:"client_rdp_access"`
	ClientSshAccess labplan.ConnectionType `tfschema:"client_ssh_access"`
	WebRdpAccess    labplan.ConnectionType `tfschema:"web_rdp_access"`
	WebSshAccess    labplan.ConnectionType `tfschema:"web_ssh_access"`
}

type Support struct {
	Email        string `tfschema:"email"`
	Instructions string `tfschema:"instructions"`
	Phone        string `tfschema:"phone"`
	Url          string `tfschema:"url"`
}

type LabServicePlanResource struct{}

var _ sdk.ResourceWithUpdate = LabServicePlanResource{}

func (r LabServicePlanResource) ResourceType() string {
	return "azurerm_lab_service_plan"
}

func (r LabServicePlanResource) ModelObject() interface{} {
	return &LabServicePlanModel{}
}

func (r LabServicePlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return labplan.ValidateLabPlanID
}

func (r LabServicePlanResource) Arguments() map[string]*pluginsdk.Schema {
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

		"default_auto_shutdown": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disconnect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"idle_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"no_connect_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azValidate.ISO8601Duration,
					},

					"shutdown_on_idle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ShutdownOnIdleModeUserAbsence),
							string(labplan.ShutdownOnIdleModeLowUsage),
						}, false),
					},
				},
			},
		},

		"default_connection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_rdp_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"client_ssh_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"web_rdp_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},

					"web_ssh_access": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(labplan.ConnectionTypePrivate),
							string(labplan.ConnectionTypePublic),
						}, false),
					},
				},
			},
		},

		"default_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"shared_gallery_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: galleries.ValidateGalleryID,
		},

		"support": {
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

func (r LabServicePlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServicePlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServicePlanModel
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
					AllowedRegions:             normalizeAllowedRegions(model.AllowedRegions),
					DefaultAutoShutdownProfile: expandDefaultAutoShutdown(model.DefaultAutoShutdown),
					DefaultConnectionProfile:   expandDefaultConnection(model.DefaultConnection),
					DefaultNetworkProfile:      expandDefaultNetwork(model.DefaultNetworkSubnetId),
					SupportInfo:                expandSupportInfo(model.Support),
				},
				Tags: &model.Tags,
			}

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

func (r LabServicePlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabPlanClient

			id, err := labplan.ParseLabPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServicePlanModel
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
				properties.Properties.AllowedRegions = normalizeAllowedRegions(model.AllowedRegions)
			}

			if metadata.ResourceData.HasChange("default_auto_shutdown") {
				properties.Properties.DefaultAutoShutdownProfile = expandDefaultAutoShutdown(model.DefaultAutoShutdown)
			}

			if metadata.ResourceData.HasChange("default_connection") {
				properties.Properties.DefaultConnectionProfile = expandDefaultConnection(model.DefaultConnection)
			}

			if metadata.ResourceData.HasChange("default_network_subnet_id") {
				properties.Properties.DefaultNetworkProfile = expandDefaultNetwork(model.DefaultNetworkSubnetId)
			}

			if metadata.ResourceData.HasChange("support") {
				properties.Properties.SupportInfo = expandSupportInfo(model.Support)
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

func (r LabServicePlanResource) Read() sdk.ResourceFunc {
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

			state := LabServicePlanModel{
				Name:              id.LabPlanName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties
			state.DefaultNetworkSubnetId = flattenDefaultNetwork(properties.DefaultNetworkProfile)
			state.DefaultAutoShutdown = flattenDefaultAutoShutdown(properties.DefaultAutoShutdownProfile)
			state.DefaultConnection = flattenDefaultConnection(properties.DefaultConnectionProfile)
			state.Support = flattenSupportInfo(properties.SupportInfo)

			if properties.AllowedRegions != nil {
				state.AllowedRegions = *normalizeAllowedRegions(*properties.AllowedRegions)
			}

			if galleryId := properties.SharedGalleryId; galleryId != nil {
				state.SharedGalleryId = *galleryId
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServicePlanResource) Delete() sdk.ResourceFunc {
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

func expandDefaultAutoShutdown(input []DefaultAutoShutdown) *labplan.AutoShutdownProfile {
	if len(input) == 0 {
		return nil
	}

	defaultAutoShutdownProfile := &input[0]
	result := labplan.AutoShutdownProfile{}

	shutdownOnDisconnectEnabled := labplan.EnableStateDisabled
	if defaultAutoShutdownProfile.DisconnectDelay != "" {
		result.DisconnectDelay = &defaultAutoShutdownProfile.DisconnectDelay
		shutdownOnDisconnectEnabled = labplan.EnableStateEnabled
	}
	result.ShutdownOnDisconnect = &shutdownOnDisconnectEnabled

	if defaultAutoShutdownProfile.IdleDelay != "" {
		result.IdleDelay = &defaultAutoShutdownProfile.IdleDelay
	}

	shutdownWhenNotConnectedEnabled := labplan.EnableStateDisabled
	if defaultAutoShutdownProfile.NoConnectDelay != "" {
		result.NoConnectDelay = &defaultAutoShutdownProfile.NoConnectDelay
		shutdownWhenNotConnectedEnabled = labplan.EnableStateEnabled
	}
	result.ShutdownWhenNotConnected = &shutdownWhenNotConnectedEnabled

	shutdownOnIdle := labplan.ShutdownOnIdleModeNone
	if defaultAutoShutdownProfile.ShutdownOnIdle != "" {
		shutdownOnIdle = defaultAutoShutdownProfile.ShutdownOnIdle
	}
	result.ShutdownOnIdle = &shutdownOnIdle

	return &result
}

func flattenDefaultAutoShutdown(input *labplan.AutoShutdownProfile) []DefaultAutoShutdown {
	var defaultAutoShutdownProfiles []DefaultAutoShutdown
	if input == nil {
		return defaultAutoShutdownProfiles
	}

	defaultAutoShutdownProfile := DefaultAutoShutdown{}

	if input.DisconnectDelay != nil {
		defaultAutoShutdownProfile.DisconnectDelay = *input.DisconnectDelay
	}

	if input.IdleDelay != nil {
		defaultAutoShutdownProfile.IdleDelay = *input.IdleDelay
	}

	if input.NoConnectDelay != nil {
		defaultAutoShutdownProfile.NoConnectDelay = *input.NoConnectDelay
	}

	if shutdownOnIdle := input.ShutdownOnIdle; shutdownOnIdle != nil && *shutdownOnIdle != labplan.ShutdownOnIdleModeNone {
		defaultAutoShutdownProfile.ShutdownOnIdle = *shutdownOnIdle
	}

	return append(defaultAutoShutdownProfiles, defaultAutoShutdownProfile)
}

func expandDefaultConnection(input []DefaultConnection) *labplan.ConnectionProfile {
	if len(input) == 0 {
		return nil
	}

	defaultConnectionProfile := &input[0]
	result := labplan.ConnectionProfile{}

	clientRdpAccess := labplan.ConnectionTypeNone
	if defaultConnectionProfile.ClientRdpAccess != "" {
		clientRdpAccess = defaultConnectionProfile.ClientRdpAccess
	}
	result.ClientRdpAccess = &clientRdpAccess

	clientSshAccess := labplan.ConnectionTypeNone
	if defaultConnectionProfile.ClientSshAccess != "" {
		clientSshAccess = defaultConnectionProfile.ClientSshAccess
	}
	result.ClientSshAccess = &clientSshAccess

	webRdpAccess := labplan.ConnectionTypeNone
	if defaultConnectionProfile.WebRdpAccess != "" {
		webRdpAccess = defaultConnectionProfile.WebRdpAccess
	}
	result.WebRdpAccess = &webRdpAccess

	webSshAccess := labplan.ConnectionTypeNone
	if defaultConnectionProfile.WebSshAccess != "" {
		webSshAccess = defaultConnectionProfile.WebSshAccess
	}
	result.WebSshAccess = &webSshAccess

	return &result
}

func flattenDefaultConnection(input *labplan.ConnectionProfile) []DefaultConnection {
	var defaultConnectionProfiles []DefaultConnection
	if input == nil {
		return defaultConnectionProfiles
	}

	defaultConnectionProfile := DefaultConnection{}

	if clientRdpAccess := input.ClientRdpAccess; clientRdpAccess != nil && *clientRdpAccess != labplan.ConnectionTypeNone {
		defaultConnectionProfile.ClientRdpAccess = *clientRdpAccess
	}

	if clientSshAccess := input.ClientSshAccess; clientSshAccess != nil && *clientSshAccess != labplan.ConnectionTypeNone {
		defaultConnectionProfile.ClientSshAccess = *clientSshAccess
	}

	if webRdpAccess := input.WebRdpAccess; webRdpAccess != nil && *webRdpAccess != labplan.ConnectionTypeNone {
		defaultConnectionProfile.WebRdpAccess = *webRdpAccess
	}

	if webSshAccess := input.WebSshAccess; webSshAccess != nil && *webSshAccess != labplan.ConnectionTypeNone {
		defaultConnectionProfile.WebSshAccess = *webSshAccess
	}

	return append(defaultConnectionProfiles, defaultConnectionProfile)
}

func expandDefaultNetwork(input string) *labplan.LabPlanNetworkProfile {
	if input == "" {
		return nil
	}

	result := labplan.LabPlanNetworkProfile{
		SubnetId: utils.String(input),
	}

	return &result
}

func flattenDefaultNetwork(input *labplan.LabPlanNetworkProfile) string {
	var defaultNetworkSubnetId string
	if input == nil {
		return defaultNetworkSubnetId
	}

	if input.SubnetId != nil {
		defaultNetworkSubnetId = *input.SubnetId
	}

	return defaultNetworkSubnetId
}

func expandSupportInfo(input []Support) *labplan.SupportInfo {
	if len(input) == 0 {
		return nil
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

	return &result
}

func flattenSupportInfo(input *labplan.SupportInfo) []Support {
	var supportInfos []Support
	if input == nil {
		return supportInfos
	}

	supportInfo := Support{}

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

	return append(supportInfos, supportInfo)
}

func normalizeAllowedRegions(input []string) *[]string {
	regions := make([]string, 0)

	for _, v := range input {
		region := location.Normalize(v)
		regions = append(regions, region)
	}

	return &regions
}
