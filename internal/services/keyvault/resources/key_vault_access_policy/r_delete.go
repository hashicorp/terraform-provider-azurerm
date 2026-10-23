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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultAccessPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	vaultId := id.KeyVaultId()

	// Locking to prevent parallel changes causing issues
	locks.ByName(vaultId.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(vaultId.VaultName, keyVaultResourceName)

	keyVault, err := client.Get(ctx, vaultId)
	if err != nil {
		return fmt.Errorf("retrieving parent %s: %+v", vaultId, err)
	}

	// To remove a policy correctly, we need to send it with all permissions in the correct case which may have drifted
	// in config over time so we read it back from the vault by objectId
	var accessPolicyRaw *vaults.AccessPolicyEntry
	if model := keyVault.Model; model != nil {
		accessPolicyRaw = findKeyVaultAccessPolicy(model.Properties.AccessPolicies, id.ObjectID(), id.ApplicationId())
	}
	if accessPolicyRaw == nil {
		return fmt.Errorf("unable to find Access Policy (Object ID %q / Application ID %q) on %s", id.ObjectID(), id.ApplicationId(), vaultId)
	}
	accessPolicy := *accessPolicyRaw
	if id.ApplicationId() != "" {
		accessPolicy.ApplicationId = pointer.To(id.ApplicationId())
	}
	parameters := vaults.VaultAccessPolicyParameters{
		Name: pointer.To(vaultId.VaultName),
		Properties: vaults.VaultAccessPolicyProperties{
			AccessPolicies: []vaults.AccessPolicyEntry{
				accessPolicy,
			},
		},
	}

	keyVaultId := id.KeyVaultId()
	updateId := vaults.NewOperationKindID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, vaults.AccessPolicyUpdateKindRemove)
	if _, err = client.UpdateAccessPolicy(ctx, updateId, parameters); err != nil {
		return fmt.Errorf("removing Access Policy (Object ID %q / Application ID %q) for %s: %+v", id.ObjectID(), id.ApplicationId(), vaultId, err)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"found", "vaultnotfound"},
		Target:                    []string{"notfound"},
		Refresh:                   accessPolicyRefreshFunc(ctx, client, vaultId, id.ObjectID(), id.ApplicationId()),
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for removal of Access Policy (Object ID: %q) from %s: %+v", id.ObjectID(), keyVaultId, err)
	}

	return nil
}
