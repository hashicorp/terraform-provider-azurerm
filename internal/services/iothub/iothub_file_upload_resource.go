// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

type IotHubFileUploadResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotHubFileUploadResource{}
)

type IotHubFileUploadResourceModel struct {
	AuthenticationType   string `tfschema:"authentication_type"`
	ConnectionString     string `tfschema:"connection_string"`
	ContainerName        string `tfschema:"container_name"`
	DefaultTTL           string `tfschema:"default_ttl"`
	IdentityId           string `tfschema:"identity_id"`
	IotHubId             string `tfschema:"iothub_id"`
	LockDuration         string `tfschema:"lock_duration"`
	MaxDeliveryCount     int64  `tfschema:"max_delivery_count"`
	NotificationsEnabled bool   `tfschema:"notifications_enabled"`
	SasTTL               string `tfschema:"sas_ttl"`
}

func (r IotHubFileUploadResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"connection_string": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: fileUploadConnectionStringDiffSuppress,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageContainerName,
		},

		"authentication_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.AuthenticationTypeKeyBased),
				string(devices.AuthenticationTypeIdentityBased),
			}, false),
		},

		"default_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H",
			ValidateFunc: azValidate.ISO8601Duration,
		},

		"identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"lock_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1M",
			ValidateFunc: azValidate.ISO8601Duration,
		},

		"max_delivery_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"notifications_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"sas_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H",
			ValidateFunc: azValidate.ISO8601Duration,
		},
	}
}

func (r IotHubFileUploadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotHubFileUploadResource) ResourceType() string {
	return "azurerm_iothub_file_upload"
}

func (r IotHubFileUploadResource) ModelObject() interface{} {
	return &IotHubFileUploadResourceModel{}
}

func (r IotHubFileUploadResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.IotHubID
}

