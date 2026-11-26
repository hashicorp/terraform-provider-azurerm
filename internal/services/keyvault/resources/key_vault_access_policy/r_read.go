// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, meta any) error {
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

	return resourceKeyVaultAccessPolicyFlatten(d, id, accessPolicy)
}

func resourceKeyVaultAccessPolicyFlatten(d *pluginsdk.ResourceData, id *parse.AccessPolicyId, accessPolicy *vaults.AccessPolicyEntry) error {
	d.Set("key_vault_id", id.KeyVaultId().ID())
	d.Set("application_id", id.ApplicationId())
	d.Set("object_id", id.ObjectID())

	if accessPolicy != nil {
		d.Set("tenant_id", accessPolicy.TenantId)

		if err := d.Set("certificate_permissions", helpers.FlattenCertificatePermissions(accessPolicy.Permissions.Certificates)); err != nil {
			return fmt.Errorf("setting `certificate_permissions`: %+v", err)
		}

		if err := d.Set("key_permissions", helpers.FlattenKeyPermissions(accessPolicy.Permissions.Keys)); err != nil {
			return fmt.Errorf("setting `key_permissions`: %+v", err)
		}

		if err := d.Set("secret_permissions", helpers.FlattenSecretPermissions(accessPolicy.Permissions.Secrets)); err != nil {
			return fmt.Errorf("setting `secret_permissions`: %+v", err)
		}

		if err := d.Set("storage_permissions", helpers.FlattenStoragePermissions(accessPolicy.Permissions.Storage)); err != nil {
			return fmt.Errorf("setting `storage_permissions`: %+v", err)
		}
	}

	return setKeyVaultAccessPolicyIdentityData(d, id)
}

func setKeyVaultAccessPolicyIdentityData(d *pluginsdk.ResourceData, id *parse.AccessPolicyId) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}
	if err := identity.Set("key_vault_id", id.KeyVaultId().ID()); err != nil {
		return fmt.Errorf("setting `key_vault_id` in resource identity: %+v", err)
	}
	if err := identity.Set("object_id", id.ObjectID()); err != nil {
		return fmt.Errorf("setting `object_id` in resource identity: %+v", err)
	}
	if err := identity.Set("application_id", id.ApplicationId()); err != nil {
		return fmt.Errorf("setting `application_id` in resource identity: %+v", err)
	}
	return nil
}

func findKeyVaultAccessPolicy(policies *[]vaults.AccessPolicyEntry, objectId string, applicationId string) *vaults.AccessPolicyEntry {
	if policies == nil {
		return nil
	}

	for _, policy := range *policies {
		if strings.EqualFold(policy.ObjectId, objectId) {
			aid := pointer.From(policy.ApplicationId)

			if strings.EqualFold(aid, applicationId) {
				return &policy
			}
		}
	}

	return nil
}

func accessPolicyRefreshFunc(ctx context.Context, client *vaults.VaultsClient, keyVaultId commonids.KeyVaultId, objectId string, applicationId string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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
