// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_file_upload

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient

			id, err := parse.IotHubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			iotHub, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if response.WasNotFound(iotHub.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			state := IotHubFileUploadResourceModel{
				AuthenticationType: string(devices.AuthenticationTypeKeyBased),
				DefaultTTL:         "PT1H",
				IotHubId:           id.ID(),
				LockDuration:       "PT1M",
				MaxDeliveryCount:   10,
				SasTTL:             "PT1H",
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
