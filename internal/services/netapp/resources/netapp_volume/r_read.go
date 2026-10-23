// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppVolumeRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.VolumeName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("pool_name", id.CapacityPoolName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		zone := ""
		if model.Zones != nil {
			if zones := *model.Zones; len(zones) > 0 {
				zone = zones[0]
			}
		}
		d.Set("zone", zone)

		props := model.Properties
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", pointer.FromEnum(props.ServiceLevel))
		d.Set("subnet_id", props.SubnetId)
		d.Set("kerberos_enabled", props.KerberosEnabled)
		d.Set("smb_continuous_availability_enabled", props.SmbContinuouslyAvailable)
		d.Set("smb3_protocol_encryption_enabled", props.SmbEncryption)
		d.Set("network_features", pointer.FromEnum(props.NetworkFeatures))
		d.Set("protocols", props.ProtocolTypes)
		d.Set("security_style", pointer.FromEnum(props.SecurityStyle))
		d.Set("snapshot_directory_visible", props.SnapshotDirectoryVisible)
		d.Set("throughput_in_mibps", props.ThroughputMibps)
		d.Set("storage_quota_in_gb", props.UsageThreshold/1073741824)
		d.Set("encryption_key_source", pointer.FromEnum(props.EncryptionKeySource))
		d.Set("key_vault_private_endpoint_id", props.KeyVaultPrivateEndpointResourceId)
		d.Set("large_volume_enabled", props.IsLargeVolume)
		d.Set("breakthrough_mode_enabled", pointer.From(props.BreakthroughMode) == volumes.BreakthroughModeEnabled)
		d.Set("accept_grow_capacity_pool_for_short_term_clone_split", pointer.FromEnum(props.AcceptGrowCapacityPoolForShortTermCloneSplit))
		d.Set("unix_permissions", pointer.From(props.UnixPermissions))

		if pointer.From(props.CoolAccess) {
			// enums returned from the API are inconsistent so normalize them here
			// https://github.com/Azure/azure-rest-api-specs/issues/35371
			coolAccess := map[string]any{
				"retrieval_policy":        normalizeCoolAccessRetrievalPolicy(pointer.From(props.CoolAccessRetrievalPolicy)),
				"tiering_policy":          normalizeCoolAccessTieringPolicy(pointer.From(props.CoolAccessTieringPolicy)),
				"coolness_period_in_days": pointer.From(props.CoolnessPeriod),
			}
			d.Set("cool_access", []any{coolAccess})
		} else {
			d.Set("cool_access", []any{})
		}

		smbNonBrowsable := false
		if props.SmbNonBrowsable != nil {
			smbNonBrowsable = strings.EqualFold(string(*props.SmbNonBrowsable), string(volumes.SmbNonBrowsableEnabled))
		}
		d.Set("smb_non_browsable_enabled", smbNonBrowsable)

		smbAccessBasedEnumeration := false
		if props.SmbAccessBasedEnumeration != nil {
			smbAccessBasedEnumeration = strings.EqualFold(string(*props.SmbAccessBasedEnumeration), string(volumes.SmbAccessBasedEnumerationEnabled))
		}
		d.Set("smb_access_based_enumeration_enabled", smbAccessBasedEnumeration)

		avsDataStore := false
		if props.AvsDataStore != nil {
			avsDataStore = strings.EqualFold(string(*props.AvsDataStore), string(volumes.AvsDataStoreEnabled))
		}
		d.Set("azure_vmware_data_store_enabled", avsDataStore)

		if err := d.Set("export_policy_rule", flattenNetAppVolumeExportPolicyRule(props.ExportPolicy)); err != nil {
			return fmt.Errorf("setting `export_policy_rule`: %+v", err)
		}
		if err := d.Set("mount_target", flattenNetAppVolumeMountTargets(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_target`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
		if err := d.Set("data_protection_snapshot_policy", flattenNetAppVolumeDataProtectionSnapshotPolicy(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_snapshot_policy`: %+v", err)
		}
		if err := d.Set("data_protection_backup_policy", flattenNetAppVolumeDataProtectionBackupPolicy(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_backup_policy`: %+v", err)
		}
		if err := d.Set("data_protection_advanced_ransomware", flattenNetAppVolumeDataProtectionAdvancedRansomwareProtection(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_advanced_ransomware`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

const (
	// Volume types used in Azure NetApp Files
	volumeTypeDataProtection = "DataProtection"
	volumeTypeShortTermClone = "ShortTermClone"
)

func normalizeCoolAccessTieringPolicy(coolAccessTieringPolicy volumes.CoolAccessTieringPolicy) string {
	switch strings.ToLower(string(coolAccessTieringPolicy)) {
	case "auto":
		return string(volumes.CoolAccessTieringPolicyAuto)
	case "snapshot-only":
		return string(volumes.CoolAccessTieringPolicySnapshotOnly)
	default:
		return string(coolAccessTieringPolicy)
	}
}

func normalizeCoolAccessRetrievalPolicy(coolAccessRetrievalPolicy volumes.CoolAccessRetrievalPolicy) string {
	switch strings.ToLower(string(coolAccessRetrievalPolicy)) {
	case "never":
		return string(volumes.CoolAccessRetrievalPolicyNever)
	case "default":
		return string(volumes.CoolAccessRetrievalPolicyDefault)
	case "on-read":
		return string(volumes.CoolAccessRetrievalPolicyOnRead)
	default:
		return string(coolAccessRetrievalPolicy)
	}
}
