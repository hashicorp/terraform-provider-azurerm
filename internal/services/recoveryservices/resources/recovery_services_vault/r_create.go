// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesVaultCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	settingsClient := meta.(*clients.Client).RecoveryServices.VaultsSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(vaults.StandardTierStorageRedundancyGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(vaults.StandardTierStorageRedundancyGeoRedundant), id)
	}

	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]any)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_recovery_services_vault", id.ID())
		}
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
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
			MonitoringSettings:  expandRecoveryServicesVaultMonitorSettings(d.Get("monitoring").([]any)),
			RedundancySettings: &vaults.VaultPropertiesRedundancySettings{
				CrossRegionRestore:            &crossRegionRestoreEnabled,
				StandardTierStorageRedundancy: pointer.ToEnum[vaults.StandardTierStorageRedundancy](d.Get("storage_mode_type").(string)),
			},
		},
	}

	if vaults.SkuName(sku) == vaults.SkuNameRSZero {
		vault.Sku.Tier = pointer.To("Standard")
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

	// Async Operation of creation with `UserAssigned` identity is returned with 404
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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, vault, vaults.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	if requireAdditionalUpdate {
		if err := client.UpdateThenPoll(ctx, id, updatePatch, vaults.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("creating Recovery Service %s: %+v", id, err)
		}
	}

	if d.Get("classic_vmware_replication_enabled").(bool) {
		settingsId := replicationvaultsetting.NewReplicationVaultSettingID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, "default")
		settingsInput := replicationvaultsetting.VaultSettingCreationInput{
			Properties: replicationvaultsetting.VaultSettingCreationInputProperties{
				VMwareToAzureProviderType: pointer.To("Vmware"),
			},
		}
		if err := settingsClient.CreateThenPoll(ctx, settingsId, settingsInput); err != nil {
			return fmt.Errorf("creating %s: %+v", settingsId, err)
		}
	}

	return resourceRecoveryServicesVaultRead(d, meta)
}
