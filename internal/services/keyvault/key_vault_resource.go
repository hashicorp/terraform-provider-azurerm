// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

var keyVaultResourceName = "azurerm_key_vault"

func resourceKeyVault() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceKeyVaultCreate,
		Read:   resourceKeyVaultRead,
		Update: resourceKeyVaultUpdate,
		Delete: resourceKeyVaultDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseKeyVaultID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KeyVaultV0ToV1{},
			1: migration.KeyVaultV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VaultName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(vaults.SkuNameStandard),
					string(vaults.SkuNamePremium),
				}, false),
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"access_policy": {
				Type:       pluginsdk.TypeList,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				MaxItems:   1024,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"object_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"application_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.IsUUIDOrEmpty,
						},
						"certificate_permissions": schemaCertificatePermissions(),
						"key_permissions":         schemaKeyPermissions(),
						"secret_permissions":      schemaSecretPermissions(),
						"storage_permissions":     schemaStoragePermissions(),
					},
				},
			},

			// NOTE: To unblock customers where they had previously deployed a key vault with
			// contacts, but cannot re-deploy the key vault. I am adding support for the contact
			// field back into the resource for UPDATE ONLY. If this is a new resource and the
			// contact field is defined in the configuration file it will now throw an error.
			// This will allow legacy key vaults to continue to work as the previously have
			// and enforces our new model of separating out the data plane call into its
			// own resource (e.g., contacts)...
			"contact": {
				Type:       pluginsdk.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "As the `contact` property requires reaching out to the dataplane, to better support private endpoints and keyvaults with public network access disabled, new key vaults with the `contact` field defined in the configuration file will now be required to use the `azurerm_key_vault_certificate_contacts` resource instead of the exposed `contact` field in the key vault resource itself.",
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"phone": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"enabled_for_deployment": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enabled_for_disk_encryption": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enabled_for_template_deployment": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_rbac_authorization": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(vaults.NetworkRuleActionAllow),
								string(vaults.NetworkRuleActionDeny),
							}, false),
						},
						"bypass": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(vaults.NetworkRuleBypassOptionsNone),
								string(vaults.NetworkRuleBypassOptionsAzureServices),
							}, false),
						},
						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									commonValidate.IPv4Address,
									commonValidate.CIDR,
								),
							},
							Set: set.HashIPv4AddressOrCIDR,
						},
						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      set.HashStringIgnoreCase,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"soft_delete_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      90,
				ValidateFunc: validation.IntBetween(7, 90),
			},

			"tags": commonschema.Tags(),

			// Computed
			"vault_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	return resource
}

func resourceKeyVaultCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	managementClient := meta.(*clients.Client).KeyVault.ManagementClient // TODO: Remove in 4.0
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKeyVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))

	// Locking this resource so we don't make modifications to it at the same time if there is a
	// key vault access policy trying to update it as well
	locks.ByName(id.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultName, keyVaultResourceName)

	isPublic := d.Get("public_network_access_enabled").(bool)
	contactRaw := d.Get("contact").(*pluginsdk.Set).List()
	contactCount := len(contactRaw)

	if contactCount > 0 {
		if features.FourPointOhBeta() {
			// In v4.0 providers block creation of all key vaults if the configuration
			// file contains a 'contact' field...
			return fmt.Errorf("%s: `contact` field is not supported for new key vaults", id)
		} else if !isPublic {
			// In v3.x providers block creation of key vaults if 'public_network_access_enabled'
			// is 'false'...
			return fmt.Errorf("%s: `contact` cannot be specified when `public_network_access_enabled` is set to `false`", id)
		}
	}

	// check for the presence of an existing, live one which should be imported into the state
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_key_vault", id.ID())
	}

	// before creating check to see if the key vault exists in the soft delete state
	deletedVaultId := vaults.NewDeletedVaultID(id.SubscriptionId, location, id.VaultName)
	softDeletedKeyVault, err := client.GetDeleted(ctx, deletedVaultId)
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
			return fmt.Errorf(optedOutOfRecoveringSoftDeletedKeyVaultErrorFmt(id.VaultName, location))
		}

		recoverSoftDeletedKeyVault = true
	}

	tenantUUID := d.Get("tenant_id").(string)
	enabledForDeployment := d.Get("enabled_for_deployment").(bool)
	enabledForDiskEncryption := d.Get("enabled_for_disk_encryption").(bool)
	enabledForTemplateDeployment := d.Get("enabled_for_template_deployment").(bool)
	enableRbacAuthorization := d.Get("enable_rbac_authorization").(bool)
	t := d.Get("tags").(map[string]interface{})

	policies := d.Get("access_policy").([]interface{})
	accessPolicies := expandAccessPolicies(policies)

	networkAclsRaw := d.Get("network_acls").([]interface{})
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
			EnableRbacAuthorization:      &enableRbacAuthorization,
			NetworkAcls:                  networkAcls,

			// @tombuildsstuff: as of 2020-12-15 this is now defaulted on, and appears to be so in all regions
			// This has been confirmed in Azure Public and Azure China - but I couldn't find any more
			// documentation with further details
			EnableSoftDelete: utils.Bool(true),
		},
		Tags: tags.Expand(t),
	}

	if isPublic {
		parameters.Properties.PublicNetworkAccess = utils.String("Enabled")
	} else {
		parameters.Properties.PublicNetworkAccess = utils.String("Disabled")
	}

	if purgeProtectionEnabled := d.Get("purge_protection_enabled").(bool); purgeProtectionEnabled {
		parameters.Properties.EnablePurgeProtection = utils.Bool(purgeProtectionEnabled)
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
		if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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

	d.SetId(id.ID())

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
			Refresh:                   keyVaultRefreshFunc(vaultUri),
			Delay:                     30 * time.Second,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 10,
			Timeout:                   time.Until(deadline),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %s", id, err)
		}
	}

	// Only call the data plane if the 'contact' field has been defined...
	if contactCount > 0 {
		contacts := dataplane.Contacts{
			ContactList: expandKeyVaultCertificateContactList(contactRaw),
		}

		if _, err := managementClient.SetCertificateContacts(ctx, vaultUri, contacts); err != nil {
			return fmt.Errorf("failed to set Contacts for %s: %+v", id, err)
		}
	}

	return resourceKeyVaultRead(d, meta)
}

func resourceKeyVaultUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	managementClient := meta.(*clients.Client).KeyVault.ManagementClient
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

	d.Partial(true)

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

	if d.HasChange("access_policy") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		policiesRaw := d.Get("access_policy").([]interface{})
		accessPolicies := expandAccessPolicies(policiesRaw)
		update.Properties.AccessPolicies = accessPolicies
	}

	if d.HasChange("enabled_for_deployment") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.EnabledForDeployment = utils.Bool(d.Get("enabled_for_deployment").(bool))
	}

	if d.HasChange("enabled_for_disk_encryption") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.EnabledForDiskEncryption = utils.Bool(d.Get("enabled_for_disk_encryption").(bool))
	}

	if d.HasChange("enabled_for_template_deployment") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.EnabledForTemplateDeployment = utils.Bool(d.Get("enabled_for_template_deployment").(bool))
	}

	if d.HasChange("enable_rbac_authorization") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		update.Properties.EnableRbacAuthorization = utils.Bool(d.Get("enable_rbac_authorization").(bool))
	}

	if d.HasChange("network_acls") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		networkAclsRaw := d.Get("network_acls").([]interface{})
		networkAcls, subnetIds := expandKeyVaultNetworkAcls(networkAclsRaw)

		// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
		virtualNetworkNames := make([]string, 0)
		for _, v := range subnetIds {
			id, err := commonids.ParseSubnetIDInsensitively(v)
			if err != nil {
				return err
			}

			if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
				virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
			}
		}

		locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

		update.Properties.NetworkAcls = networkAcls
	}

	if d.HasChange("purge_protection_enabled") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		newValue := d.Get("purge_protection_enabled").(bool)

		// existing.Properties guaranteed non-nil above
		oldValue := false
		if existing.Model.Properties.EnablePurgeProtection != nil {
			oldValue = *existing.Model.Properties.EnablePurgeProtection
		}

		// whilst this should have got caught in the customizeDiff this won't work if that fields interpolated
		// hence the double-checking here
		if oldValue && !newValue {
			return fmt.Errorf("updating %s: once Purge Protection has been Enabled it's not possible to disable it", *id)
		}

		update.Properties.EnablePurgeProtection = utils.Bool(newValue)

		if newValue {
			// When the KV was created with a version prior to v2.42 and the `soft_delete_enabled` is set to false, setting `purge_protection_enabled` to `true` would not work when updating KV with v2.42 or later of terraform provider.
			// This is because the `purge_protection_enabled` only works when soft delete is enabled.
			// Since version v2.42 of the Azure Provider and later force the value of `soft_delete_enabled` to be true, we should set `EnableSoftDelete` to true when `purge_protection_enabled` is enabled to make sure it works in this case.
			update.Properties.EnableSoftDelete = utils.Bool(true)
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if update.Properties == nil {
			update.Properties = &vaults.VaultPatchProperties{}
		}

		if isPublic {
			update.Properties.PublicNetworkAccess = utils.String("Enabled")
		} else {
			update.Properties.PublicNetworkAccess = utils.String("Disabled")
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
		var oldValue int64 = 0
		if existing.Model.Properties.SoftDeleteRetentionInDays != nil {
			oldValue = *existing.Model.Properties.SoftDeleteRetentionInDays
		}

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
		t := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(t)
	}

	if _, err := client.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("contact") {
		contacts := dataplane.Contacts{
			ContactList: expandKeyVaultCertificateContactList(d.Get("contact").(*pluginsdk.Set).List()),
		}

		vaultUri := ""
		if existing.Model != nil && existing.Model.Properties.VaultUri != nil {
			vaultUri = *existing.Model.Properties.VaultUri
		}

		if vaultUri == "" {
			return fmt.Errorf("failed to get vault base url for %s: %s", *id, err)
		}

		var err error
		if len(*contacts.ContactList) == 0 {
			_, err = managementClient.DeleteCertificateContacts(ctx, vaultUri)
		} else {
			_, err = managementClient.SetCertificateContacts(ctx, vaultUri, contacts)
		}

		if err != nil {
			var extendedErrorMsg string
			if !isPublic {
				extendedErrorMsg = "\n\nWARNING: public network access for this key vault has been disabled, access to the key vault is only allowed through private endpoints"
			}
			return fmt.Errorf("updating Contacts for %s: %+v %s", *id, err, extendedErrorMsg)
		}
	}

	d.Partial(false)

	return resourceKeyVaultRead(d, meta)
}

func resourceKeyVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	managementClient := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKeyVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.VaultUri != nil {
			vaultUri = *model.Properties.VaultUri
		}
	}

	if vaultUri != "" {
		meta.(*clients.Client).KeyVault.AddToCache(*id, vaultUri)
	}

	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("vault_uri", vaultUri)

	publicNetworkAccessEnabled := true

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("tenant_id", model.Properties.TenantId)
		d.Set("enabled_for_deployment", model.Properties.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", model.Properties.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", model.Properties.EnabledForTemplateDeployment)
		d.Set("enable_rbac_authorization", model.Properties.EnableRbacAuthorization)
		d.Set("purge_protection_enabled", model.Properties.EnablePurgeProtection)

		if model.Properties.PublicNetworkAccess != nil {
			publicNetworkAccessEnabled = strings.EqualFold(*model.Properties.PublicNetworkAccess, "Enabled")
		}
		d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

		// @tombuildsstuff: the API doesn't return this field if it's not configured
		// however https://docs.microsoft.com/en-us/azure/key-vault/general/soft-delete-overview
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

		flattenedPolicies := flattenAccessPolicies(model.Properties.AccessPolicies)
		if err := d.Set("access_policy", flattenedPolicies); err != nil {
			return fmt.Errorf("setting `access_policy`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

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
	contacts, err := managementClient.GetCertificateContacts(ctx, vaultUri)
	if err != nil {
		if publicNetworkAccessEnabled && (!utils.ResponseWasForbidden(contacts.Response) && !utils.ResponseWasNotFound(contacts.Response)) {
			return fmt.Errorf("retrieving `contact` for KeyVault: %+v", err)
		}
	}

	if err := d.Set("contact", flattenKeyVaultCertificateContactList(&contacts)); err != nil {
		return fmt.Errorf("setting `contact` for KeyVault: %+v", err)
	}

	return nil
}

func resourceKeyVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKeyVaultID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultName, keyVaultResourceName)

	read, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	location := ""
	purgeProtectionEnabled := false
	softDeleteEnabled := false
	virtualNetworkNames := make([]string, 0)
	if model := read.Model; model != nil {
		if model.Location != nil {
			location = *model.Location
		}

		// Check to see if purge protection is enabled or not...
		if model.Properties.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *model.Properties.EnablePurgeProtection
		}
		if model.Properties.EnableSoftDelete != nil {
			softDeleteEnabled = *model.Properties.EnableSoftDelete
		}

		// ensure we lock on the latest network names, to ensure we handle Azure's networking layer being limited to one change at a time
		if acls := model.Properties.NetworkAcls; acls != nil {
			if rules := acls.VirtualNetworkRules; rules != nil {
				for _, v := range *rules {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
					if err != nil {
						return err
					}

					if !utils.SliceContainsValue(virtualNetworkNames, subnetId.VirtualNetworkName) {
						virtualNetworkNames = append(virtualNetworkNames, subnetId.VirtualNetworkName)
					}
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Purge the soft deleted key vault permanently if the feature flag is enabled
	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeleteOnDestroy && softDeleteEnabled {
		deletedVaultId := vaults.NewDeletedVaultID(id.SubscriptionId, location, id.VaultName)

		// KeyVaults with Purge Protection Enabled cannot be deleted unless done by Azure
		if purgeProtectionEnabled {
			deletedInfo, err := getSoftDeletedStateForKeyVault(ctx, client, deletedVaultId)
			if err != nil {
				return fmt.Errorf("retrieving the Deletion Details for %s: %+v", *id, err)
			}

			// in the future it'd be nice to raise a warning, but this is the best we can do for now
			if deletedInfo != nil {
				log.Printf("[DEBUG] The Key Vault %q has Purge Protection Enabled and was deleted on %q. Azure will purge this on %q", id.VaultName, deletedInfo.deleteDate, deletedInfo.purgeDate)
			} else {
				log.Printf("[DEBUG] The Key Vault %q has Purge Protection Enabled and will be purged automatically by Azure", id.VaultName)
			}
			return nil
		}

		log.Printf("[DEBUG] KeyVault %q marked for purge - executing purge", id.VaultName)
		if err := client.PurgeDeletedThenPoll(ctx, deletedVaultId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
		log.Printf("[DEBUG] Purged KeyVault %q.", id.VaultName)
	}

	meta.(*clients.Client).KeyVault.Purge(*id)

	return nil
}

func keyVaultRefreshFunc(vaultUri string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if KeyVault %q is available..", vaultUri)

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}

		conn, err := client.Get(vaultUri)
		if err != nil {
			log.Printf("[DEBUG] Didn't find KeyVault at %q", vaultUri)
			return nil, "pending", fmt.Errorf("connecting to %q: %s", vaultUri, err)
		}

		defer conn.Body.Close()

		log.Printf("[DEBUG] Found KeyVault at %q", vaultUri)
		return "available", "available", nil
	}
}

func expandKeyVaultNetworkAcls(input []interface{}) (*vaults.NetworkRuleSet, []string) {
	subnetIds := make([]string, 0)
	if len(input) == 0 {
		return nil, subnetIds
	}

	v := input[0].(map[string]interface{})

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
		Bypass:              pointer.To(vaults.NetworkRuleBypassOptions(bypass)),
		DefaultAction:       pointer.To(vaults.NetworkRuleAction(defaultAction)),
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}
	return &ruleSet, subnetIds
}

