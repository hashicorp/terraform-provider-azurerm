// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicatedItemUpdate(d *pluginsdk.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	// We are only allowed to update the configuration once the VM is fully protected
	state, err := waitForReplicationToBeHealthy(ctx, d, meta)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	testNetworkId := d.Get("test_network_id").(string)

	id := replicationprotecteditems.NewReplicationProtectedItemID(subscriptionId, resGroup, vaultName, fabricName, sourceProtectionContainerName, name)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		targetAvailabilitySetID = pointer.To(id.(string))
	} else {
		targetAvailabilitySetID = nil
	}

	nicList := d.Get("network_interface").([]any)
	vmNics := make([]replicationprotecteditems.VMNicInputDetails, 0, len(nicList))
	for _, raw := range nicList {
		vmNicInput := raw.(map[string]any)
		sourceNicId := vmNicInput["source_network_interface_id"].(string)
		nicId := findNicId(state, sourceNicId)
		if nicId == nil {
			return fmt.Errorf("updating replicated vm %s (vault %s): Trying to update NIC that is not known by Azure %s", name, vaultName, sourceNicId)
		}
		ipConfig := expandSiteRecoveryReplicatedVMIPConfig(vmNicInput)
		vmNics = append(vmNics, replicationprotecteditems.VMNicInputDetails{
			NicId:     nicId,
			IPConfigs: &ipConfig,
		})
	}

	oldManagedDisks, newManagedDisks := d.GetChange("managed_disk")
	oldSet := oldManagedDisks.(*pluginsdk.Set)
	newSet := newManagedDisks.(*pluginsdk.Set)

	addedDisks := newSet.Difference(oldSet).List()
	existingDisks := newSet.Intersection(oldSet).List()

	if len(addedDisks) > 0 {
		var vmManagedDisks []replicationprotecteditems.A2AVMManagedDiskInputDetails
		for _, raw := range addedDisks {
			diskInput := raw.(map[string]any)
			diskId := diskInput["disk_id"].(string)
			stagingStorageAccountId := diskInput["staging_storage_account_id"].(string)
			targetResourceGroupId := diskInput["target_resource_group_id"].(string)

			diskInputDetails := replicationprotecteditems.A2AVMManagedDiskInputDetails{
				DiskId:                              diskId,
				PrimaryStagingAzureStorageAccountId: stagingStorageAccountId,
				RecoveryResourceGroupId:             targetResourceGroupId,
			}

			if v := diskInput["target_replica_disk_type"].(string); v != "" {
				diskInputDetails.RecoveryReplicaDiskAccountType = &v
			}
			if v := diskInput["target_disk_type"].(string); v != "" {
				diskInputDetails.RecoveryTargetDiskAccountType = &v
			}
			if v := diskInput["target_disk_encryption_set_id"].(string); v != "" {
				diskInputDetails.RecoveryDiskEncryptionSetId = &v
			}

			diskEncryptionRaw := diskInput["target_disk_encryption"].([]any)
			if len(diskEncryptionRaw) > 0 {
				diskInputDetails.DiskEncryptionInfo = expandDiskEncryption(diskEncryptionRaw)
			}

			vmManagedDisks = append(vmManagedDisks, diskInputDetails)
		}

		if err = client.AddDisksThenPoll(ctx, id, replicationprotecteditems.AddDisksInput{
			Properties: &replicationprotecteditems.AddDisksInputProperties{
				ProviderSpecificDetails: replicationprotecteditems.A2AAddDisksInput{
					VMManagedDisks: &vmManagedDisks,
				},
			},
		}); err != nil {
			return fmt.Errorf("adding disks to replicated vm %s (vault %s): %+v", name, vaultName, err)
		}
	}

	managedDisks := make([]replicationprotecteditems.A2AVMManagedDiskUpdateDetails, 0, len(existingDisks))
	for _, raw := range existingDisks {
		diskInput := raw.(map[string]any)

		managedDisks = append(managedDisks, replicationprotecteditems.A2AVMManagedDiskUpdateDetails{
			DiskId:                         pointer.To(diskInput["disk_id"].(string)),
			RecoveryReplicaDiskAccountType: pointer.To(diskInput["target_replica_disk_type"].(string)),
			RecoveryTargetDiskAccountType:  pointer.To(diskInput["target_disk_type"].(string)),
			DiskEncryptionInfo:             expandDiskEncryption(diskInput["target_disk_encryption"].([]any)),
		})
	}

	if targetNetworkId == "" {
		// No target network id was specified, so we want to preserve what was selected
		if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.SelectedRecoveryAzureNetworkId != nil {
				targetNetworkId = *a2aDetails.SelectedRecoveryAzureNetworkId
			} else {
				return fmt.Errorf("target_network_id must be set when a network_interface is configured")
			}
		} else {
			return fmt.Errorf("target_network_id must be set when a network_interface is configured")
		}
	}

	if testNetworkId == "" {
		// No test network id was specified, so we want to preserve what was selected
		if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.SelectedTfoAzureNetworkId != nil {
				testNetworkId = *a2aDetails.SelectedTfoAzureNetworkId
			}
		}
	}

	parameters := replicationprotecteditems.UpdateReplicationProtectedItemInput{
		Properties: &replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
			RecoveryAzureVMName:            &name,
			SelectedRecoveryAzureNetworkId: &targetNetworkId,
			SelectedTfoAzureNetworkId:      &testNetworkId,
			VMNics:                         &vmNics,
			RecoveryAvailabilitySetId:      targetAvailabilitySetID,
			RecoveryAzureVMSize:            pointer.To(d.Get("target_virtual_machine_size").(string)),
			ProviderSpecificDetails: replicationprotecteditems.A2AUpdateReplicationProtectedItemInput{
				ManagedDiskUpdateDetails:           &managedDisks,
				RecoveryProximityPlacementGroupId:  pointer.To(d.Get("target_proximity_placement_group_id").(string)),
				RecoveryBootDiagStorageAccountId:   pointer.To(d.Get("target_boot_diagnostic_storage_account_id").(string)),
				RecoveryCapacityReservationGroupId: pointer.To(d.Get("target_capacity_reservation_group_id").(string)),
				RecoveryVirtualMachineScaleSetId:   pointer.To(d.Get("target_virtual_machine_scale_set_id").(string)),
			},
		},
	}

	if err = client.UpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching updated replicated vm %s (vault %s): %+v", name, vaultName, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("fetching updated replicated vm %s: model is nil", name)
	}

	return flattenSiteRecoveryReplicatedItem(d, resp.Model)
}
