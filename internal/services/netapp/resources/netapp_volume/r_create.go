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
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppVolumeCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := volumes.NewVolumeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("pool_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_volume", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))

	zones := &[]string{}
	if v, ok := d.GetOk("zone"); ok {
		zones = &[]string{
			v.(string),
		}
	}

	volumePath := d.Get("volume_path").(string)
	serviceLevel := volumes.ServiceLevel(d.Get("service_level").(string))
	subnetID := d.Get("subnet_id").(string)
	kerberosEnabled := d.Get("kerberos_enabled").(bool)
	smbContinuouslyAvailable := d.Get("smb_continuous_availability_enabled").(bool)
	smbEncryption := d.Get("smb3_protocol_encryption_enabled").(bool)
	networkFeatures := volumes.NetworkFeatures(d.Get("network_features").(string))

	smbNonBrowsable := volumes.SmbNonBrowsableDisabled
	if d.Get("smb_non_browsable_enabled").(bool) {
		smbNonBrowsable = volumes.SmbNonBrowsableEnabled
	}

	smbAccessBasedEnumeration := volumes.SmbAccessBasedEnumerationDisabled
	if d.Get("smb_access_based_enumeration_enabled").(bool) {
		smbAccessBasedEnumeration = volumes.SmbAccessBasedEnumerationEnabled
	}

	protocols := d.Get("protocols").(*pluginsdk.Set).List()
	if len(protocols) == 0 {
		protocols = append(protocols, "NFSv3")
	}

	// Handling security style property
	securityStyle := volumes.SecurityStyle(d.Get("security_style").(string))
	if strings.EqualFold(string(securityStyle), "unix") && len(protocols) == 1 && strings.EqualFold(protocols[0].(string), "cifs") {
		return fmt.Errorf("unix security style cannot be used in a CIFS enabled volume for %s", id)
	}
	if strings.EqualFold(string(securityStyle), "ntfs") && len(protocols) == 1 && (strings.EqualFold(protocols[0].(string), "nfsv3") || strings.EqualFold(protocols[0].(string), "nfsv4.1")) {
		return fmt.Errorf("ntfs security style cannot be used in a NFSv3/NFSv4.1 enabled volume for %s", id)
	}

	acceptGrowCapacityPool := d.Get("accept_grow_capacity_pool_for_short_term_clone_split").(string)

	storageQuotaInGB := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)

	exportPolicyRuleRaw := d.Get("export_policy_rule").([]any)
	exportPolicyRule := expandNetAppVolumeExportPolicyRule(exportPolicyRuleRaw)

	dataProtectionReplicationRaw := d.Get("data_protection_replication").([]any)
	dataProtectionReplication := helpers.ExpandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

	dataProtectionSnapshotPolicyRaw := d.Get("data_protection_snapshot_policy").([]any)
	dataProtectionSnapshotPolicy := expandNetAppVolumeDataProtectionSnapshotPolicy(dataProtectionSnapshotPolicyRaw)

	dataProtectionBackupPolicyRaw := d.Get("data_protection_backup_policy").([]any)
	dataProtectionBackupPolicy := expandNetAppVolumeDataProtectionBackupPolicy(dataProtectionBackupPolicyRaw)

	dataProtectionARPRaw := d.Get("data_protection_advanced_ransomware").([]any)
	dataProtectionARP := expandNetAppVolumeDataProtectionAdvancedRansomwareProtection(dataProtectionARPRaw)

	authorizeReplication := false
	volumeType := ""
	endpointType := ""
	if dataProtectionReplication != nil {
		if dataProtectionReplication.EndpointType != nil {
			endpointType = string(*dataProtectionReplication.EndpointType)
		}
		if strings.EqualFold(endpointType, "dst") {
			authorizeReplication = true
			volumeType = volumeTypeDataProtection
		}
	}

	// Validate applicability of backup policies
	if dataProtectionBackupPolicy != nil {
		// Validate that backup policies are not being enforced in a data protection replication destination volume
		if strings.EqualFold(volumeType, "dst") && dataProtectionBackupPolicy.PolicyEnforced == pointer.To(true) {
			return fmt.Errorf("backup policy cannot be enforced on a data protection destination volume, NetApp Volume %q (Resource Group %q)", id.VolumeName, id.ResourceGroupName)
		}
	}

	// Validating that snapshot policies are not being created in a data protection replication volume
	if dataProtectionSnapshotPolicy != nil && volumeType != "" {
		return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, NetApp Volume %q (Resource Group %q)", id.VolumeName, id.ResourceGroupName)
	}

	snapshotDirectoryVisible := d.Get("snapshot_directory_visible").(bool)

	// Handling volume creation from snapshot case
	snapshotResourceID := d.Get("create_from_snapshot_resource_id").(string)
	snapshotID := ""
	if snapshotResourceID != "" {
		// Get snapshot ID GUID value
		parsedSnapshotResourceID, err := snapshots.ParseSnapshotID(snapshotResourceID)
		if err != nil {
			return fmt.Errorf("parsing snapshotResourceID %q: %+v", snapshotResourceID, err)
		}

		snapshotClient := meta.(*clients.Client).NetApp.SnapshotClient
		snapshotResponse, err := snapshotClient.Get(ctx, *parsedSnapshotResourceID)
		if err != nil {
			return fmt.Errorf("getting snapshot from %s: %+v", id, err)
		}
		if model := snapshotResponse.Model; model != nil && model.Id != nil {
			snapshotID = *model.Id
		}

		sourceVolumeId := volumes.NewVolumeID(parsedSnapshotResourceID.SubscriptionId, parsedSnapshotResourceID.ResourceGroupName, parsedSnapshotResourceID.NetAppAccountName, parsedSnapshotResourceID.CapacityPoolName, parsedSnapshotResourceID.VolumeName)
		// Validate if properties that cannot be changed matches (protocols, subnet_id, location, resource group, account_name, pool_name, service_level)
		sourceVolume, err := client.Get(ctx, sourceVolumeId)
		if err != nil {
			return fmt.Errorf("getting source NetApp Volume (snapshot's parent resource) %q (Resource Group %q): %+v", parsedSnapshotResourceID.VolumeName, parsedSnapshotResourceID.ResourceGroupName, err)
		}

		propertyMismatch := []string{}
		if model := sourceVolume.Model; model != nil {
			props := model.Properties
			if !ValidateSlicesEquality(*props.ProtocolTypes, *pluginsdk.ExpandStringSlice(protocols), false) {
				propertyMismatch = append(propertyMismatch, "protocols")
			}
			if !strings.EqualFold(props.SubnetId, subnetID) {
				propertyMismatch = append(propertyMismatch, "subnet_id")
			}
			if !strings.EqualFold(model.Location, location) {
				propertyMismatch = append(propertyMismatch, "location")
			}
			if volumeServiceLevel := props.ServiceLevel; volumeServiceLevel != nil {
				if !strings.EqualFold(string(*props.ServiceLevel), string(serviceLevel)) {
					propertyMismatch = append(propertyMismatch, "service_level")
				}
			}
			if !strings.EqualFold(sourceVolumeId.ResourceGroupName, id.ResourceGroupName) {
				propertyMismatch = append(propertyMismatch, "resource_group_name")
			}
			if !strings.EqualFold(sourceVolumeId.NetAppAccountName, id.NetAppAccountName) {
				propertyMismatch = append(propertyMismatch, "account_name")
			}
			if len(propertyMismatch) > 0 {
				return fmt.Errorf("the following properties to create a new NetApp Volume from a Snapshot do not match:\n%s", strings.Join(propertyMismatch, "\n"))
			}
		}

		if acceptGrowCapacityPool != "" {
			volumeType = volumeTypeShortTermClone
		}
	}

	avsDataStoreEnabled := volumes.AvsDataStoreDisabled
	if d.Get("azure_vmware_data_store_enabled").(bool) {
		avsDataStoreEnabled = volumes.AvsDataStoreEnabled
	}

	parameters := volumes.Volume{
		Location: location,
		Properties: volumes.VolumeProperties{
			CreationToken:             volumePath,
			ServiceLevel:              &serviceLevel,
			SubnetId:                  subnetID,
			KerberosEnabled:           &kerberosEnabled,
			SmbContinuouslyAvailable:  &smbContinuouslyAvailable,
			SmbEncryption:             &smbEncryption,
			NetworkFeatures:           &networkFeatures,
			SmbNonBrowsable:           &smbNonBrowsable,
			SmbAccessBasedEnumeration: &smbAccessBasedEnumeration,
			ProtocolTypes:             pluginsdk.ExpandStringSlice(protocols),
			SecurityStyle:             &securityStyle,
			UsageThreshold:            storageQuotaInGB,
			ExportPolicy:              exportPolicyRule,
			VolumeType:                pointer.To(volumeType),
			SnapshotId:                pointer.To(snapshotID),
			DataProtection: &volumes.VolumePropertiesDataProtection{
				Replication:          dataProtectionReplication,
				Snapshot:             dataProtectionSnapshotPolicy,
				Backup:               dataProtectionBackupPolicy,
				RansomwareProtection: dataProtectionARP,
			},
			AvsDataStore:             &avsDataStoreEnabled,
			SnapshotDirectoryVisible: pointer.To(snapshotDirectoryVisible),
			IsLargeVolume:            pointer.To(d.Get("large_volume_enabled").(bool)),
			UnixPermissions:          pointer.To(d.Get("unix_permissions").(string)),
		},
		Tags:  tags.Expand(d.Get("tags").(map[string]any)),
		Zones: zones,
	}

	if acceptGrowCapacityPool != "" {
		parameters.Properties.AcceptGrowCapacityPoolForShortTermCloneSplit = pointer.ToEnum[volumes.AcceptGrowCapacityPoolForShortTermCloneSplit](acceptGrowCapacityPool)
	}

	// Breakthrough Mode is only sent when requested, omitting it keeps the volume on the standard configuration
	// and avoids sending an unsupported property to subscriptions that are not onboarded to the feature
	if d.Get("breakthrough_mode_enabled").(bool) {
		parameters.Properties.BreakthroughMode = pointer.To(volumes.BreakthroughModeEnabled)
	}

	if throughputMibps, ok := d.GetOk("throughput_in_mibps"); ok {
		parameters.Properties.ThroughputMibps = pointer.To(throughputMibps.(float64))
	}

	if len(d.Get("cool_access").([]any)) > 0 {
		coolAccess := d.Get("cool_access").([]any)[0].(map[string]any)
		parameters.Properties.CoolAccess = pointer.To(true)
		parameters.Properties.CoolAccessRetrievalPolicy = pointer.ToEnum[volumes.CoolAccessRetrievalPolicy](coolAccess["retrieval_policy"].(string))
		parameters.Properties.CoolAccessTieringPolicy = pointer.ToEnum[volumes.CoolAccessTieringPolicy](coolAccess["tiering_policy"].(string))
		parameters.Properties.CoolnessPeriod = pointer.To(int64(coolAccess["coolness_period_in_days"].(int)))
	}

	if encryptionKeySource, ok := d.GetOk("encryption_key_source"); ok {
		// Validating Microsoft.KeyVault encryption key provider is enabled only on Standard network features
		if volumes.EncryptionKeySource(encryptionKeySource.(string)) == volumes.EncryptionKeySourceMicrosoftPointKeyVault && networkFeatures == volumes.NetworkFeaturesBasic {
			return fmt.Errorf("volume encryption cannot be enabled when network features is set to basic: %s", id.ID())
		}

		parameters.Properties.EncryptionKeySource = pointer.ToEnum[volumes.EncryptionKeySource](encryptionKeySource.(string))
	}

	if keyVaultPrivateEndpointID, ok := d.GetOk("key_vault_private_endpoint_id"); ok {
		parameters.Properties.KeyVaultPrivateEndpointResourceId = pointer.To(keyVaultPrivateEndpointID.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// Waiting for volume be completely provisioned
	if err := helpers.WaitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
		return err
	}

	// If this is a data replication secondary volume, authorize replication on primary volume
	if authorizeReplication {
		replVolID, err := volumes.ParseVolumeID(pointer.From(dataProtectionReplication.RemoteVolumeResourceId))
		if err != nil {
			return err
		}

		if err = client.AuthorizeReplicationThenPoll(
			ctx, *replVolID, volumes.AuthorizeRequest{
				RemoteVolumeResourceId: pointer.To(id.ID()),
			},
		); err != nil {
			return fmt.Errorf("cannot authorize volume replication: %v", err)
		}

		// Wait for volume replication authorization to complete on the destination volume
		// Note: We check the destination (current volume being created), not the source
		// This is important for one-to-many replication where checking the source would fail
		destinationReplID, err := volumes.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Waiting for replication authorization on destination volume %s to complete", id)
		if err := helpers.WaitForReplAuthorization(ctx, client, *destinationReplID); err != nil {
			return err
		}
	}

	return resourceNetAppVolumeRead(d, meta)
}

