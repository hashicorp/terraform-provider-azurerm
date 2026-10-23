// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultAccessPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccessPolicyID(d.Id())
	if err != nil {
		return err
	}
	keyVaultId := id.KeyVaultId()

	// Locking to prevent parallel changes causing issues
	locks.ByName(keyVaultId.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultId.VaultName, keyVaultResourceName)

	certPermissionsRaw := d.Get("certificate_permissions").([]any)
	certPermissions := helpers.ExpandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]any)
	keyPermissions := helpers.ExpandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]any)
	secretPermissions := helpers.ExpandSecretPermissions(secretPermissionsRaw)

	storagePermissionsRaw := d.Get("storage_permissions").([]any)
	storagePermissions := helpers.ExpandStoragePermissions(storagePermissionsRaw)

	accessPolicy := vaults.AccessPolicyEntry{
		ObjectId: id.ObjectID(),
		TenantId: d.Get("tenant_id").(string),
		Permissions: vaults.Permissions{
			Certificates: certPermissions,
			Keys:         keyPermissions,
			Secrets:      secretPermissions,
			Storage:      storagePermissions,
		},
	}
	if id.ApplicationId() != "" {
		accessPolicy.ApplicationId = pointer.To(id.ApplicationId())
	}

	parameters := vaults.VaultAccessPolicyParameters{
		Name: pointer.To(keyVaultId.VaultName),
		Properties: vaults.VaultAccessPolicyProperties{
			AccessPolicies: []vaults.AccessPolicyEntry{
				accessPolicy,
			},
		},
	}

	updateId := vaults.NewOperationKindID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, vaults.AccessPolicyUpdateKindReplace)
	if _, err = client.UpdateAccessPolicy(ctx, updateId, parameters); err != nil {
		return fmt.Errorf("updating Access Policy (Object ID %q / Application ID %q) for %s: %+v", id.ObjectID(), id.ApplicationId(), keyVaultId, err)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound", "vaultnotfound"},
		Target:                    []string{"found"},
		Refresh:                   accessPolicyRefreshFunc(ctx, client, keyVaultId, id.ObjectID(), id.ApplicationId()),
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for update of Access Policy (Object ID: %q) for %s: %+v", id.ObjectID(), keyVaultId, err)
	}

	return nil
}
