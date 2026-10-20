// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKeyVaultID(d.Id())
	if err != nil {
		return err
	}

	// Locking this resource so we don't make modifications to it at the same time if there is a
	// key vault access policy trying to update it as well
	locks.ByName(id.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultName, keyVaultResourceName)

	// first pull the existing key vault since we need to lock on several bits of its information
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", *id)
	}

	update := vaults.VaultPatchParameters{}
	isPublic := d.Get("public_network_access_enabled").(bool)
	if update.Properties == nil {
		update.Properties = &vaults.VaultPatchProperties{}
	}

	if d.HasChange("access_policy") {
		policiesRaw := d.Get("access_policy").([]any)
		update.Properties.AccessPolicies = helpers.ExpandAccessPolicies(policiesRaw)
	}

	if d.HasChange("enabled_for_deployment") {
		update.Properties.EnabledForDeployment = pointer.To(d.Get("enabled_for_deployment").(bool))
	}

	if d.HasChange("enabled_for_disk_encryption") {
		update.Properties.EnabledForDiskEncryption = pointer.To(d.Get("enabled_for_disk_encryption").(bool))
	}

	if d.HasChange("enabled_for_template_deployment") {
		update.Properties.EnabledForTemplateDeployment = pointer.To(d.Get("enabled_for_template_deployment").(bool))
	}

	if d.HasChange("rbac_authorization_enabled") {
		update.Properties.EnableRbacAuthorization = pointer.To(d.Get("rbac_authorization_enabled").(bool))
	}

	if d.HasChange("network_acls") {
		networkAclsRaw := d.Get("network_acls").([]any)
		networkAcls, subnetIds := expandKeyVaultNetworkAcls(networkAclsRaw)

		// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
		virtualNetworkNames := make([]string, 0)
		for _, v := range subnetIds {
			id, err := commonids.ParseSubnetIDInsensitively(v)
			if err != nil {
				return err
			}

			if !slices.Contains(virtualNetworkNames, id.VirtualNetworkName) {
				virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
			}
		}

		locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

		update.Properties.NetworkAcls = networkAcls
	}

	if d.HasChange("purge_protection_enabled") {
		newValue := d.Get("purge_protection_enabled").(bool)

		// existing.Properties guaranteed non-nil above
		oldValue := pointer.From(existing.Model.Properties.EnablePurgeProtection)

		// whilst this should have got caught in the customizeDiff this won't work if that fields interpolated
		// hence the double-checking here
		if oldValue && !newValue {
			return fmt.Errorf("updating %s: once Purge Protection has been Enabled it's not possible to disable it", *id)
		}

		update.Properties.EnablePurgeProtection = pointer.To(newValue)

		if newValue {
			// When the KV was created with a version prior to v2.42 and the `soft_delete_enabled` is set to false, setting `purge_protection_enabled` to `true` would not work when updating KV with v2.42 or later of terraform provider.
			// This is because the `purge_protection_enabled` only works when soft delete is enabled.
			// Since version v2.42 of the Azure Provider and later force the value of `soft_delete_enabled` to be true, we should set `EnableSoftDelete` to true when `purge_protection_enabled` is enabled to make sure it works in this case.
			update.Properties.EnableSoftDelete = pointer.To(true)
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		if isPublic {
			update.Properties.PublicNetworkAccess = pointer.To("Enabled")
		} else {
			update.Properties.PublicNetworkAccess = pointer.To("Disabled")
		}
	}

	if d.HasChange("sku_name") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.Sku = &vaults.Sku{
			Family: vaults.SkuFamilyA,
			Name:   vaults.SkuName(d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("soft_delete_retention_days") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		// existing.Properties guaranteed non-nil above
		oldValue := pointer.From(existing.Model.Properties.SoftDeleteRetentionInDays)

		// whilst this should have got caught in the customizeDiff this won't work if that fields interpolated
		// hence the double-checking here
		if oldValue != 0 {
			// Code="BadRequest" Message="The property \"softDeleteRetentionInDays\" has been set already and it can't be modified."
			return fmt.Errorf("updating %s: once `soft_delete_retention_days` has been configured it cannot be modified", *id)
		}

		update.Properties.SoftDeleteRetentionInDays = pointer.To(int64(d.Get("soft_delete_retention_days").(int)))
	}

	if d.HasChange("tenant_id") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.TenantId = pointer.To(d.Get("tenant_id").(string))
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]any)
		update.Tags = tags.Expand(t)
	}

	if _, err := client.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.Partial(false)

	return resourceKeyVaultRead(d, meta)
}
