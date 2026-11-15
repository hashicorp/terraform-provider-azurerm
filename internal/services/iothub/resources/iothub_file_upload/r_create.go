// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_file_upload

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Create() sdk.ResourceFunc {
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
				if response.WasNotFound(iotHub.Response.Response) {
					return fmt.Errorf("%q was not found", id)
				}

				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if iotHub.Properties != nil && iotHub.Properties.MessagingEndpoints != nil {
					if storageEndpoint, ok := iotHub.Properties.StorageEndpoints["$default"]; ok {
						if storageEndpoint.ConnectionString != nil && *storageEndpoint.ConnectionString != "" && storageEndpoint.ContainerName != nil && *storageEndpoint.ContainerName != "" {
							return metadata.ResourceRequiresImport(r.ResourceType(), id)
						}
					}
				}
			}

			messagingEndpointProperties := make(map[string]*devices.MessagingEndpointProperties)
			storageEndpointProperties := make(map[string]*devices.StorageEndpointProperties)

			messagingEndpointProperties["fileNotifications"] = &devices.MessagingEndpointProperties{
				LockDurationAsIso8601: pointer.To(state.LockDuration),
				MaxDeliveryCount:      pointer.To(int32(state.MaxDeliveryCount)),
				TTLAsIso8601:          pointer.To(state.DefaultTTL),
			}

			storageEndpointProperties["$default"] = &devices.StorageEndpointProperties{
				AuthenticationType: devices.AuthenticationType(state.AuthenticationType),
				ConnectionString:   pointer.To(state.ConnectionString),
				ContainerName:      pointer.To(state.ContainerName),
				SasTTLAsIso8601:    pointer.To(state.SasTTL),
			}

			if state.IdentityId != "" {
				if state.AuthenticationType != string(devices.AuthenticationTypeIdentityBased) {
					return fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
				}
				storageEndpointProperties["$default"].Identity = &devices.ManagedIdentity{
					UserAssignedIdentity: pointer.To(state.IdentityId),
				}
			}

			if iotHub.Properties == nil {
				iotHub.Properties = &devices.IotHubProperties{}
			}

			iotHub.Properties.EnableFileUploadNotifications = pointer.To(state.NotificationsEnabled)
			iotHub.Properties.MessagingEndpoints = messagingEndpointProperties
			iotHub.Properties.StorageEndpoints = storageEndpointProperties

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, iotHub, "")
			if err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			metadata.SetID(id)

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the completion of the creation of %q: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
