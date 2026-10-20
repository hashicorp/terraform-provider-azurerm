// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_oracle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			volumeClient := metadata.Client.NetApp.VolumeClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppVolumeGroupOracleModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("volume") {
				for i := 0; i < metadata.ResourceData.Get("volume.#").(int); i++ {
					// Checking if individual volume has a change
					volumeItem := fmt.Sprintf("volume.%v", i)

					capacityPoolId, err := capacitypools.ParseCapacityPoolID(metadata.ResourceData.Get(fmt.Sprintf("%v.capacity_pool_id", volumeItem)).(string))
					if err != nil {
						return err
					}

					if metadata.ResourceData.HasChange(volumeItem) {
						volumeId := volumes.NewVolumeID(id.SubscriptionId,
							id.ResourceGroupName,
							id.NetAppAccountName,
							capacityPoolId.CapacityPoolName,
							metadata.ResourceData.Get(fmt.Sprintf("%v.name", volumeItem)).(string))

						update := volumes.VolumePatch{
							Properties: &volumes.VolumePatchProperties{},
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)) {
							storageQuotaInBytes := int64(metadata.ResourceData.Get(fmt.Sprintf("%v.storage_quota_in_gb", volumeItem)).(int) * 1073741824)
							update.Properties.UsageThreshold = pointer.To(storageQuotaInBytes)
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.export_policy_rule", volumeItem)) {
							exportPolicyRuleRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.export_policy_rule", volumeItem)).([]any)

							// Validating export policy rules
							volumeProtocolRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.protocols", volumeItem)).([]any)[0]
							volumeProtocol := volumeProtocolRaw.(string)

							errors := make([]error, 0)
							for _, ruleRaw := range exportPolicyRuleRaw {
								if ruleRaw != nil {
									rule := volumegroups.ExportPolicyRule{}

									v := ruleRaw.(map[string]any)
									rule.Nfsv3 = pointer.To(v["nfsv3_enabled"].(bool))
									rule.Nfsv41 = pointer.To(v["nfsv41_enabled"].(bool))

									errors = append(errors, validate.ValidateNetAppVolumeGroupExportPolicyRule(rule, volumeProtocol)...)
								}
							}

							if len(errors) > 0 {
								return fmt.Errorf("one or more issues found while performing export policies validations for %s:\n%+v", id, errors)
							}

							var protocolOverride []string
							// Only override export policy protocols if we're also changing volume protocols
							if metadata.ResourceData.HasChange(fmt.Sprintf("%v.protocols", volumeItem)) {
								protocolsRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.protocols", volumeItem)).([]any)
								protocolOverride = make([]string, len(protocolsRaw))
								for i, p := range protocolsRaw {
									protocolOverride[i] = p.(string)
								}
							}

							update.Properties.ExportPolicy = helpers.ExpandNetAppVolumeGroupVolumeExportPolicyRulePatchWithProtocolConversion(exportPolicyRuleRaw, protocolOverride)
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.protocols", volumeItem)) {
							protocolsRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.protocols", volumeItem)).([]any)
							protocols := make([]string, len(protocolsRaw))
							for i, p := range protocolsRaw {
								protocols[i] = p.(string)
							}
							update.Properties.ProtocolTypes = pointer.To(protocols)
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.data_protection_snapshot_policy", volumeItem)) {
							// Validating that snapshot policies are not being created in a data protection volume
							dataProtectionReplicationRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.data_protection_replication", volumeItem)).([]any)
							dataProtectionReplication := helpers.ExpandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

							if dataProtectionReplication != nil &&
								dataProtectionReplication.EndpointType != nil &&
								strings.EqualFold(pointer.FromEnum(dataProtectionReplication.EndpointType), string(volumegroups.EndpointTypeDst)) {
								return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, %s", volumeId)
							}

							dataProtectionSnapshotPolicyRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.data_protection_snapshot_policy", volumeItem)).([]any)
							dataProtectionSnapshotPolicy := helpers.ExpandNetAppVolumeDataProtectionSnapshotPolicyPatch(dataProtectionSnapshotPolicyRaw)
							update.Properties.DataProtection = &volumes.VolumePatchPropertiesDataProtection{
								Snapshot: dataProtectionSnapshotPolicy,
							}
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.throughput_in_mibps", volumeItem)) {
							throughputMibps := metadata.ResourceData.Get(fmt.Sprintf("%v.throughput_in_mibps", volumeItem))
							update.Properties.ThroughputMibps = pointer.To(throughputMibps.(float64))
						}

						if metadata.ResourceData.HasChange(fmt.Sprintf("%v.tags", volumeItem)) {
							tagsRaw := metadata.ResourceData.Get(fmt.Sprintf("%v.tags", volumeItem)).(map[string]any)
							update.Tags = tags.Expand(tagsRaw)
						}

						if err = volumeClient.UpdateThenPoll(ctx, volumeId, update); err != nil {
							return fmt.Errorf("updating %s: %+v", volumeId, err)
						}

						// Waiting for volume to fully complete an update
						if err := helpers.WaitForVolumeCreateOrUpdate(ctx, volumeClient, volumeId); err != nil {
							return fmt.Errorf("waiting update %s: %+v", volumeId, err)
						}
					}
				}
			}

			return nil
		},
	}
}
