// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan VaultGuardProxyModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %w", err)
			}
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			vaultId, err := vaults.ParseVaultID(plan.VaultId)
			if err != nil {
				return fmt.Errorf("parsing vault id %w", err)
			}

			name := "VaultProxy"
			id := resourceguardproxy.NewBackupResourceGuardProxyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking presence of %s:%+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_recovery_services_vault_resource_guard_association", id.ID())
				}
			}

			proxy := resourceguardproxy.ResourceGuardProxyBaseResource{
				Id:   pointer.To(id.ID()),
				Type: pointer.To(VaultGuardResourceType),
				Properties: pointer.To(resourceguardproxy.ResourceGuardProxyBase{
					ResourceGuardResourceId: pointer.To(plan.ResourceGuardId),
				}),
			}

			if _, err = client.Put(ctx, id, proxy); err != nil {
				return fmt.Errorf("creating %s:%w", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
