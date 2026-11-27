// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2025-04-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.VideoIndexer.AccountClient

			var account AccountModel
			if err := metadata.Decode(&account); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := accounts.AccountPatch{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(account.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(account.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("storage") {
				payload.Properties = &accounts.AccountPropertiesForPatchRequest{
					StorageServices: expandStorageForUpdate(account.Storage),
				}
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				if payload.Properties == nil {
					payload.Properties = &accounts.AccountPropertiesForPatchRequest{}
				}
				payload.Properties.PublicNetworkAccess = pointer.ToEnum[accounts.PublicNetworkAccess](account.PublicNetworkAccess)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
