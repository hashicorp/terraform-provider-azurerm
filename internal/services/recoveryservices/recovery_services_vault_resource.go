// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupresourcevaultconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesVaultCreate,
		Read:   resourceRecoveryServicesVaultRead,
		Update: resourceRecoveryServicesVaultUpdate,
		Delete: resourceRecoveryServicesVaultDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := vaults.ParseVaultID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"encryption": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				RequiredWith: []string{"identity"},
				MaxItems:     1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
						},
						"infrastructure_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"use_system_assigned_identity": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// set `immutability` to Computed, because it will start to return from the service once it has been set.
			"immutability": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(vaults.ImmutabilityStateLocked),
					string(vaults.ImmutabilityStateUnlocked),
					string(vaults.ImmutabilityStateDisabled),
				}, false),
			},

			"tags": commonschema.Tags(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(vaults.SkuNameRSZero),
					string(vaults.SkuNameStandard),
				}, false),
			},

			"storage_mode_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  vaults.StandardTierStorageRedundancyGeoRedundant,
				ValidateFunc: validation.StringInSlice([]string{
					string(vaults.StandardTierStorageRedundancyGeoRedundant),
					string(vaults.StandardTierStorageRedundancyLocallyRedundant),
					string(vaults.StandardTierStorageRedundancyZoneRedundant),
				}, false),
			},

			"cross_region_restore_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"soft_delete_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"monitoring": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"alerts_for_all_job_failures_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"alerts_for_critical_operation_failures_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"classic_vmware_replication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && !new.(bool)
			}),
			pluginsdk.ForceNewIfChange("immutability", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) == string(vaults.ImmutabilityStateLocked)
			}),
		),
	}
}

func resourceRecoveryServicesVaultCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	settingsClient := meta.(*clients.Client).RecoveryServices.VaultsSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	cfgId := backupresourcevaultconfigs.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(vaults.StandardTierStorageRedundancyGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(vaults.StandardTierStorageRedundancyGeoRedundant), id.String())
	}

	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating Recovery Service %s", id.String())

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_recovery_services_vault", id.ID())
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	sku := d.Get("sku").(string)
	crossRegionRestoreEnabled := vaults.CrossRegionRestoreDisabled
	if d.Get("cross_region_restore_enabled").(bool) {
		crossRegionRestoreEnabled = vaults.CrossRegionRestoreEnabled
	}

	vault := vaults.Vault{
		Location: location,
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
		Sku: &vaults.Sku{
			Name: vaults.SkuName(sku),
		},
		Properties: &vaults.VaultProperties{
			PublicNetworkAccess: expandRecoveryServicesVaultPublicNetworkAccess(d.Get("public_network_access_enabled").(bool)),
			MonitoringSettings:  expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]interface{})),
			RedundancySettings: &vaults.VaultPropertiesRedundancySettings{
				CrossRegionRestore:            &crossRegionRestoreEnabled,
				StandardTierStorageRedundancy: pointer.To(vaults.StandardTierStorageRedundancy(d.Get("storage_mode_type").(string))),
			},
		},
	}

	if vaults.SkuName(sku) == vaults.SkuNameRSZero {
		vault.Sku.Tier = utils.String("Standard")
	}

	if _, ok := d.GetOk("encryption"); ok {
		encryption, err := expandEncryption(d)
		if err != nil {
			return err
		}
		vault.Properties.Encryption = encryption
	}

	requireAdditionalUpdate := false
	updatePatch := vaults.PatchVault{
		Properties: &vaults.VaultProperties{},
	}
	if immutability, ok := d.GetOk("immutability"); ok {
		// The API doesn't allow to set the immutability to "Locked" on creation.
		// Here we firstly make it "Unlocked", and once created, we will update it to "Locked".
		// Note: The `immutability` could be transitioned only in the limited directions.
		// Locked <- Unlocked <-> Disabled
		if immutability == string(vaults.ImmutabilityStateLocked) {
			updatePatch.Properties.SecuritySettings = expandRecoveryServicesVaultSecuritySettings(immutability)
			requireAdditionalUpdate = true
			immutability = string(vaults.ImmutabilityStateUnlocked)
		}
		vault.Properties.SecuritySettings = expandRecoveryServicesVaultSecuritySettings(immutability)
	}

	// Async Operaation of creation with `UserAssigned` identity is returned with 404
	// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/27869
	// `SystemAssigned, UserAssigned` Identity require an additional update to work
	// Trakced on https://github.com/Azure/azure-rest-api-specs/issues/27851
	if expandedIdentity.Type == identity.TypeUserAssigned || expandedIdentity.Type == identity.TypeSystemAssignedUserAssigned {
		requireAdditionalUpdate = true
		updatePatch.Identity = expandedIdentity
		vault.Identity = &identity.SystemAndUserAssignedMap{
			Type: identity.TypeNone,
		}
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, vault)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id.String(), err)
	}

	if requireAdditionalUpdate {
		err := client.UpdateThenPoll(ctx, id, updatePatch)
		if err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}

	}
	// an update on the vault will reset the vault config to default, so we handle it at last.
	enhancedSecurityState := backupresourcevaultconfigs.EnhancedSecurityStateEnabled
	cfg := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
		Properties: &backupresourcevaultconfigs.BackupResourceVaultConfig{
			EnhancedSecurityState: &enhancedSecurityState, // always enabled
		},
	}

	var StateRefreshPendingStrings []string
	var StateRefreshTargetStrings []string
	if sd := d.Get("soft_delete_enabled").(bool); sd {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
	} else {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
	}

	_, err = cfgsClient.Update(ctx, cfgId, cfg)
	if err != nil {
		return err
	}

	// sometimes update sync succeed but READ returns with old value, so we refresh till the value is correct.
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   StateRefreshPendingStrings,
		Target:                    StateRefreshTargetStrings,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 3,
		Refresh:                   resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx, cfgsClient, cfgId),
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service %s: %+v", id.String(), err)
	}

	if d.Get("classic_vmware_replication_enabled").(bool) {
		settingsId := replicationvaultsetting.NewReplicationVaultSettingID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, "default")
		settingsInput := replicationvaultsetting.VaultSettingCreationInput{
			Properties: replicationvaultsetting.VaultSettingCreationInputProperties{
				VMwareToAzureProviderType: utils.String("Vmware"),
			},
		}
		if err := settingsClient.CreateThenPoll(ctx, settingsId, settingsInput); err != nil {
			return fmt.Errorf("creating %s: %+v", settingsId, err)
		}
	}

	d.SetId(id.ID())
	return resourceRecoveryServicesVaultRead(d, meta)
}

func resourceRecoveryServicesVaultUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	cfgId := backupresourcevaultconfigs.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

	encryption, err := expandEncryption(d)
	if err != nil {
		return err
	}
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
	}
	if existing.Model == nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: `model` was nil", id.String())
	}
	model := existing.Model

	if model.Properties != nil {
		prop := model.Properties
		if prop.Encryption != nil {
			if encryption == nil {
				return fmt.Errorf("once encryption with your own key has been enabled it's not possible to disable it")
			}
			if *encryption.InfrastructureEncryption != *prop.Encryption.InfrastructureEncryption {
				return fmt.Errorf("once `infrastructure_encryption_enabled` has been set it's not possible to change it")
			}
			if d.HasChange("sku") {
				// Once encryption has been enabled, calling `CreateOrUpdate` without it is not allowed.
				// But `sku` can only be updated by `CreateOrUpdate` and the support for `encryption` in `CreateOrUpdate` is still under preview (https://docs.microsoft.com/azure/backup/encryption-at-rest-with-cmk?tabs=portal#enable-encryption-using-customer-managed-keys-at-vault-creation-in-preview).
				// TODO remove this restriction and add `encryption` to below `sku` update block when `encryption` in `CreateOrUpdate` is GA
				return fmt.Errorf("`sku` cannot be changed when encryption with your own key has been enabled")
			}
		}
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	if model.Identity != nil && !validateIdentityUpdate(*existing.Model.Identity, *expandedIdentity) {
		return fmt.Errorf("`Once `identity` specified, the managed identity must not be disabled (even temporarily). Disabling the managed identity may lead to inconsistent behavior. Details could be found on https://learn.microsoft.com/en-us/azure/backup/encryption-at-rest-with-cmk?tabs=portal#enable-system-assigned-managed-identity-for-the-vault")
	}

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(vaults.StandardTierStorageRedundancyGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(vaults.StandardTierStorageRedundancyGeoRedundant), id.String())
	}

	enhanchedSecurityState := backupresourcevaultconfigs.EnhancedSecurityStateEnabled
	cfg := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
		Properties: &backupresourcevaultconfigs.BackupResourceVaultConfig{
			EnhancedSecurityState: &enhanchedSecurityState, // always enabled
		},
	}

	// `sku` can only be updated by `CreateOrUpdate` but not `Update`, so use `CreateOrUpdate` with required and unchangeable properties
	if d.HasChange("sku") {
		sku := d.Get("sku").(string)
		vault := vaults.Vault{
			Location: d.Get("location").(string),
			Identity: expandedIdentity,
			Sku: &vaults.Sku{
				Name: vaults.SkuName(sku),
			},
			Properties: &vaults.VaultProperties{
				PublicNetworkAccess: expandRecoveryServicesVaultPublicNetworkAccess(d.Get("public_network_access_enabled").(bool)), // It's required to call CreateOrUpdate.
				MonitoringSettings:  expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]interface{})),
			},
		}

		if vaults.SkuName(sku) == vaults.SkuNameRSZero {
			vault.Sku.Tier = utils.String("Standard")
		}

		err = client.CreateOrUpdateThenPoll(ctx, id, vault)
		if err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v", id.String(), err)
		}
	}

	requireAdditionalUpdate := false
	additionalUpdatePatch := vaults.PatchVault{
		Properties: &vaults.VaultProperties{},
	}
	vault := vaults.PatchVault{
		Properties: &vaults.VaultProperties{},
	}

	if d.HasChange("public_network_access_enabled") {
		vault.Properties.PublicNetworkAccess = expandRecoveryServicesVaultPublicNetworkAccess(d.Get("public_network_access_enabled").(bool))
	}

	if d.HasChanges("monitoring") {
		vault.Properties.MonitoringSettings = expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]interface{}))
	}

	if d.HasChange("identity") {
		vault.Identity = expandedIdentity
	}

	if d.HasChange("encryption") {
		vault.Properties.Encryption = encryption
	}

	if d.HasChange("tags") {
		vault.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("immutability") {
		// The API does not allow to set the immutability from `Disabled` to `Locked` directly,
		// Hence we firstly make it `Unlocked`, and once created, we will update it to `Locked`.
		// Note: The `immutability` could be transitioned only in the limited directions.
		// Locked <- Unlocked <-> Disabled

		// When the service returns `null`, it equals `disabled`
		currentImmutability := pointer.To(vaults.ImmutabilityStateDisabled)
		if model.Properties != nil && model.Properties.SecuritySettings != nil && model.Properties.SecuritySettings.ImmutabilitySettings != nil && model.Properties.SecuritySettings.ImmutabilitySettings.State != nil {
			currentImmutability = model.Properties.SecuritySettings.ImmutabilitySettings.State
		}
		immutability := d.Get("immutability")
		if string(*currentImmutability) == string(vaults.ImmutabilityStateDisabled) && immutability == string(vaults.ImmutabilityStateLocked) {
			additionalUpdatePatch.Properties.SecuritySettings = expandRecoveryServicesVaultSecuritySettings(immutability)
			requireAdditionalUpdate = true
			immutability = string(vaults.ImmutabilityStateUnlocked)
		}
		vault.Properties.SecuritySettings = expandRecoveryServicesVaultSecuritySettings(immutability)
	}

	crossRegionRestoreEnabled := vaults.CrossRegionRestoreDisabled
	if crossRegionRestore {
		crossRegionRestoreEnabled = vaults.CrossRegionRestoreEnabled
	}

	if d.HasChanges("storage_mode_type", "cross_region_restore_enabled") {
		vault.Properties.RedundancySettings = &vaults.VaultPropertiesRedundancySettings{
			CrossRegionRestore:            &crossRegionRestoreEnabled,
			StandardTierStorageRedundancy: pointer.To(vaults.StandardTierStorageRedundancy(storageMode)),
		}
	}

	err = client.UpdateThenPoll(ctx, id, vault)
	if err != nil {
		return fmt.Errorf("updating  %s: %+v", id, err)
	}

	if requireAdditionalUpdate {
		err := client.UpdateThenPoll(ctx, id, additionalUpdatePatch)
		if err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}
	}

	// an update on vault will cause the vault config reset to default, so whether the config has change or not, it needs to be updated.
	var StateRefreshPendingStrings []string
	var StateRefreshTargetStrings []string
	if sd := d.Get("soft_delete_enabled").(bool); sd {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
	} else {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
	}

	_, err = cfgsClient.Update(ctx, cfgId, cfg)
	if err != nil {
		return err
	}

	// sometimes update sync succeed but READ returns with old value, so we refresh till the value is correct.
	// tracked by https://github.com/Azure/azure-rest-api-specs/issues/21548
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   StateRefreshPendingStrings,
		Target:                    StateRefreshTargetStrings,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 3,
		Refresh:                   resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx, cfgsClient, cfgId),
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())
	return resourceRecoveryServicesVaultRead(d, meta)
}

func resourceRecoveryServicesVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	vaultSettingsClient := meta.(*clients.Client).RecoveryServices.VaultsSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

	cfgId := backupresourcevaultconfigs.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
		}

		if prop := model.Properties; prop != nil {

			immutability := vaults.ImmutabilityStateDisabled
			if prop.SecuritySettings != nil && prop.SecuritySettings.ImmutabilitySettings != nil {
				immutability = pointer.From(prop.SecuritySettings.ImmutabilitySettings.State)
			}
			d.Set("immutability", string(immutability))

			d.Set("public_network_access_enabled", flattenRecoveryServicesVaultPublicNetworkAccess(model.Properties.PublicNetworkAccess))

			d.Set("monitoring", flattenRecoveryServicesVaultMonitorSettings(prop.MonitoringSettings))

			storageModeType := vaults.StandardTierStorageRedundancyInvalid
			crossRegionRestoreEnabled := false
			if prop.RedundancySettings != nil {
				storageModeType = pointer.From(prop.RedundancySettings.StandardTierStorageRedundancy)
				if prop.RedundancySettings.CrossRegionRestore != nil {
					crossRegionRestoreEnabled = *prop.RedundancySettings.CrossRegionRestore == vaults.CrossRegionRestoreEnabled
				}
			}
			d.Set("cross_region_restore_enabled", crossRegionRestoreEnabled)
			d.Set("storage_mode_type", string(storageModeType))
		}

		cfg, err := cfgsClient.Get(ctx, cfgId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", cfgId, err)
		}

		softDeleteEnabled := false
		if cfg.Model != nil && cfg.Model.Properties != nil && cfg.Model.Properties.SoftDeleteFeatureState != nil {
			softDeleteEnabled = *cfg.Model.Properties.SoftDeleteFeatureState == backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled
		}

		d.Set("soft_delete_enabled", softDeleteEnabled)

		flattenIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		encryption := flattenVaultEncryption(*model)
		if encryption != nil {
			d.Set("encryption", []interface{}{encryption})
		}

		vaultSettingsId := replicationvaultsetting.NewReplicationVaultSettingID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, "default")
		vaultSetting, err := vaultSettingsClient.Get(ctx, vaultSettingsId)
		if err != nil {
			return fmt.Errorf("reading Recovery Service Vault Setting %s: %+v", id.String(), err)
		}

		classicVmwareReplicationEnabled := false
		if vaultSetting.Model != nil && vaultSetting.Model.Properties != nil {
			if v := vaultSetting.Model.Properties.VMwareToAzureProviderType; v != nil {
				classicVmwareReplicationEnabled = strings.EqualFold(*v, "vmware")
			}
		}
		d.Set("classic_vmware_replication_enabled", classicVmwareReplicationEnabled)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceRecoveryServicesVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	protectedItemsClient := meta.(*clients.Client).RecoveryServices.ProtectedItemsGroupClient
	protectedItemClient := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opResultClient := meta.(*clients.Client).RecoveryServices.BackupOperationResultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

	if meta.(*clients.Client).Features.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy {
		log.Printf("[DEBUG] Purging Protected Items from %s", id.String())

		vaultId := backupprotecteditems.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

		protectedItems, err := protectedItemsClient.ListComplete(ctx, vaultId, backupprotecteditems.ListOperationOptions{})
		if err != nil {
			return fmt.Errorf("listing protected items in %s: %+v", id, err)
		}

		for _, item := range protectedItems.Items {
			if item.Id != nil {
				protectedItemId, err := protecteditems.ParseProtectedItemID(pointer.From(item.Id))
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Purging %s from %s", protectedItemId, id)

				resp, err := protectedItemClient.Delete(ctx, *protectedItemId)
				if err != nil {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("issuing delete request for %s: %+v", protectedItemId, err)
					}
				}

				operationId, err := parseBackupOperationId(resp.HttpResponse)
				if err != nil {
					return fmt.Errorf("purging %s from %s: %+v", protectedItemId, id, err)
				}

				if err = resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx, protectedItemClient, opResultClient, *protectedItemId, operationId); err != nil {
					return fmt.Errorf("waiting for %s to be purged from %s: %+v", protectedItemId, id, err)
				}
			}
		}
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.String(), err)
	}

	return nil
}

