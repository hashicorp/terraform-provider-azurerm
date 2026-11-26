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

func (r Resource) Update() sdk.ResourceFunc {
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
				if !response.WasNotFound(existing.Response.Response) {
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
				existing.Properties.EnableFileUploadNotifications = pointer.To(state.NotificationsEnabled)
			}

			if metadata.ResourceData.HasChange("default_ttl") {
				messagingEndpoint.TTLAsIso8601 = pointer.To(state.DefaultTTL)
			}

			if metadata.ResourceData.HasChange("lock_duration") {
				messagingEndpoint.LockDurationAsIso8601 = pointer.To(state.LockDuration)
			}

			if metadata.ResourceData.HasChange("max_delivery_count") {
				messagingEndpoint.MaxDeliveryCount = pointer.To(int32(state.MaxDeliveryCount))
			}

			if metadata.ResourceData.HasChange("authentication_type") {
				storageEndpoint.AuthenticationType = devices.AuthenticationType(state.AuthenticationType)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				storageEndpoint.ConnectionString = pointer.To(state.ConnectionString)
			}

			if metadata.ResourceData.HasChange("container_name") {
				storageEndpoint.ContainerName = pointer.To(state.ContainerName)
			}

			if metadata.ResourceData.HasChange("identity_id") {
				if state.IdentityId != "" {
					if state.AuthenticationType != string(devices.AuthenticationTypeIdentityBased) {
						return fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
					}
					storageEndpoint.Identity = &devices.ManagedIdentity{
						UserAssignedIdentity: pointer.To(state.IdentityId),
					}
				} else {
					storageEndpoint.Identity = nil
				}
			}

			if metadata.ResourceData.HasChange("sas_ttl") {
				storageEndpoint.SasTTLAsIso8601 = pointer.To(state.SasTTL)
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
