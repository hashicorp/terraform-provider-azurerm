// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// update to re-activate the security module
func resourceArmKeyVaultManagedHardwareSecurityModuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	kvClient := meta.(*clients.Client).ManagedHSMs
	hsmClient := kvClient.ManagedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hsmClient.Get(ctx, *id)
	if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model
	hasUpdate := false
	if d.HasChange("tags") {
		hasUpdate = true
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}
	if d.HasChange("network_acls") {
		hasUpdate = true
		model.Properties.NetworkAcls = expandMHSMNetworkAcls(d.Get("network_acls").([]any))
	}
	if d.HasChange("public_network_access_enabled") {
		hasUpdate = true
		publicNetworkAccessEnabled := managedhsms.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccessEnabled = managedhsms.PublicNetworkAccessDisabled
		}
		model.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccessEnabled)
	}
	if hasUpdate {
		if err := hsmClient.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
			return fmt.Errorf("updating %s tags: %+v", id, err)
		}
	}

	// security domain download to activate this module
	if ok := d.HasChange("security_domain_key_vault_certificate_ids"); ok {
		// get hsm uri
		resp, err := hsmClient.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
			return fmt.Errorf("got nil HSMUri for %s: %+v", id, err)
		}

		keyVaultClient := meta.(*clients.Client).KeyVault.ManagementClient
		encData, err := securityDomainDownload(ctx, kvClient.DataPlaneSecurityDomainsClient, *keyVaultClient, *resp.Model.Properties.HsmUri, d.Get("security_domain_key_vault_certificate_ids").([]any), d.Get("security_domain_quorum").(int))
		if err != nil {
			return fmt.Errorf("downloading security domain for %q: %+v", id, err)
		}
		d.Set("security_domain_encrypted_data", encData)
	}

	return nil
}