func validateIdentityUpdate(origin identity.SystemAndUserAssignedMap, target identity.SystemAndUserAssignedMap) bool {
	switch origin.Type {
	case identity.TypeSystemAssigned:
		switch target.Type {
		case identity.TypeNone:
			return false
		case identity.TypeUserAssigned:
			return false
		default:
			return true
		}
	case identity.TypeUserAssigned:
		switch target.Type {
		case identity.TypeNone:
			return false
		case identity.TypeSystemAssigned:
			return false
		default:
			return true
		}
	case identity.TypeSystemAssignedUserAssigned:
		switch target.Type {
		case identity.TypeNone:
			return false
		case identity.TypeSystemAssigned:
			return false
		case identity.TypeUserAssigned:
			return false
		default:
			return true
		}
	}
	return true
}

func expandEncryption(d *pluginsdk.ResourceData) (*vaults.VaultPropertiesEncryption, error) {
	encryptionRaw := d.Get("encryption")
	if encryptionRaw == nil {
		return nil, nil
	}
	settings := encryptionRaw.([]interface{})
	if len(settings) == 0 {
		return nil, nil
	}
	encryptionMap := settings[0].(map[string]interface{})
	keyUri := encryptionMap["key_id"].(string)
	enabledInfraEncryption := encryptionMap["infrastructure_encryption_enabled"].(bool)
	infraEncryptionState := vaults.InfrastructureEncryptionStateEnabled
	if !enabledInfraEncryption {
		infraEncryptionState = vaults.InfrastructureEncryptionStateDisabled
	}
	encryption := &vaults.VaultPropertiesEncryption{
		KeyVaultProperties: &vaults.CmkKeyVaultProperties{
			KeyUri: utils.String(keyUri),
		},
		KekIdentity: &vaults.CmkKekIdentity{
			UseSystemAssignedIdentity: utils.Bool(encryptionMap["use_system_assigned_identity"].(bool)),
		},
		InfrastructureEncryption: &infraEncryptionState,
	}
	if v, ok := encryptionMap["user_assigned_identity_id"].(string); ok && v != "" {
		if *encryption.KekIdentity.UseSystemAssignedIdentity {
			return nil, fmt.Errorf(" `use_system_assigned_identity` must be disabled when `user_assigned_identity_id` is set.")
		}
		encryption.KekIdentity.UserAssignedIdentity = utils.String(v)
	}
	return encryption, nil
}

