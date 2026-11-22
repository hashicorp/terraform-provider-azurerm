// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IotHubDeviceUpdateAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := resp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if existing.Properties == nil {
				existing.Properties = &deviceupdates.AccountProperties{}
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := deviceupdates.PublicNetworkAccessEnabled
				if !model.PublicNetworkAccessEnabled {
					publicNetworkAccess = deviceupdates.PublicNetworkAccessDisabled
				}
				existing.Properties.PublicNetworkAccess = &publicNetworkAccess
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = &model.Tags
			}

			if err := client.AccountsCreateThenPoll(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
