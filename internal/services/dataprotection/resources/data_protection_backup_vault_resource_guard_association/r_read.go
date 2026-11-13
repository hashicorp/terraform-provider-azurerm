// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_resource_guard_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.ResourceGuardProxyClient

			id, err := resourceguardproxybaseresources.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DppResourceGuardProxyGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state, err := r.flatten(*id, resp.Model)
			if err != nil {
				return err
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(state)
		},
	}
}

func (Resource) flatten(id resourceguardproxybaseresources.BackupResourceGuardProxyId, model *resourceguardproxybaseresources.ResourceGuardProxyBaseResource) (*DataProtectionBackupVaultResourceGuardAssociationModel, error) {
	vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	state := &DataProtectionBackupVaultResourceGuardAssociationModel{
		DataProtectionBackupVaultId: vaultId.ID(),
	}

	if model != nil && model.Properties != nil && model.Properties.ResourceGuardResourceId != nil {
		guardId, err := resourceguardresources.ParseResourceGuardIDInsensitively(*model.Properties.ResourceGuardResourceId)
		if err != nil {
			return nil, err
		}
		state.DataProtectionResourceGuardId = guardId.ID()
	}

	return state, nil
}
