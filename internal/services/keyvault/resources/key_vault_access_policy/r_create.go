// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultAccessPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tenantId := d.Get("tenant_id").(string)
	objectId := d.Get("object_id").(string)
	applicationId := d.Get("application_id").(string)

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewAccessPolicyId(*keyVaultId, objectId, applicationId)

	// Locking to prevent parallel changes causing issues
	locks.ByName(keyVaultId.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultId.VaultName, keyVaultResourceName)

	keyVault, err := client.Get(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("retrieving parent %s: %+v", *keyVaultId, err)
	}

	if model := keyVault.Model; model != nil {
		// we don't reuse findKeyVaultAccessPolicy since we're also checking the Tenant ID
		if model.Properties.AccessPolicies != nil {
			for _, policy := range *model.Properties.AccessPolicies {
				tenantIdMatches := policy.TenantId == tenantId
				objectIdMatches := policy.ObjectId == objectId

				applicationIdMatches := pointer.From(policy.ApplicationId) == applicationId
				if tenantIdMatches && objectIdMatches && applicationIdMatches {
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_key_vault_access_policy", id.ID())
					}
				}
			}
		}
	}

	var accessPolicy vaults.AccessPolicyEntry
	certPermissionsRaw := d.Get("certificate_permissions").([]any)
	certPermissions := helpers.ExpandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]any)
	keyPermissions := helpers.ExpandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]any)
	secretPermissions := helpers.ExpandSecretPermissions(secretPermissionsRaw)

	storagePermissionsRaw := d.Get("storage_permissions").([]any)
	storagePermissions := helpers.ExpandStoragePermissions(storagePermissionsRaw)

	accessPolicy = vaults.AccessPolicyEntry{
		ObjectId: objectId,
		TenantId: tenantId,
		Permissions: vaults.Permissions{
			Certificates: certPermissions,
			Keys:         keyPermissions,
			Secrets:      secretPermissions,
			Storage:      storagePermissions,
		},
	}

	if applicationId != "" {
		accessPolicy.ApplicationId = pointer.To(applicationId)
	}

	parameters := vaults.VaultAccessPolicyParameters{
		Name: pointer.To(keyVaultId.VaultName),
		Properties: vaults.VaultAccessPolicyProperties{
			AccessPolicies: []vaults.AccessPolicyEntry{
				accessPolicy,
			},
		},
	}

	updateId := vaults.NewOperationKindID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, vaults.AccessPolicyUpdateKindAdd)
	if _, err = client.UpdateAccessPolicy(ctx, updateId, parameters); err != nil {
		return fmt.Errorf("creating Access Policy (Object ID %q / Application ID %q) within %s: %+v", objectId, applicationId, *keyVaultId, err)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound", "vaultnotfound"},
		Target:                    []string{"found"},
		Refresh:                   accessPolicyRefreshFunc(ctx, client, *keyVaultId, objectId, applicationId),
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("creating Access Policy (Object ID: %q) within %s: %+v", objectId, keyVaultId, err)
	}

	d.SetId(id.ID())
	return nil
}
