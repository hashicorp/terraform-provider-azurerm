// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmKeyVaultManagedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedHSMs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managedhsms.NewManagedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ManagedHsmClient.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_key_vault_managed_hardware_security_module", id.ID())
		}
	}

	publicNetworkAccessEnabled := managedhsms.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccessEnabled = managedhsms.PublicNetworkAccessDisabled
	}
	hsm := managedhsms.ManagedHsm{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &managedhsms.ManagedHsmProperties{
			InitialAdminObjectIds:     pluginsdk.ExpandStringSlice(d.Get("admin_object_ids").(*pluginsdk.Set).List()),
			CreateMode:                pointer.To(managedhsms.CreateModeDefault),
			EnableSoftDelete:          pointer.To(true),
			SoftDeleteRetentionInDays: pointer.To(int64(d.Get("soft_delete_retention_days").(int))),
			EnablePurgeProtection:     pointer.To(d.Get("purge_protection_enabled").(bool)),
			PublicNetworkAccess:       pointer.To(publicNetworkAccessEnabled),
			NetworkAcls:               expandMHSMNetworkAcls(d.Get("network_acls").([]any)),
		},
		Sku: &managedhsms.ManagedHsmSku{
			Family: managedhsms.ManagedHsmSkuFamilyB,
			Name:   managedhsms.ManagedHsmSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}
	if tenantId := d.Get("tenant_id").(string); tenantId != "" {
		hsm.Properties.TenantId = pointer.To(tenantId)
	}

	if err := client.ManagedHsmClient.CreateOrUpdateCallbackThenPoll(ctx, id, hsm, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	dataPlaneUri := ""
	resp, err := client.ManagedHsmClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if model := resp.Model; model != nil && model.Properties != nil && model.Properties.HsmUri != nil {
		dataPlaneUri = *model.Properties.HsmUri
	}
	if dataPlaneUri == "" {
		return fmt.Errorf("retrieving %s: `properties.HsmUri` was nil", id)
	}
	client.AddToCache(id, dataPlaneUri)

	// security domain download to activate this module
	if ok := d.HasChange("security_domain_key_vault_certificate_ids"); ok {
		// get hsm uri
		resp, err := client.ManagedHsmClient.Get(ctx, id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
			return fmt.Errorf("got nil HSMUri for %s: %+v", id, err)
		}

		keyVaultClient := meta.(*clients.Client).KeyVault.ManagementClient
		encData, err := securityDomainDownload(ctx, client.DataPlaneSecurityDomainsClient, *keyVaultClient, *resp.Model.Properties.HsmUri, d.Get("security_domain_key_vault_certificate_ids").([]any), d.Get("security_domain_quorum").(int))
		if err != nil {
			return fmt.Errorf("downloading security domain for %q: %+v", id, err)
		}
		d.Set("security_domain_encrypted_data", encData)
	}

	return resourceArmKeyVaultManagedHardwareSecurityModuleRead(d, meta)
}

func expandMHSMNetworkAcls(input []any) *managedhsms.MHSMNetworkRuleSet {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &managedhsms.MHSMNetworkRuleSet{
		Bypass:        pointer.ToEnum[managedhsms.NetworkRuleBypassOptions](v["bypass"].(string)),
		DefaultAction: pointer.ToEnum[managedhsms.NetworkRuleAction](v["default_action"].(string)),
	}
}
