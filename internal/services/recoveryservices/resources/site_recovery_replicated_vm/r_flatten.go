// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"fmt"

	//nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenSiteRecoveryReplicatedItem(d *pluginsdk.ResourceData, model *replicationprotecteditems.ReplicationProtectedItem) error {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", id.ReplicationProtectedItemName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("source_recovery_fabric_name", id.ReplicationFabricName)
	d.Set("source_recovery_protection_container_name", id.ReplicationProtectionContainerName)

	if prop := model.Properties; prop != nil {
		recoveryFabricId := ""
		if fabricId := pointer.From(prop.RecoveryFabricId); fabricId != "" {
			parsedFabricId, err := replicationfabrics.ParseReplicationFabricIDInsensitively(fabricId)
			if err != nil {
				return err
			}
			recoveryFabricId = parsedFabricId.ID()
		}
		d.Set("target_recovery_fabric_id", recoveryFabricId)

		recoveryPolicyId := ""
		if policyId := pointer.From(prop.PolicyId); policyId != "" {
			parsedPolicyId, err := replicationpolicies.ParseReplicationPolicyIDInsensitively(policyId)
			if err != nil {
				return err
			}
			recoveryPolicyId = parsedPolicyId.ID()
		}
		d.Set("recovery_replication_policy_id", recoveryPolicyId)

		recoveryContainerId := ""
		if containerId := pointer.From(prop.RecoveryContainerId); containerId != "" {
			parsedContainerId, err := replicationprotecteditems.ParseReplicationProtectionContainerIDInsensitively(containerId)
			if err != nil {
				return err
			}
			recoveryContainerId = parsedContainerId.ID()
		}
		d.Set("target_recovery_protection_container_id", recoveryContainerId)

		if a2aDetails, isA2a := prop.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
			sourceVmId := ""
			if objId := pointer.From(a2aDetails.FabricObjectId); objId != "" {
				parsedVmID, err := virtualmachines.ParseVirtualMachineIDInsensitively(objId)
				if err != nil {
					return err
				}
				sourceVmId = parsedVmID.ID()
			}
			d.Set("source_vm_id", sourceVmId)

			recoveryGroupId := ""
			if groupId := pointer.From(a2aDetails.RecoveryAzureResourceGroupId); groupId != "" {
				parsedGroupId, err := commonids.ParseResourceGroupIDInsensitively(groupId)
				if err != nil {
					return err
				}
				recoveryGroupId = parsedGroupId.ID()
			}
			d.Set("target_resource_group_id", recoveryGroupId)

			availabilitySetId := ""
			if respAvailabilitySetId := pointer.From(a2aDetails.RecoveryAvailabilitySet); respAvailabilitySetId != "" {
				parsedAvailabilitySetId, err := commonids.ParseAvailabilitySetIDInsensitively(respAvailabilitySetId)
				if err != nil {
					return err
				}
				availabilitySetId = parsedAvailabilitySetId.ID()
			}
			d.Set("target_availability_set_id", availabilitySetId)

			targetNetworkId := ""
			if respTargetNetworkId := pointer.From(a2aDetails.SelectedRecoveryAzureNetworkId); respTargetNetworkId != "" {
				parsedTargetNetworkId, err := commonids.ParseVirtualNetworkIDInsensitively(respTargetNetworkId)
				if err != nil {
					return err
				}
				targetNetworkId = parsedTargetNetworkId.ID()
			}
			d.Set("target_network_id", targetNetworkId)

			testNetworkId := ""
			if respTfoNetworkId := pointer.From(a2aDetails.SelectedTfoAzureNetworkId); respTfoNetworkId != "" {
				parsedTfoNetworkId, err := commonids.ParseVirtualNetworkIDInsensitively(respTfoNetworkId)
				if err != nil {
					return err
				}
				testNetworkId = parsedTfoNetworkId.ID()
			}
			d.Set("test_network_id", testNetworkId)

			proximityPlacementGroupId := ""
			if respProximityPlacementGroupId := pointer.From(a2aDetails.RecoveryProximityPlacementGroupId); respProximityPlacementGroupId != "" {
				parsedProximityPlacementGroupId, err := proximityplacementgroups.ParseProximityPlacementGroupIDInsensitively(respProximityPlacementGroupId)
				if err != nil {
					return err
				}
				proximityPlacementGroupId = parsedProximityPlacementGroupId.ID()
			}
			d.Set("target_proximity_placement_group_id", proximityPlacementGroupId)

			recoveryBootDiagStorageAccount := ""
			if respBootDiagStorageAccountId := pointer.From(a2aDetails.RecoveryBootDiagStorageAccountId); respBootDiagStorageAccountId != "" {
				parsedRecoveryBootDiagStorageAccount, err := commonids.ParseStorageAccountIDInsensitively(respBootDiagStorageAccountId)
				if err != nil {
					return err
				}
				recoveryBootDiagStorageAccount = parsedRecoveryBootDiagStorageAccount.ID()
			}
			d.Set("target_boot_diagnostic_storage_account_id", recoveryBootDiagStorageAccount)

			capReservationGroupId := ""
			if respCapacityGroupId := pointer.From(a2aDetails.RecoveryCapacityReservationGroupId); respCapacityGroupId != "" {
				parsedCapReservaGroupId, err := capacityreservationgroups.ParseCapacityReservationGroupIDInsensitively(respCapacityGroupId)
				if err != nil {
					return err
				}
				capReservationGroupId = parsedCapReservaGroupId.ID()
			}
			d.Set("target_capacity_reservation_group_id", capReservationGroupId)

			vmssId := ""
			if respVmssId := pointer.From(a2aDetails.RecoveryVirtualMachineScaleSetId); respVmssId != "" {
				parsedVmssId, err := commonids.ParseVirtualMachineScaleSetIDInsensitively(respVmssId)
				if err != nil {
					return err
				}
				vmssId = parsedVmssId.ID()
			}
			d.Set("target_virtual_machine_scale_set_id", vmssId)

			d.Set("target_virtual_machine_size", pointer.From(a2aDetails.RecoveryAzureVMSize))
			d.Set("target_zone", a2aDetails.RecoveryAvailabilityZone)
			d.Set("target_edge_zone", helpers.FlattenEdgeZone(a2aDetails.RecoveryExtendedLocation))
			d.Set("multi_vm_group_name", a2aDetails.MultiVMGroupName)
			d.Set("test_network_id", a2aDetails.SelectedTfoAzureNetworkId)
			if a2aDetails.ProtectedDisks != nil {
				disksOutput := make([]any, 0)
				for _, disk := range *a2aDetails.ProtectedDisks {
					disksOutput = append(disksOutput, map[string]any{
						"disk_uri":                   disk.DiskUri,
						"staging_storage_account_id": disk.PrimaryStagingAzureStorageAccountId,
						"target_storage_account_id":  disk.RecoveryAzureStorageAccountId,
					})
				}
				d.Set("unmanaged_disk", disksOutput)
			}

			if a2aDetails.ProtectedManagedDisks != nil {
				disksOutput := make([]any, 0)
				for _, disk := range *a2aDetails.ProtectedManagedDisks {
					diskOutput := make(map[string]any)
					diskId := ""
					if respDiskId := pointer.From(disk.DiskId); respDiskId != "" {
						parsedDiskId, err := commonids.ParseManagedDiskIDInsensitively(respDiskId)
						if err != nil {
							return err
						}
						diskId = parsedDiskId.ID()
					}
					diskOutput["disk_id"] = diskId

					primaryStagingAzureStorageAccountID := ""
					if respStorageAccId := pointer.From(disk.PrimaryStagingAzureStorageAccountId); respStorageAccId != "" {
						parsedStorageAccountId, err := commonids.ParseStorageAccountIDInsensitively(respStorageAccId)
						if err != nil {
							return err
						}
						primaryStagingAzureStorageAccountID = parsedStorageAccountId.ID()
					}
					diskOutput["staging_storage_account_id"] = primaryStagingAzureStorageAccountID

					recoveryResourceGroupID := ""
					if respRGId := pointer.From(disk.RecoveryResourceGroupId); respRGId != "" {
						parsedResourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(respRGId)
						if err != nil {
							return err
						}
						recoveryResourceGroupID = parsedResourceGroupId.ID()
					}
					diskOutput["target_resource_group_id"] = recoveryResourceGroupID

					diskOutput["target_replica_disk_type"] = pointer.From(disk.RecoveryReplicaDiskAccountType)

					diskOutput["target_disk_type"] = pointer.From(disk.RecoveryTargetDiskAccountType)

					recoveryEncryptionSetId := ""
					if respDESId := pointer.From(disk.RecoveryDiskEncryptionSetId); respDESId != "" {
						parsedEncryptionSetId, err := commonids.ParseDiskEncryptionSetIDInsensitively(respDESId)
						if err != nil {
							return err
						}
						recoveryEncryptionSetId = parsedEncryptionSetId.ID()
					}
					diskOutput["target_disk_encryption_set_id"] = recoveryEncryptionSetId

					diskOutput["target_disk_encryption"] = flattenTargetDiskEncryption(disk)

					disksOutput = append(disksOutput, diskOutput)
				}
				d.Set("managed_disk", pluginsdk.NewSet(resourceSiteRecoveryReplicatedVMDiskHash, disksOutput))
			}

			if a2aDetails.VMNics != nil {
				nicsOutput := make([]any, 0)
				for _, nic := range *a2aDetails.VMNics {
					nicOutput := make(map[string]any)
					if nic.SourceNicArmId != nil {
						nicOutput["source_network_interface_id"] = *nic.SourceNicArmId
					}
					if nic.IPConfigs != nil && len(*nic.IPConfigs) > 0 {
						nicOutput["ip_configuration"] = flattenSiteRecoveryReplicatedVMIPConfig(nic.IPConfigs)
					}
					nicsOutput = append(nicsOutput, nicOutput)
				}

				if err := d.Set("network_interface", nicsOutput); err != nil {
					return fmt.Errorf("setting `network_interface`: %v", err)
				}
			}
		}
	}

	return nil
}

