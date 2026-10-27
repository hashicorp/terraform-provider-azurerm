// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesVaultUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	encryption, err := expandEncryption(d)
	if err != nil {
		return err
	}
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: `model` was nil", id)
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

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	if model.Identity != nil && !validateIdentityUpdate(*existing.Model.Identity, *expandedIdentity) {
		return fmt.Errorf("`Once `identity` specified, the managed identity must not be disabled (even temporarily). Disabling the managed identity may lead to inconsistent behaviour. Details could be found on https://learn.microsoft.com/azure/backup/encryption-at-rest-with-cmk?tabs=portal#enable-system-assigned-managed-identity-for-the-vault")
	}

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(vaults.StandardTierStorageRedundancyGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(vaults.StandardTierStorageRedundancyGeoRedundant), id)
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
				MonitoringSettings:  expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]any)),
			},
		}

		if vaults.SkuName(sku) == vaults.SkuNameRSZero {
			vault.Sku.Tier = pointer.To("Standard")
		}

		if err = client.CreateOrUpdateThenPoll(ctx, id, vault, vaults.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v", id, err)
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
		vault.Properties.MonitoringSettings = expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]any))
	}

	if d.HasChange("identity") {
		vault.Identity = expandedIdentity
	}

	if d.HasChange("encryption") {
		vault.Properties.Encryption = encryption
	}

	if d.HasChange("tags") {
		vault.Tags = tags.Expand(d.Get("tags").(map[string]any))
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
			StandardTierStorageRedundancy: pointer.ToEnum[vaults.StandardTierStorageRedundancy](storageMode),
		}
	}

	if err = client.UpdateThenPoll(ctx, id, vault, vaults.DefaultUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating  %s: %+v", id, err)
	}

	if requireAdditionalUpdate {
		if err := client.UpdateThenPoll(ctx, id, additionalUpdatePatch, vaults.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v, but recovery vault was created, a manually import might be required", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceRecoveryServicesVaultRead(d, meta)
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
