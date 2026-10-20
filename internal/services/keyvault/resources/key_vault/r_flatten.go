// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	dataplane "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultFlatten(ctx context.Context, managementClient *dataplane.BaseClient, d *pluginsdk.ResourceData, id *commonids.KeyVaultId, model *vaults.Vault, includeResource bool) error {
	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	publicNetworkAccessEnabled := true

	if model != nil {
		d.Set("vault_uri", pointer.From(model.Properties.VaultUri))
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("tenant_id", model.Properties.TenantId)
		d.Set("enabled_for_deployment", model.Properties.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", model.Properties.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", model.Properties.EnabledForTemplateDeployment)
		d.Set("rbac_authorization_enabled", model.Properties.EnableRbacAuthorization)
		d.Set("purge_protection_enabled", model.Properties.EnablePurgeProtection)

		if model.Properties.PublicNetworkAccess != nil {
			publicNetworkAccessEnabled = strings.EqualFold(*model.Properties.PublicNetworkAccess, "Enabled")
		}
		d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

		// @tombuildsstuff: the API doesn't return this field if it's not configured
		// however https://docs.microsoft.com/azure/key-vault/general/soft-delete-overview
		// defaults this to 90 days, as such we're going to have to assume that for the moment
		// in lieu of anything being returned
		softDeleteRetentionDays := 90
		if model.Properties.SoftDeleteRetentionInDays != nil && *model.Properties.SoftDeleteRetentionInDays != 0 {
			softDeleteRetentionDays = int(*model.Properties.SoftDeleteRetentionInDays)
		}
		d.Set("soft_delete_retention_days", softDeleteRetentionDays)

		skuName := ""
		// The Azure API is inconsistent here, so rewrite this into the casing we expect
		// TODO: this can be removed when the new base layer is enabled?
		for _, v := range vaults.PossibleValuesForSkuName() {
			if strings.EqualFold(v, string(model.Properties.Sku.Name)) {
				skuName = v
			}
		}
		d.Set("sku_name", skuName)

		if err := d.Set("network_acls", flattenKeyVaultNetworkAcls(model.Properties.NetworkAcls)); err != nil {
			return fmt.Errorf("setting `network_acls`: %+v", err)
		}

		if err := d.Set("access_policy", helpers.FlattenAccessPolicies(model.Properties.AccessPolicies)); err != nil {
			return fmt.Errorf("setting `access_policy`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	} else {
		d.Set("vault_uri", "")
	}

	if includeResource {
		// If publicNetworkAccessEnabled is true, the data plane call should succeed.
		// (if the caller has the 'ManageContacts' certificate permissions)
		//
		// If an error is returned from the data plane call we need to return that error.
		//
		// If publicNetworkAccessEnabled is false, the data plane call should fail unless
		// there is a private endpoint connected to the key vault.
		// (and the caller has the 'ManageContacts' certificate permissions)
		//
		// We don't know if the private endpoint has been created yet, so we need
		// to ignore the error if the data plane call fails.
		contacts, err := managementClient.GetCertificateContacts(ctx, d.Get("vault_uri").(string))
		if err != nil {
			if publicNetworkAccessEnabled && (!response.WasForbidden(contacts.Response.Response) && !response.WasNotFound(contacts.Response.Response)) {
				return fmt.Errorf("retrieving `contact` for KeyVault: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenKeyVaultNetworkAcls(input *vaults.NetworkRuleSet) []any {
	bypass := string(vaults.NetworkRuleBypassOptionsAzureServices)
	defaultAction := string(vaults.NetworkRuleActionAllow)
	ipRules := make([]any, 0)
	virtualNetworkSubnetIds := make([]any, 0)

	if input != nil {
		if input.Bypass != nil {
			bypass = string(*input.Bypass)
		}
		if input.DefaultAction != nil {
			defaultAction = string(*input.DefaultAction)
		}
		if input.IPRules != nil {
			for _, v := range *input.IPRules {
				ipRules = append(ipRules, v.Value)
			}
		}
		if input.VirtualNetworkRules != nil {
			for _, v := range *input.VirtualNetworkRules {
				subnetIdRaw := v.Id
				subnetId, err := commonids.ParseSubnetIDInsensitively(subnetIdRaw)
				if err == nil {
					subnetIdRaw = subnetId.ID()
				}
				virtualNetworkSubnetIds = append(virtualNetworkSubnetIds, subnetIdRaw)
			}
		}
	}

	return []any{
		map[string]any{
			"bypass":                     bypass,
			"default_action":             defaultAction,
			"ip_rules":                   pluginsdk.NewSet(pluginsdk.HashString, ipRules),
			"virtual_network_subnet_ids": pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkSubnetIds),
		},
	}
}
