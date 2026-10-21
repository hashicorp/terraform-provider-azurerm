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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerStaticMemberModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.StaticMembers
			networkGroupId, err := staticmembers.ParseNetworkGroupID(model.NetworkGroupId)
			if err != nil {
				return err
			}

			id := staticmembers.NewStaticMemberID(networkGroupId.SubscriptionId, networkGroupId.ResourceGroupName, networkGroupId.NetworkManagerName, networkGroupId.NetworkGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			staticMember := staticmembers.StaticMember{
				Properties: &staticmembers.StaticMemberProperties{
					ResourceId: &model.TargetVNetId,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, staticMember); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
