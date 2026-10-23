// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).KeyVault.VaultsClient
	deletedVaultsClient := meta.(*clients.Client).KeyVault.DeletedVaultsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKeyVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))

	// Locking this resource so we don't make modifications to it at the same time if there is a
	// key vault access policy trying to update it as well
	locks.ByName(id.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultName, keyVaultResourceName)

	isPublic := d.Get("public_network_access_enabled").(bool)

	// check for the presence of an existing, live one which should be imported into the state
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_key_vault", id.ID())
		}
	}

	// before creating check to see if the key vault exists in the soft delete state
	deletedVaultId := deletedvaults.NewDeletedVaultID(id.SubscriptionId, location, id.VaultName)
	softDeletedKeyVault, err := deletedVaultsClient.VaultsGetDeleted(ctx, deletedVaultId)
	if err != nil {
		// If Terraform lacks permission to read at the Subscription we'll get 409, not 404
		if !response.WasNotFound(softDeletedKeyVault.HttpResponse) && !response.WasStatusCode(softDeletedKeyVault.HttpResponse, http.StatusForbidden) {
			return fmt.Errorf("checking for the presence of an existing Soft-Deleted Key Vault %q (Location %q): %+v", id.VaultName, location, err)
		}
	}

	// if so, does the user want us to recover it?
	recoverSoftDeletedKeyVault := false
	if !response.WasNotFound(softDeletedKeyVault.HttpResponse) && !response.WasStatusCode(softDeletedKeyVault.HttpResponse, http.StatusForbidden) {
		if !meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults {
			// this exists but the users opted out so they must import this it out-of-band
			return errors.New(optedOutOfRecoveringSoftDeletedKeyVaultErrorFmt(id.VaultName, location))
		}

		recoverSoftDeletedKeyVault = true
	}

	tenantUUID := d.Get("tenant_id").(string)
	enabledForDeployment := d.Get("enabled_for_deployment").(bool)
	enabledForDiskEncryption := d.Get("enabled_for_disk_encryption").(bool)
	enabledForTemplateDeployment := d.Get("enabled_for_template_deployment").(bool)
	rbacAuthorizationEnabled := d.Get("rbac_authorization_enabled").(bool)
	t := d.Get("tags").(map[string]any)

	policies := d.Get("access_policy").([]any)
	accessPolicies := helpers.ExpandAccessPolicies(policies)

	networkAclsRaw := d.Get("network_acls").([]any)
	networkAcls, subnetIds := expandKeyVaultNetworkAcls(networkAclsRaw)

	sku := vaults.Sku{
		Family: vaults.SkuFamilyA,
		Name:   vaults.SkuName(d.Get("sku_name").(string)),
	}

	parameters := vaults.VaultCreateOrUpdateParameters{
		Location: location,
		Properties: vaults.VaultProperties{
			TenantId:                     tenantUUID,
			Sku:                          sku,
			AccessPolicies:               accessPolicies,
			EnabledForDeployment:         &enabledForDeployment,
			EnabledForDiskEncryption:     &enabledForDiskEncryption,
			EnabledForTemplateDeployment: &enabledForTemplateDeployment,
			EnableRbacAuthorization:      pointer.To(rbacAuthorizationEnabled),
			NetworkAcls:                  networkAcls,

			// @tombuildsstuff: as of 2020-12-15 this is now defaulted on, and appears to be so in all regions
			// This has been confirmed in Azure Public and Azure China - but I couldn't find any more
			// documentation with further details
			EnableSoftDelete: pointer.To(true),
		},
		Tags: tags.Expand(t),
	}

	if isPublic {
		parameters.Properties.PublicNetworkAccess = pointer.To("Enabled")
	} else {
		parameters.Properties.PublicNetworkAccess = pointer.To("Disabled")
	}

	if purgeProtectionEnabled := d.Get("purge_protection_enabled").(bool); purgeProtectionEnabled {
		parameters.Properties.EnablePurgeProtection = pointer.To(purgeProtectionEnabled)
	}

	if v := d.Get("soft_delete_retention_days"); v != 90 {
		parameters.Properties.SoftDeleteRetentionInDays = pointer.To(int64(v.(int)))
	}

	parameters.Properties.CreateMode = pointer.To(vaults.CreateModeDefault)
	if recoverSoftDeletedKeyVault {
		parameters.Properties.CreateMode = pointer.To(vaults.CreateModeRecover)
	}

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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	read, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	vaultUri := ""
	if model := read.Model; model != nil {
		if model.Properties.VaultUri != nil {
			vaultUri = *model.Properties.VaultUri
		}
	}

	if vaultUri == "" {
		return fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", id)
	}

	meta.(*clients.Client).KeyVault.AddToCache(id, vaultUri)

	// When Public Network Access is Enabled (i.e. it's Public) we can hit the Data Plane API until
	// we get a valid response repeatedly - ensuring that the API is fully online before proceeding.
	//
	// This works around an issue where the provisioning of dependent resources fails, due to the
	// Key Vault not being fully online - which is a particular issue when recreating the Key Vault.
	//
	// When Public Network Access is Disabled (i.e. it's Private) we don't poll - meaning that users
	// are more likely to encounter issues in downstream resources (particularly when using Private
	// Link due to DNS replication delays) - however there isn't a great deal we can do about that
	// given the Data Plane API isn't going to be publicly available.
	//
	// As such we poll to check the Key Vault is available, if it's public, to ensure that downstream
	// operations can succeed.
	if isPublic {
		log.Printf("[DEBUG] Waiting for %s to become available", id)
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("internal-error: context had no deadline")
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"pending"},
			Target:                    []string{"available"},
			Refresh:                   keyVaultRefreshFunc(ctx, vaultUri),
			Delay:                     30 * time.Second,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 10,
			Timeout:                   time.Until(deadline),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %s", id, err)
		}
	}

	return resourceKeyVaultRead(d, meta)
}

func expandKeyVaultNetworkAcls(input []any) (*vaults.NetworkRuleSet, []string) {
	subnetIds := make([]string, 0)
	if len(input) == 0 {
		return nil, subnetIds
	}

	v := input[0].(map[string]any)

	bypass := v["bypass"].(string)
	defaultAction := v["default_action"].(string)

	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]vaults.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := vaults.IPRule{
			Value: v.(string),
		}
		ipRules = append(ipRules, rule)
	}

	networkRulesRaw := v["virtual_network_subnet_ids"].(*pluginsdk.Set)
	networkRules := make([]vaults.VirtualNetworkRule, 0)
	for _, v := range networkRulesRaw.List() {
		rawId := v.(string)
		subnetIds = append(subnetIds, rawId)
		rule := vaults.VirtualNetworkRule{
			Id: rawId,
		}
		networkRules = append(networkRules, rule)
	}

	ruleSet := vaults.NetworkRuleSet{
		Bypass:              pointer.ToEnum[vaults.NetworkRuleBypassOptions](bypass),
		DefaultAction:       pointer.ToEnum[vaults.NetworkRuleAction](defaultAction),
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}
	return &ruleSet, subnetIds
}
