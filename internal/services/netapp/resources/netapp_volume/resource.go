// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppVolumeCreate,
		Read:   resourceNetAppVolumeRead,
		Update: resourceNetAppVolumeUpdate,
		Delete: resourceNetAppVolumeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := volumes.ParseVolumeID(id)
			return err
		}),

		Schema: netappVolumeSchema(),
		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			// Validate large volume and storage_quota_in_gb based on Azure NetApp Files requirements
			isLargeVolume := d.Get("large_volume_enabled").(bool)
			isBreakthroughMode := d.Get("breakthrough_mode_enabled").(bool)
			storageQuotaInGB := d.Get("storage_quota_in_gb").(int)
			coolAccessConfig := d.Get("cool_access").([]any)

			// Breakthrough Mode places the volume on dedicated capacity and is only supported on large volumes
			if isBreakthroughMode && !isLargeVolume {
				return fmt.Errorf("`breakthrough_mode_enabled` can only be used on large volumes; set `large_volume_enabled` to true")
			}

			// Breakthrough Mode and cool access cannot both be requested when the volume is created, cool access can
			// however be enabled on a subsequent update of a Breakthrough Mode volume
			if isBreakthroughMode && len(coolAccessConfig) > 0 && d.Id() == "" {
				return fmt.Errorf("`cool_access` cannot be configured when `breakthrough_mode_enabled` is true at creation time; enable `cool_access` in a subsequent update instead")
			}

			switch {
			case isBreakthroughMode && storageQuotaInGB < 2400:
				// Breakthrough Mode volumes must be at least 2,400 GiB
				return fmt.Errorf("when `breakthrough_mode_enabled` is true, `storage_quota_in_gb` must be at least 2,400 GB (2,400 GiB)")
			case isBreakthroughMode && storageQuotaInGB > 2457600:
				// Breakthrough Mode volumes can grow up to 2,400 TiB (2,457,600 GB)
				return fmt.Errorf("when `breakthrough_mode_enabled` is true, `storage_quota_in_gb` must not exceed 2,457,600 GB (2,400 TiB)")
			case !isBreakthroughMode && isLargeVolume && storageQuotaInGB < 51200:
				// Large volumes must be at least 50 TiB (51,200 GB)
				return fmt.Errorf("when `large_volume_enabled` is true, `storage_quota_in_gb` must be at least 51,200 GB (50 TiB)")
			case !isBreakthroughMode && isLargeVolume && storageQuotaInGB > 1048576:
				// Validate against the maximum (1 PiB / 1,048,576 GB)
				return fmt.Errorf("`storage_quota_in_gb` must not exceed 1,048,576 GB (1 PiB); larger sizes require requesting special quota or setting `breakthrough_mode_enabled` to true")
			case !isLargeVolume && storageQuotaInGB > 102400:
				// Non-large volumes cannot be larger than 100 TiB (102,400 GB)
				return fmt.Errorf("when `large_volume_enabled` is false, `storage_quota_in_gb` must not exceed 102,400 GB (100 TiB); set `large_volume_enabled` to true for larger volumes")
			default:
				// All validations passed - no action needed
			}

			// Validate short-term clone properties
			acceptGrowCapacityPool := d.Get("accept_grow_capacity_pool_for_short_term_clone_split").(string)

			// Only validate if accept_grow_capacity_pool_for_short_term_clone_split is set
			if acceptGrowCapacityPool != "" {
				rawConfig := d.GetRawConfig().AsValueMap()
				// Check if create_from_snapshot_resource_id is explicitly null or empty
				// Allow computed values (where value is empty but raw config is not null)
				if rawConfig["create_from_snapshot_resource_id"].IsNull() {
					return fmt.Errorf("`accept_grow_capacity_pool_for_short_term_clone_split` can only be used in conjunction with `create_from_snapshot_resource_id`")
				}
			}

			// Validate that short-term clones are not supported on large volumes
			if acceptGrowCapacityPool != "" && isLargeVolume {
				return fmt.Errorf("short-term clones are not supported on large volumes; `accept_grow_capacity_pool_for_short_term_clone_split` cannot be used when `large_volume_enabled` is true")
			}

			// Validate that short-term clones are not supported with cool access
			if acceptGrowCapacityPool != "" && len(coolAccessConfig) > 0 {
				return fmt.Errorf("short-term clones are not supported on volumes enabled for cool access; `accept_grow_capacity_pool_for_short_term_clone_split` cannot be used when `cool_access` is configured")
			}

			if d.HasChanges("service_level", "pool_name") {
				serviceLevelChange := d.HasChange("service_level")
				poolNameChange := d.HasChange("pool_name")

				// `service_level` and `pool_name` must be updated together or we ForceNew the resource
				// https://learn.microsoft.com/azure/azure-netapp-files/dynamic-change-volume-service-level
				if serviceLevelChange && !poolNameChange {
					return d.ForceNew("service_level")
				}
				if !serviceLevelChange && poolNameChange {
					return d.ForceNew("pool_name")
				}
			}

			// Validate NFSv3 to NFSv4.1 protocol conversion restrictions
			if d.HasChange("protocols") {
				old, new := d.GetChange("protocols")
				oldProtocols := old.(*pluginsdk.Set).List()
				newProtocols := new.(*pluginsdk.Set).List()

				// Convert to string slices for validation
				oldProtocolsStr := make([]string, len(oldProtocols))
				newProtocolsStr := make([]string, len(newProtocols))

				for i, v := range oldProtocols {
					oldProtocolsStr[i] = v.(string)
				}
				for i, v := range newProtocols {
					newProtocolsStr[i] = v.(string)
				}

				kerberosEnabled := d.Get("kerberos_enabled").(bool)
				dataReplication := d.Get("data_protection_replication").([]any)

				// Get the new export policy rules configuration to validate against new protocols
				// Always use the new configuration when protocols are changing to ensure validation against intended state
				exportPolicyRules := d.Get("export_policy_rule").([]any)

				validationErrors := validate.ValidateNetAppVolumeProtocolConversion(oldProtocolsStr, newProtocolsStr, kerberosEnabled, dataReplication, exportPolicyRules)
				for _, err := range validationErrors {
					return err
				}
			}

			// Validate cross-zone-region replication requirements
			// According to Azure documentation, for cross-zone replication, both source and destination volumes must have zones
			dataReplicationRaw := d.Get("data_protection_replication").([]any)
			if len(dataReplicationRaw) > 0 && dataReplicationRaw[0] != nil {
				// This is a destination volume with data_protection_replication configured
				dataReplication := dataReplicationRaw[0].(map[string]any)
				remoteVolumeLocation := dataReplication["remote_volume_location"].(string)
				currentLocation := d.Get("location").(string)

				// Check if this is cross-zone replication (same region)
				if strings.EqualFold(location.Normalize(remoteVolumeLocation), location.Normalize(currentLocation)) {
					// Cross-zone replication: both source and destination must have zones assigned
					zone := d.Get("zone").(string)
					if zone == "" {
						return fmt.Errorf("when configuring cross-zone replication (data_protection_replication with same region), the destination volume must have a `zone` assigned")
					}
					// Note: We cannot validate the source volume's zone here since we don't have access to it during plan/diff
					// The documentation states the source must also have a zone, which users must ensure separately
				}
			}

			return nil
		},
	}
}