func (r IotHubFileUploadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			var state IotHubFileUploadResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.IotHubID(state.IotHubId)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			iotHub, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(iotHub.Response) {
					return fmt.Errorf("%q was not found", id)
				}

				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iotHub.Properties != nil && iotHub.Properties.MessagingEndpoints != nil {
				if storageEndpoint, ok := iotHub.Properties.StorageEndpoints["$default"]; ok {
					if storageEndpoint.ConnectionString != nil && *storageEndpoint.ConnectionString != "" && storageEndpoint.ContainerName != nil && *storageEndpoint.ContainerName != "" {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			messagingEndpointProperties := make(map[string]*devices.MessagingEndpointProperties)
			storageEndpointProperties := make(map[string]*devices.StorageEndpointProperties)

			messagingEndpointProperties["fileNotifications"] = &devices.MessagingEndpointProperties{
				LockDurationAsIso8601: utils.String(state.LockDuration),
				MaxDeliveryCount:      utils.Int32(int32(state.MaxDeliveryCount)),
				TTLAsIso8601:          utils.String(state.DefaultTTL),
			}

			storageEndpointProperties["$default"] = &devices.StorageEndpointProperties{
				AuthenticationType: devices.AuthenticationType(state.AuthenticationType),
				ConnectionString:   utils.String(state.ConnectionString),
				ContainerName:      utils.String(state.ContainerName),
				SasTTLAsIso8601:    utils.String(state.SasTTL),
			}

			if state.IdentityId != "" {
				if state.AuthenticationType != string(devices.AuthenticationTypeIdentityBased) {
					return fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
				}
				storageEndpointProperties["$default"].Identity = &devices.ManagedIdentity{
					UserAssignedIdentity: utils.String(state.IdentityId),
				}
			}

			if iotHub.Properties == nil {
				iotHub.Properties = &devices.IotHubProperties{}
			}

			iotHub.Properties.EnableFileUploadNotifications = utils.Bool(state.NotificationsEnabled)
			iotHub.Properties.MessagingEndpoints = messagingEndpointProperties
			iotHub.Properties.StorageEndpoints = storageEndpointProperties

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, iotHub, "")
			if err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the completion of the creation of %q: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotHubFileUploadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient

			id, err := parse.IotHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			iotHub, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(iotHub.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			state := IotHubFileUploadResourceModel{
				AuthenticationType:   string(devices.AuthenticationTypeKeyBased),
				ConnectionString:     "",
				ContainerName:        "",
				DefaultTTL:           "PT1H",
				IdentityId:           "",
				IotHubId:             id.ID(),
				LockDuration:         "PT1M",
				MaxDeliveryCount:     10,
				NotificationsEnabled: false,
				SasTTL:               "PT1H",
			}

			if props := iotHub.Properties; props != nil {
				if v := props.EnableFileUploadNotifications; v != nil {
					state.NotificationsEnabled = *v
				}

				if messagingEndpoint, ok := props.MessagingEndpoints["fileNotifications"]; ok {
					if v := messagingEndpoint.TTLAsIso8601; v != nil {
						state.DefaultTTL = *v
					}
					if v := messagingEndpoint.LockDurationAsIso8601; v != nil {
						state.LockDuration = *v
					}
					if v := messagingEndpoint.MaxDeliveryCount; v != nil {
						state.MaxDeliveryCount = int64(*v)
					}
				}

				if storageEndpoint, ok := props.StorageEndpoints["$default"]; ok {
					if v := string(storageEndpoint.AuthenticationType); v != "" {
						state.AuthenticationType = v
					}
					if v := storageEndpoint.ConnectionString; v != nil {
						state.ConnectionString = *v
					}
					if v := storageEndpoint.ContainerName; v != nil {
						state.ContainerName = *v
					}
					if v := storageEndpoint.Identity; v != nil && v.UserAssignedIdentity != nil {
						state.IdentityId = *v.UserAssignedIdentity
					}
					if v := storageEndpoint.SasTTLAsIso8601; v != nil {
						state.SasTTL = *v
					}
				}
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotHubFileUploadResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient

			var state IotHubFileUploadResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.IotHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
				}
			}

			if existing.Properties == nil {
				existing.Properties = &devices.IotHubProperties{}
			}

			if existing.Properties.MessagingEndpoints == nil {
				existing.Properties.MessagingEndpoints = make(map[string]*devices.MessagingEndpointProperties)
			}

			if _, ok := existing.Properties.MessagingEndpoints["fileNotifications"]; !ok {
				existing.Properties.MessagingEndpoints["fileNotifications"] = &devices.MessagingEndpointProperties{}
			}

			if existing.Properties.StorageEndpoints == nil {
				existing.Properties.StorageEndpoints = make(map[string]*devices.StorageEndpointProperties)
			}

			if _, ok := existing.Properties.StorageEndpoints["$default"]; !ok {
				existing.Properties.StorageEndpoints["$default"] = &devices.StorageEndpointProperties{}
			}

			messagingEndpoint := existing.Properties.MessagingEndpoints["fileNotifications"]
			storageEndpoint := existing.Properties.StorageEndpoints["$default"]

			if metadata.ResourceData.HasChange("notifications_enabled") {
				existing.Properties.EnableFileUploadNotifications = utils.Bool(state.NotificationsEnabled)
			}

			if metadata.ResourceData.HasChange("default_ttl") {
				messagingEndpoint.TTLAsIso8601 = utils.String(state.DefaultTTL)
			}

			if metadata.ResourceData.HasChange("lock_duration") {
				messagingEndpoint.LockDurationAsIso8601 = utils.String(state.LockDuration)
			}

			if metadata.ResourceData.HasChange("max_delivery_count") {
				messagingEndpoint.MaxDeliveryCount = utils.Int32(int32(state.MaxDeliveryCount))
			}

			if metadata.ResourceData.HasChange("authentication_type") {
				storageEndpoint.AuthenticationType = devices.AuthenticationType(state.AuthenticationType)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				storageEndpoint.ConnectionString = utils.String(state.ConnectionString)
			}

			if metadata.ResourceData.HasChange("container_name") {
				storageEndpoint.ContainerName = utils.String(state.ContainerName)
			}

			if metadata.ResourceData.HasChange("identity_id") {
				if state.IdentityId != "" {
					if state.AuthenticationType != string(devices.AuthenticationTypeIdentityBased) {
						return fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
					}
					storageEndpoint.Identity = &devices.ManagedIdentity{
						UserAssignedIdentity: utils.String(state.IdentityId),
					}
				} else {
					storageEndpoint.Identity = nil
				}
			}

			if metadata.ResourceData.HasChange("sas_ttl") {
				storageEndpoint.SasTTLAsIso8601 = utils.String(state.SasTTL)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, existing, "")
			if err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the update of %q: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotHubFileUploadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient

			id, err := parse.IotHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
				}
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, existing, "")
			if err != nil {
				return fmt.Errorf("deleting %q: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %q: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
