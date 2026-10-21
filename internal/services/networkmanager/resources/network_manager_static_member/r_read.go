// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_static_member

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.StaticMembers

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties
			state := ManagerStaticMemberModel{
				Name:           id.StaticMemberName,
				NetworkGroupId: staticmembers.NewNetworkGroupID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.NetworkGroupName).ID(),
			}

			if properties.Region != nil {
				state.Region = *properties.Region
			}

			if properties.ResourceId != nil {
				state.TargetVNetId = *properties.ResourceId
			}

			return metadata.Encode(&state)
		},
	}
}
