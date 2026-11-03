// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppVolumeUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	update := volumes.VolumePatch{
		Properties: &volumes.VolumePatchProperties{},
	}

	if d.HasChange("zones") {
		return fmt.Errorf("zone changes are not supported after volume is already created, %s", id)
	}

	if d.HasChange("snapshot_directory_visible") {
		snapshotDirectoryVisible := d.Get("snapshot_directory_visible").(bool)
		update.Properties.SnapshotDirectoryVisible = pointer.To(snapshotDirectoryVisible)
	}

	if d.HasChange("storage_quota_in_gb") {
		storageQuotaInBytes := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)
		update.Properties.UsageThreshold = pointer.To(storageQuotaInBytes)
	}

	if d.HasChange("export_policy_rule") {
		exportPolicyRuleRaw := d.Get("export_policy_rule").([]any)
		var protocolOverride []string
		// Only override export policy protocols if we're also changing volume protocols
		if d.HasChange("protocols") {
			protocols := d.Get("protocols").(*pluginsdk.Set).List()
			protocolOverride = *pluginsdk.ExpandStringSlice(protocols)
		}
		update.Properties.ExportPolicy = expandNetAppVolumeExportPolicyRulePatch(exportPolicyRuleRaw, protocolOverride)
	}

	if d.HasChange("protocols") {
		protocols := d.Get("protocols").(*pluginsdk.Set).List()
		update.Properties.ProtocolTypes = pluginsdk.ExpandStringSlice(protocols)
	}

	if d.HasChange("data_protection_snapshot_policy") {
		// Validating that snapshot policies are not being created in a data protection volume
		dataProtectionReplicationRaw := d.Get("data_protection_replication").([]any)
		dataProtectionReplication := helpers.ExpandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

		if dataProtectionReplication != nil && dataProtectionReplication.EndpointType != nil && strings.EqualFold(string(*dataProtectionReplication.EndpointType), "dst") {
			return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, %s", id)
		}

		dataProtectionSnapshotPolicyRaw := d.Get("data_protection_snapshot_policy").([]any)

		update.Properties.DataProtection = &volumes.VolumePatchPropertiesDataProtection{}
		update.Properties.DataProtection.Snapshot = helpers.ExpandNetAppVolumeDataProtectionSnapshotPolicyPatch(dataProtectionSnapshotPolicyRaw)
	}

	if d.HasChange("data_protection_backup_policy") {
		// Validate applicability of backup policies
		dataProtectionReplicationRaw := d.Get("data_protection_replication").([]any)
		dataProtectionReplication := helpers.ExpandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

		if dataProtectionReplication != nil && dataProtectionReplication.EndpointType != nil && strings.EqualFold(string(*dataProtectionReplication.EndpointType), "dst") {
			return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, %s", id)
		}

		dataProtectionBackupPolicyRaw := d.Get("data_protection_backup_policy").([]any)

		if update.Properties.DataProtection == nil {
			update.Properties.DataProtection = &volumes.VolumePatchPropertiesDataProtection{}
		}
		update.Properties.DataProtection.Backup = expandNetAppVolumeDataProtectionBackupPolicyPatch(dataProtectionBackupPolicyRaw)
	}

	if d.HasChange("data_protection_advanced_ransomware") {
		dataProtectionARPRaw := d.Get("data_protection_advanced_ransomware").([]any)

		if update.Properties.DataProtection == nil {
			update.Properties.DataProtection = &volumes.VolumePatchPropertiesDataProtection{}
		}
		update.Properties.DataProtection.RansomwareProtection = expandNetAppVolumeDataProtectionAdvancedRansomwareProtectionPatch(dataProtectionARPRaw)
	}

	if d.HasChange("throughput_in_mibps") {
		throughputMibps := d.Get("throughput_in_mibps")
		update.Properties.ThroughputMibps = pointer.To(throughputMibps.(float64))
	}

	if d.HasChange("smb_non_browsable_enabled") {
		update.Properties.SmbNonBrowsable = pointer.To(volumes.SmbNonBrowsableDisabled)
		if d.Get("smb_non_browsable_enabled").(bool) {
			update.Properties.SmbNonBrowsable = pointer.To(volumes.SmbNonBrowsableEnabled)
		}
	}

	if d.HasChange("smb_access_based_enumeration_enabled") {
		update.Properties.SmbAccessBasedEnumeration = pointer.To(volumes.SmbAccessBasedEnumerationDisabled)
		if d.Get("smb_access_based_enumeration_enabled").(bool) {
			update.Properties.SmbAccessBasedEnumeration = pointer.To(volumes.SmbAccessBasedEnumerationEnabled)
		}
	}

	if d.HasChange("cool_access") {
		if len(d.Get("cool_access").([]any)) > 0 {
			coolAccess := d.Get("cool_access").([]any)[0].(map[string]any)
			update.Properties.CoolAccess = pointer.To(true)

			if d.HasChange("cool_access.0.retrieval_policy") {
				update.Properties.CoolAccessRetrievalPolicy = pointer.ToEnum[volumes.CoolAccessRetrievalPolicy](coolAccess["retrieval_policy"].(string))
			}

			if d.HasChange("cool_access.0.tiering_policy") {
				update.Properties.CoolAccessTieringPolicy = pointer.ToEnum[volumes.CoolAccessTieringPolicy](coolAccess["tiering_policy"].(string))
			}

			if d.HasChange("cool_access.0.coolness_period_in_days") {
				update.Properties.CoolnessPeriod = pointer.To(int64(coolAccess["coolness_period_in_days"].(int)))
			}
		} else {
			update.Properties.CoolAccess = pointer.To(false)
		}
	}

	if d.HasChange("unix_permissions") {
		update.Properties.UnixPermissions = pointer.To(d.Get("unix_permissions").(string))
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]any)
		update.Tags = tags.Expand(tagsRaw)
	}

	if err = client.UpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Wait for volume to complete update
	if err := helpers.WaitForVolumeCreateOrUpdate(ctx, client, *id); err != nil {
		return err
	}

	if d.HasChanges("service_level", "pool_name") {
		poolName := d.Get("pool_name").(string)
		poolId := volumes.NewCapacityPoolID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, poolName)
		volumeId, err := volumes.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}
		poolChangeInput := volumes.PoolChangeRequest{
			NewPoolResourceId: poolId.ID(),
		}
		if _, err = client.PoolChange(ctx, *volumeId, poolChangeInput); err != nil {
			return fmt.Errorf("updating `service_level` for %s: %+v", id, err)
		}

		// the id needs to be updated as the pool is different
		newId := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, poolName, id.VolumeName)

		if err = helpers.WaitForVolumeCreateOrUpdate(ctx, client, newId); err != nil {
			return err
		}

		d.SetId(newId.ID())
	}

	return resourceNetAppVolumeRead(d, meta)
}
