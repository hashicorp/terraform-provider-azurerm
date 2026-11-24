// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q:%+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("creating %s:%+v", id, err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
			state := VaultGuardProxyModel{
				VaultId: vaultId.ID(),
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.ResourceGuardResourceId != nil {
				// Resource Guards created outside of Terraform may have a lowercased `resourceGroups` segment
				// which is persisted to the proxy even when the API is provided correct casing on creation, so we need to parse this insensitively before setting it into state
				resourceGuardID, err := resourceguardresources.ParseResourceGuardIDInsensitively(*resp.Model.Properties.ResourceGuardResourceId)
				if err != nil {
					return err
				}
				state.ResourceGuardId = resourceGuardID.ID()
			}

			return metadata.Encode(&state)
		},
	}
}

const (
	VaultGuardResourceType           = "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies"
	VaultGuardProxyDeleteRequestName = "default" // this name does not matter, this value comes from Portal.
)