func flattenSiteRecoveryReplicatedVMIPConfig(ipConfigs *[]replicationprotecteditems.IPConfigDetails) []any {
	outputs := make([]any, 0)

	if ipConfigs != nil {
		for _, ipConfig := range *ipConfigs {
			output := map[string]any{
				"name":                               pointer.From(ipConfig.Name),
				"primary":                            pointer.From(ipConfig.IsPrimary),
				"target_static_ip":                   pointer.From(ipConfig.RecoveryStaticIPAddress),
				"target_subnet_name":                 pointer.From(ipConfig.RecoverySubnetName),
				"recovery_public_ip_address_id":      pointer.From(ipConfig.RecoveryPublicIPAddressId),
				"failover_test_static_ip":            pointer.From(ipConfig.TfoStaticIPAddress),
				"failover_test_subnet_name":          pointer.From(ipConfig.TfoSubnetName),
				"failover_test_public_ip_address_id": pointer.From(ipConfig.TfoPublicIPAddressId),
			}
			if ipConfig.RecoveryLBBackendAddressPoolIds != nil {
				output["recovery_load_balancer_backend_address_pool_ids"] = pluginsdk.FlattenSlice(ipConfig.RecoveryLBBackendAddressPoolIds)
			}
			outputs = append(outputs, output)
		}
	}
	return outputs
}

