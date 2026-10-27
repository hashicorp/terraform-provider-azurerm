// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesVaultRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	vaultSettingsClient := meta.(*clients.Client).RecoveryServices.VaultsSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

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

		flattenIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		encryption := flattenVaultEncryption(*model)
		if encryption != nil {
			d.Set("encryption", []any{encryption})
		}

		vaultSettingsId := replicationvaultsetting.NewReplicationVaultSettingID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, "default")
		vaultSetting, err := vaultSettingsClient.Get(ctx, vaultSettingsId)
		if err != nil {
			return fmt.Errorf("reading Recovery Service Vault Setting %s: %+v", id, err)
		}

		classicVmwareReplicationEnabled := false
		if vaultSetting.Model != nil && vaultSetting.Model.Properties != nil {
			if v := vaultSetting.Model.Properties.VMwareToAzureProviderType; v != nil {
				classicVmwareReplicationEnabled = strings.EqualFold(*v, "vmware")
			}
		}
		d.Set("classic_vmware_replication_enabled", classicVmwareReplicationEnabled)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
