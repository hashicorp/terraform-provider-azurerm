// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account_encryption

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient
			keyVaultsClient := metadata.Client.KeyVault
			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			var state models.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			if metadata.ResourceData.HasChanges(
				"user_assigned_identity_id",
				"system_assigned_identity_principal_id",
				"encryption_key",
				"federated_client_id",
			) {
				encryptionExpanded, err := expandEncryption(ctx, state.EncryptionKey, keyVaultsClient, subscriptionId, pointer.To(state))
				if err != nil {
					return err
				}

				update.Properties.Encryption = encryptionExpanded

				if err := client.AccountsUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
