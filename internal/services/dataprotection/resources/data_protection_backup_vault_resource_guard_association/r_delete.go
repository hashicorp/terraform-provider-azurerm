// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_resource_guard_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataProtectionBackupVaultResourceGuardAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.ResourceGuardProxyClient

			id, err := resourceguardproxybaseresources.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			guardId, err := resourceguardresources.ParseResourceGuardID(model.DataProtectionResourceGuardId)
			if err != nil {
				return err
			}

			requestId := resourceguards.NewDeleteResourceGuardProxyRequestID(guardId.SubscriptionId, guardId.ResourceGroupName, guardId.ResourceGuardName, dataProtectionBackupVaultResourceGuardAssociationDeleteRequestName)
			unlock := resourceguardproxybaseresources.UnlockDeleteRequest{
				ResourceGuardOperationRequests: pointer.To([]string{requestId.ID()}),
			}

			if _, err = client.DppResourceGuardProxyUnlockDelete(ctx, *id, unlock, resourceguardproxybaseresources.DefaultDppResourceGuardProxyUnlockDeleteOperationOptions()); err != nil {
				return fmt.Errorf("unlocking delete for %s: %+v", *id, err)
			}

			if _, err = client.DppResourceGuardProxyDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