func expandNetAppVolumeExportPolicyRule(input []any) *volumes.VolumePropertiesExportPolicy {
	results := make([]volumes.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]any)
			ruleIndex := int64(v["rule_index"].(int))
			allowedClients := strings.Join(*pluginsdk.ExpandStringSlice(v["allowed_clients"].(*pluginsdk.Set).List()), ",")

			cifsEnabled := false
			nfsv3Enabled := false
			nfsv41Enabled := false
			if vpe := v["protocol"]; vpe != nil {
				protocolsEnabled := vpe.([]any)
				if len(protocolsEnabled) != 0 {
					for _, protocol := range protocolsEnabled {
						if protocol != nil {
							switch strings.ToLower(protocol.(string)) {
							case "cifs":
								cifsEnabled = true
							case "nfsv3":
								nfsv3Enabled = true
							case "nfsv4.1":
								nfsv41Enabled = true
							}
						}
					}
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)
			kerberos5ro := v["kerberos_5_read_only_enabled"].(bool)
			kerberos5rw := v["kerberos_5_read_write_enabled"].(bool)
			kerberos5iro := v["kerberos_5i_read_only_enabled"].(bool)
			kerberos5irw := v["kerberos_5i_read_write_enabled"].(bool)
			kerberos5pro := v["kerberos_5p_read_only_enabled"].(bool)
			kerberos5prw := v["kerberos_5p_read_write_enabled"].(bool)

			result := volumes.ExportPolicyRule{
				AllowedClients:      pointer.To(allowedClients),
				Cifs:                pointer.To(cifsEnabled),
				Nfsv3:               pointer.To(nfsv3Enabled),
				Nfsv41:              pointer.To(nfsv41Enabled),
				Kerberos5ReadOnly:   pointer.To(kerberos5ro),
				Kerberos5ReadWrite:  pointer.To(kerberos5rw),
				Kerberos5iReadOnly:  pointer.To(kerberos5iro),
				Kerberos5iReadWrite: pointer.To(kerberos5irw),
				Kerberos5pReadOnly:  pointer.To(kerberos5pro),
				Kerberos5pReadWrite: pointer.To(kerberos5prw),
				RuleIndex:           pointer.To(ruleIndex),
				UnixReadOnly:        pointer.To(unixReadOnly),
				UnixReadWrite:       pointer.To(unixReadWrite),
				HasRootAccess:       pointer.To(rootAccessEnabled),
			}

			results = append(results, result)
		}
	}

	return &volumes.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeExportPolicyRulePatch(input []any, overrideProtocols []string) *volumes.VolumePatchPropertiesExportPolicy {
	results := make([]volumes.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]any)
			ruleIndex := int64(v["rule_index"].(int))
			allowedClients := strings.Join(*pluginsdk.ExpandStringSlice(v["allowed_clients"].(*pluginsdk.Set).List()), ",")

			nfsv3Enabled := false
			nfsv41Enabled := false
			cifsEnabled := false

			// If overrideProtocols is provided (during protocol conversion), use those protocols
			// instead of reading from the export policy rule configuration
			if len(overrideProtocols) > 0 {
				for _, protocol := range overrideProtocols {
					switch strings.ToLower(protocol) {
					case "cifs":
						cifsEnabled = true
					case "nfsv3":
						nfsv3Enabled = true
					case "nfsv4.1":
						nfsv41Enabled = true
					}
				}
			} else {
				// Use existing logic when no protocol override is provided
				if vpe := v["protocol"]; vpe != nil {
					protocolsEnabled := vpe.([]any)
					if len(protocolsEnabled) != 0 {
						for _, protocol := range protocolsEnabled {
							if protocol != nil {
								switch strings.ToLower(protocol.(string)) {
								case "cifs":
									cifsEnabled = true
								case "nfsv3":
									nfsv3Enabled = true
								case "nfsv4.1":
									nfsv41Enabled = true
								}
							}
						}
					}
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)

			result := volumes.ExportPolicyRule{
				AllowedClients: pointer.To(allowedClients),
				Cifs:           pointer.To(cifsEnabled),
				Nfsv3:          pointer.To(nfsv3Enabled),
				Nfsv41:         pointer.To(nfsv41Enabled),
				RuleIndex:      pointer.To(ruleIndex),
				UnixReadOnly:   pointer.To(unixReadOnly),
				UnixReadWrite:  pointer.To(unixReadWrite),
				HasRootAccess:  pointer.To(rootAccessEnabled),
			}

			results = append(results, result)
		}
	}

	return &volumes.VolumePatchPropertiesExportPolicy{
		Rules: &results,
	}
}
