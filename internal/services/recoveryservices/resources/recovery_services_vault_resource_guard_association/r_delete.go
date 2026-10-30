// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguards"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan VaultGuardProxyModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %w", err)
			}
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			guardId, err := resourceguardresources.ParseResourceGuardID(plan.ResourceGuardId)
			if err != nil {
				return err
			}

			requestId := resourceguards.NewDeleteResourceGuardProxyRequestID(guardId.SubscriptionId, guardId.ResourceGroupName, guardId.ResourceGuardName, VaultGuardProxyDeleteRequestName)

			unlock := resourceguardproxy.UnlockDeleteRequest{
				ResourceGuardOperationRequests: pointer.To([]string{requestId.ID()}),
			}

			if _, err = client.UnlockDelete(ctx, *id, unlock); err != nil {
				return fmt.Errorf("unlocking delete %s:%+v", id, err)
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