func flattenVaultEncryption(model vaults.Vault) interface{} {
	if model.Properties == nil || model.Properties.Encryption == nil {
		return nil
	}
	encryption := model.Properties.Encryption
	if encryption.KeyVaultProperties == nil || encryption.KeyVaultProperties.KeyUri == nil {
		return nil
	}
	if encryption.KekIdentity == nil || encryption.KekIdentity.UseSystemAssignedIdentity == nil {
		return nil
	}
	encryptionMap := make(map[string]interface{})
	encryptionMap["key_id"] = encryption.KeyVaultProperties.KeyUri
	encryptionMap["use_system_assigned_identity"] = *encryption.KekIdentity.UseSystemAssignedIdentity
	encryptionMap["infrastructure_encryption_enabled"] = *encryption.InfrastructureEncryption == vaults.InfrastructureEncryptionStateEnabled
	if encryption.KekIdentity.UserAssignedIdentity != nil {
		encryptionMap["user_assigned_identity_id"] = *encryption.KekIdentity.UserAssignedIdentity
	}
	return encryptionMap
}

func expandRecoveryServicesVaultSecuritySettings(input interface{}) *vaults.SecuritySettings {
	if input == nil || len(input.(string)) == 0 {
		return nil
	}
	immutabilityState := vaults.ImmutabilityState(input.(string))
	return &vaults.SecuritySettings{
		ImmutabilitySettings: &vaults.ImmutabilitySettings{
			State: &immutabilityState,
		},
	}
}

func expandRecoveryServicesVaultPublicNetworkAccess(input bool) *vaults.PublicNetworkAccess {
	out := vaults.PublicNetworkAccessDisabled
	if input {
		out = vaults.PublicNetworkAccessEnabled
	}
	return &out
}

func flattenRecoveryServicesVaultPublicNetworkAccess(input *vaults.PublicNetworkAccess) bool {
	if input == nil {
		return false
	}
	return *input == vaults.PublicNetworkAccessEnabled
}

func expandRecoveryServicesVaultMonitorSettings(input []interface{}) *vaults.MonitoringSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	allJobAlert := vaults.AlertsStateDisabled
	if v["alerts_for_all_job_failures_enabled"].(bool) {
		allJobAlert = vaults.AlertsStateEnabled
	}

	criticalOperation := vaults.AlertsStateDisabled
	if v["alerts_for_critical_operation_failures_enabled"].(bool) {
		criticalOperation = vaults.AlertsStateEnabled
	}

	return pointer.To(vaults.MonitoringSettings{
		AzureMonitorAlertSettings: pointer.To(vaults.AzureMonitorAlertSettings{
			AlertsForAllJobFailures: pointer.To(allJobAlert),
		}),
		ClassicAlertSettings: pointer.To(vaults.ClassicAlertSettings{
			AlertsForCriticalOperations: pointer.To(criticalOperation),
		}),
	})
}

func flattenRecoveryServicesVaultMonitorSettings(input *vaults.MonitoringSettings) []interface{} {
	// `Monitor` is an optional parameters, and won't be returned from API if it has not been specified.
	if input == nil {
		return []interface{}{}
	}
	allJobAlert := false
	criticalAlert := false

	if input != nil {
		if input.AzureMonitorAlertSettings != nil && input.AzureMonitorAlertSettings.AlertsForAllJobFailures != nil {
			allJobAlert = *input.AzureMonitorAlertSettings.AlertsForAllJobFailures == vaults.AlertsStateEnabled
		}
		if input.ClassicAlertSettings != nil && input.ClassicAlertSettings.AlertsForCriticalOperations != nil {
			criticalAlert = *input.ClassicAlertSettings.AlertsForCriticalOperations == vaults.AlertsStateEnabled
		}
	}

	return []interface{}{
		map[string]interface{}{
			"alerts_for_all_job_failures_enabled":            allJobAlert,
			"alerts_for_critical_operation_failures_enabled": criticalAlert,
		},
	}
}

func resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx context.Context, cfgsClient *backupresourcevaultconfigs.BackupResourceVaultConfigsClient, id backupresourcevaultconfigs.VaultId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := cfgsClient.Get(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotYetSynced") {
				return resp, "syncing", nil
			}
			return resp, "error", fmt.Errorf("refreshing Recovery Service Vault Cfg %s: %+v", id.String(), err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.SoftDeleteFeatureState != nil {
			return resp.Model, string(*resp.Model.Properties.SoftDeleteFeatureState), nil
		}

		return resp, "error", fmt.Errorf("refreshing Recovery Service Vault Cfg %s: Properties is nil", id.String())
	}
}
