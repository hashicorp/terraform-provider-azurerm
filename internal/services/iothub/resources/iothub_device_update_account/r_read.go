// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountsGet(ctx, *id)
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

			state := IotHubDeviceUpdateAccountModel{
				Name:              id.AccountName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				if properties.HostName != nil {
					state.HostName = *properties.HostName
				}

				publicNetworkAccessEnabled := true
				if properties.PublicNetworkAccess != nil && *properties.PublicNetworkAccess == deviceupdates.PublicNetworkAccessDisabled {
					publicNetworkAccessEnabled = false
				}
				state.PublicNetworkAccessEnabled = publicNetworkAccessEnabled

				sku := deviceupdates.SKUStandard
				if properties.Sku != nil {
					sku = *properties.Sku
				}
				state.Sku = sku
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}