// TODO: Remove in 4.0
func expandKeyVaultCertificateContactList(input []interface{}) *[]dataplane.Contact {
	results := make([]dataplane.Contact, 0)
	if len(input) == 0 || input[0] == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, dataplane.Contact{
			Name:         utils.String(v["name"].(string)),
			EmailAddress: utils.String(v["email"].(string)),
			Phone:        utils.String(v["phone"].(string)),
		})
	}

	return &results
}

func flattenKeyVaultNetworkAcls(input *vaults.NetworkRuleSet) []interface{} {
	bypass := string(vaults.NetworkRuleBypassOptionsAzureServices)
	defaultAction := string(vaults.NetworkRuleActionAllow)
	ipRules := make([]interface{}, 0)
	virtualNetworkSubnetIds := make([]interface{}, 0)

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

	return []interface{}{
		map[string]interface{}{
			"bypass":                     bypass,
			"default_action":             defaultAction,
			"ip_rules":                   pluginsdk.NewSet(pluginsdk.HashString, ipRules),
			"virtual_network_subnet_ids": pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkSubnetIds),
		},
	}
}

func flattenKeyVaultCertificateContactList(input *dataplane.Contacts) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.ContactList == nil {
		return results
	}

	for _, contact := range *input.ContactList {
		emailAddress := ""
		if contact.EmailAddress != nil {
			emailAddress = *contact.EmailAddress
		}

		name := ""
		if contact.Name != nil {
			name = *contact.Name
		}

		phone := ""
		if contact.Phone != nil {
			phone = *contact.Phone
		}

		results = append(results, map[string]interface{}{
			"email": emailAddress,
			"name":  name,
			"phone": phone,
		})
	}

	return results
}

func optedOutOfRecoveringSoftDeletedKeyVaultErrorFmt(name, location string) string {
	return fmt.Sprintf(`
An existing soft-deleted Key Vault exists with the Name %q in the location %q, however
automatically recovering this KeyVault has been disabled via the "features" block.

Terraform can automatically recover the soft-deleted Key Vault when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import", or pick a different name/location.
`, name, location)
}

type keyVaultDeletionStatus struct {
	deleteDate string
	purgeDate  string
}

func getSoftDeletedStateForKeyVault(ctx context.Context, client *vaults.VaultsClient, deletedVaultId vaults.DeletedVaultId) (*keyVaultDeletionStatus, error) {
	resp, err := client.GetDeleted(ctx, deletedVaultId)
	if err != nil {
		return nil, err
	}

	if model := resp.Model; model != nil {
		// the logic is this way because the GetDeleted call will return an existing key vault
		// that is not soft deleted, but the Deleted Vault properties will be nil
		result := keyVaultDeletionStatus{}
		if props := model.Properties; props != nil {
			if props.DeletionDate != nil {
				t, err := props.GetDeletionDateAsTime()
				if err != nil {
					return nil, fmt.Errorf("parsing `DeletionDate`: %+v", err)
				}
				result.deleteDate = t.Format(time.RFC3339)
			}
			if props.ScheduledPurgeDate != nil {
				t, err := props.GetScheduledPurgeDateAsTime()
				if err != nil {
					return nil, fmt.Errorf("parsing `ScheduledPurgeDate`: %+v", err)
				}
				result.purgeDate = t.Format(time.RFC3339)
			}
		}

		if result.deleteDate != "" && result.purgeDate != "" {
			return &result, nil
		}
	}

	// otherwise we've found an existing key vault that is not soft deleted
	return nil, nil
}
