// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKeyVaultAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultAccessPolicyCreate,
		Read:   resourceKeyVaultAccessPolicyRead,
		Update: resourceKeyVaultAccessPolicyUpdate,
		Delete: resourceKeyVaultAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccessPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"application_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"certificate_permissions": schemaCertificatePermissions(),

			"key_permissions": schemaKeyPermissions(),

			"secret_permissions": schemaSecretPermissions(),

			"storage_permissions": schemaStoragePermissions(),
		},
	}
}

func resourceKeyVaultAccessPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

				appId := ""
				if policy.ApplicationId != nil {
					appId = *policy.ApplicationId
				}
				applicationIdMatches := appId == applicationId
				if tenantIdMatches && objectIdMatches && applicationIdMatches {
					return tf.ImportAsExistsError("azurerm_key_vault_access_policy", id.ID())
				}
			}
		}
	}

	var accessPolicy vaults.AccessPolicyEntry
	certPermissionsRaw := d.Get("certificate_permissions").([]interface{})
	certPermissions := expandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]interface{})
	keyPermissions := expandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
	secretPermissions := expandSecretPermissions(secretPermissionsRaw)

	storagePermissionsRaw := d.Get("storage_permissions").([]interface{})
	storagePermissions := expandStoragePermissions(storagePermissionsRaw)

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
		Name: utils.String(keyVaultId.VaultName),
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

func resourceKeyVaultAccessPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	certPermissionsRaw := d.Get("certificate_permissions").([]interface{})
	certPermissions := expandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]interface{})
	keyPermissions := expandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
	secretPermissions := expandSecretPermissions(secretPermissionsRaw)

	storagePermissionsRaw := d.Get("storage_permissions").([]interface{})
	storagePermissions := expandStoragePermissions(storagePermissionsRaw)

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
		Name: utils.String(keyVaultId.VaultName),
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

func resourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	vaultId := id.KeyVaultId()

	resp, err := client.Get(ctx, vaultId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] parent %q was not found - removing from state", vaultId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving parent %s: %+v", vaultId, err)
	}
	var accessPolicy *vaults.AccessPolicyEntry
	if model := resp.Model; model != nil {
		accessPolicy = findKeyVaultAccessPolicy(model.Properties.AccessPolicies, id.ObjectID(), id.ApplicationId())
	}
	if accessPolicy == nil {
		log.Printf("[ERROR] Access Policy (Object ID %q / Application ID %q) was not found in %s - removing from state", id.ObjectID(), id.ApplicationId(), vaultId)
		d.SetId("")
		return nil
	}

	d.Set("key_vault_id", id.KeyVaultId().ID())
	d.Set("application_id", id.ApplicationId())
	d.Set("object_id", id.ObjectID())
	d.Set("tenant_id", accessPolicy.TenantId)

	certificatePermissions := flattenCertificatePermissions(accessPolicy.Permissions.Certificates)
	if err := d.Set("certificate_permissions", certificatePermissions); err != nil {
		return fmt.Errorf("setting `certificate_permissions`: %+v", err)
	}

	keyPermissions := flattenKeyPermissions(accessPolicy.Permissions.Keys)
	if err := d.Set("key_permissions", keyPermissions); err != nil {
		return fmt.Errorf("setting `key_permissions`: %+v", err)
	}

	secretPermissions := flattenSecretPermissions(accessPolicy.Permissions.Secrets)
	if err := d.Set("secret_permissions", secretPermissions); err != nil {
		return fmt.Errorf("setting `secret_permissions`: %+v", err)
	}

	storagePermissions := flattenStoragePermissions(accessPolicy.Permissions.Storage)
	if err := d.Set("storage_permissions", storagePermissions); err != nil {
		return fmt.Errorf("setting `storage_permissions`: %+v", err)
	}

	return nil
}

func resourceKeyVaultAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		Name: utils.String(vaultId.VaultName),
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

func findKeyVaultAccessPolicy(policies *[]vaults.AccessPolicyEntry, objectId string, applicationId string) *vaults.AccessPolicyEntry {
	if policies == nil {
		return nil
	}

	for _, policy := range *policies {
		if strings.EqualFold(policy.ObjectId, objectId) {
			aid := ""
			if policy.ApplicationId != nil {
				aid = *policy.ApplicationId
			}

			if strings.EqualFold(aid, applicationId) {
				return &policy
			}
		}
	}

	return nil
}

func accessPolicyRefreshFunc(ctx context.Context, client *vaults.VaultsClient, keyVaultId commonids.KeyVaultId, objectId string, applicationId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking for completion of Access Policy create/update")

		read, err := client.Get(ctx, keyVaultId)
		if err != nil {
			if response.WasNotFound(read.HttpResponse) {
				return "vaultnotfound", "vaultnotfound", fmt.Errorf("%s was not found", keyVaultId)
			}
		}

		var accessPolicy *vaults.AccessPolicyEntry
		if model := read.Model; model != nil {
			accessPolicy = findKeyVaultAccessPolicy(model.Properties.AccessPolicies, objectId, applicationId)
		}

		if accessPolicy != nil {
			return "found", "found", nil
		}

		return "notfound", "notfound", nil
	}
}