func flattenTargetDiskEncryption(disk replicationprotecteditems.A2AProtectedManagedDiskDetails) []any {
	secretUrl := ""
	dekVaultId := ""
	keyUrl := ""
	kekVaultId := ""

	if disk.SecretIdentifier != nil {
		secretUrl = *disk.SecretIdentifier
	}
	if disk.DekKeyVaultArmId != nil {
		dekVaultId = *disk.DekKeyVaultArmId
	}
	if disk.KeyIdentifier != nil {
		keyUrl = *disk.KeyIdentifier
	}
	if disk.KekKeyVaultArmId != nil {
		kekVaultId = *disk.KekKeyVaultArmId
	}

	if secretUrl == "" && dekVaultId == "" && keyUrl == "" && kekVaultId == "" {
		return []any{}
	}

	diskEncryptionKeys := make([]any, 0)
	if secretUrl != "" || dekVaultId != "" {
		diskEncryptionKeys = append(diskEncryptionKeys, map[string]any{
			"secret_url": secretUrl,
			"vault_id":   dekVaultId,
		})
	}

	keyEncryptionKeys := make([]any, 0)
	if keyUrl != "" || kekVaultId != "" {
		keyEncryptionKeys = append(keyEncryptionKeys, map[string]any{
			"key_url":  keyUrl,
			"vault_id": kekVaultId,
		})
	}

	return []any{
		map[string]any{
			"disk_encryption_key": diskEncryptionKeys,
			"key_encryption_key":  keyEncryptionKeys,
		},
	}
}
