// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"fmt"

	//nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicatedItemCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	name := d.Get("name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceVmId := d.Get("source_vm_id").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetProtectionContainerId := d.Get("target_recovery_protection_container_id").(string)
	targetResourceGroupId := d.Get("target_resource_group_id").(string)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		targetAvailabilitySetID = pointer.To(id.(string))
	} else {
		targetAvailabilitySetID = nil
	}

	var targetAvailabilityZone *string
	if zone, isSet := d.GetOk("target_zone"); isSet {
		targetAvailabilityZone = pointer.To(zone.(string))
	} else {
		targetAvailabilityZone = nil
	}

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationprotecteditems.NewReplicationProtectedItemID(subscriptionId, resGroup, vaultName, fabricName, sourceProtectionContainerName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replicated_vm", *existing.Model.Id)
		}
	}

	managedDisksGet := d.Get("managed_disk").(*pluginsdk.Set).List()
	managedDisks := make([]replicationprotecteditems.A2AVMManagedDiskInputDetails, 0, len(managedDisksGet))
	for _, raw := range managedDisksGet {
		diskInput := raw.(map[string]any)
		managedDisks = append(managedDisks, replicationprotecteditems.A2AVMManagedDiskInputDetails{
			DiskId:                              diskInput["disk_id"].(string),
			PrimaryStagingAzureStorageAccountId: diskInput["staging_storage_account_id"].(string),
			RecoveryResourceGroupId:             diskInput["target_resource_group_id"].(string),
			RecoveryReplicaDiskAccountType:      pointer.To(diskInput["target_replica_disk_type"].(string)),
			RecoveryTargetDiskAccountType:       pointer.To(diskInput["target_disk_type"].(string)),
			RecoveryDiskEncryptionSetId:         pointer.To(diskInput["target_disk_encryption_set_id"].(string)),
			DiskEncryptionInfo:                  expandDiskEncryption(diskInput["target_disk_encryption"].([]any)),
		})
	}

	vmDisksGet := d.Get("unmanaged_disk").(*pluginsdk.Set).List()
	vmDisks := make([]replicationprotecteditems.A2AVMDiskInputDetails, 0, len(vmDisksGet))
	for _, raw := range vmDisksGet {
		diskInput := raw.(map[string]any)
		vmDisks = append(vmDisks, replicationprotecteditems.A2AVMDiskInputDetails{
			DiskUri:                             diskInput["disk_uri"].(string),
			PrimaryStagingAzureStorageAccountId: diskInput["staging_storage_account_id"].(string),
			RecoveryAzureStorageAccountId:       diskInput["target_storage_account_id"].(string),
		})
	}

	parameters := replicationprotecteditems.EnableProtectionInput{
		Properties: &replicationprotecteditems.EnableProtectionInputProperties{
			PolicyId: &policyId,
			ProviderSpecificDetails: replicationprotecteditems.A2AEnableProtectionInput{
				FabricObjectId:                     sourceVmId,
				RecoveryContainerId:                &targetProtectionContainerId,
				RecoveryResourceGroupId:            &targetResourceGroupId,
				RecoveryAvailabilitySetId:          targetAvailabilitySetID,
				RecoveryAvailabilityZone:           targetAvailabilityZone,
				MultiVMGroupName:                   pointer.To(d.Get("multi_vm_group_name").(string)),
				RecoveryProximityPlacementGroupId:  pointer.To(d.Get("target_proximity_placement_group_id").(string)),
				RecoveryBootDiagStorageAccountId:   pointer.To(d.Get("target_boot_diagnostic_storage_account_id").(string)),
				RecoveryCapacityReservationGroupId: pointer.To(d.Get("target_capacity_reservation_group_id").(string)),
				RecoveryVirtualMachineScaleSetId:   pointer.To(d.Get("target_virtual_machine_scale_set_id").(string)),
				VMManagedDisks:                     &managedDisks,
				VMDisks:                            &vmDisks,
				RecoveryExtendedLocation:           helpers.ExpandEdgeZone(d.Get("target_edge_zone").(string)),
			},
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}
	d.SetId(id.ID())

	// We are not allowed to configure the NIC on the initial setup, and the VM has to be replicated before
	// we can reconfigure. Hence this call to update when we create.
	state, err := waitForReplicationToBeHealthy(ctx, d, meta)
	if err != nil {
		return err
	}

	targetNetworkId := d.Get("target_network_id").(string)
	testNetworkId := d.Get("test_network_id").(string)

	if targetNetworkId == "" {
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
		if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.SelectedTfoAzureNetworkId != nil {
				testNetworkId = *a2aDetails.SelectedTfoAzureNetworkId
			}
		}
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

	updateParameters := replicationprotecteditems.UpdateReplicationProtectedItemInput{
		Properties: &replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
			RecoveryAzureVMName:            &name,
			SelectedRecoveryAzureNetworkId: &targetNetworkId,
			SelectedTfoAzureNetworkId:      &testNetworkId,
			VMNics:                         &vmNics,
			RecoveryAvailabilitySetId:      targetAvailabilitySetID,
			RecoveryAzureVMSize:            pointer.To(d.Get("target_virtual_machine_size").(string)),
			ProviderSpecificDetails: replicationprotecteditems.A2AUpdateReplicationProtectedItemInput{
				RecoveryProximityPlacementGroupId:  pointer.To(d.Get("target_proximity_placement_group_id").(string)),
				RecoveryBootDiagStorageAccountId:   pointer.To(d.Get("target_boot_diagnostic_storage_account_id").(string)),
				RecoveryCapacityReservationGroupId: pointer.To(d.Get("target_capacity_reservation_group_id").(string)),
				RecoveryVirtualMachineScaleSetId:   pointer.To(d.Get("target_virtual_machine_scale_set_id").(string)),
			},
		},
	}

	if err = client.UpdateThenPoll(ctx, id, updateParameters); err != nil {
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

func expandSiteRecoveryReplicatedVMIPConfig(nicInput map[string]any) []replicationprotecteditems.IPConfigInputDetails {
	output := make([]replicationprotecteditems.IPConfigInputDetails, 0)
	ipConfigs := nicInput["ip_configuration"].([]any)
	if len(ipConfigs) > 0 {
		for _, ipConfig := range ipConfigs {
			ipConfig := ipConfig.(map[string]any)
			var ipConfigName *string
			if name, ok := ipConfig["name"].(string); ok && name != "" {
				ipConfigName = pointer.To(name)
			}
			var recoveryLoadBalancerBackendPoolIds *[]string
			if ids, ok := ipConfig["recovery_load_balancer_backend_address_pool_ids"].(*schema.Set); ok && ids.Len() > 0 {
				recoveryLoadBalancerBackendPoolIds = pluginsdk.ExpandStringSlice(ids.List())
			}
			output = append(output, replicationprotecteditems.IPConfigInputDetails{
				IPConfigName:                    ipConfigName,
				RecoverySubnetName:              pointer.To(ipConfig["target_subnet_name"].(string)),
				RecoveryStaticIPAddress:         pointer.To(ipConfig["target_static_ip"].(string)),
				RecoveryLBBackendAddressPoolIds: recoveryLoadBalancerBackendPoolIds,
				RecoveryPublicIPAddressId:       pointer.To(ipConfig["recovery_public_ip_address_id"].(string)),
				TfoStaticIPAddress:              pointer.To(ipConfig["failover_test_static_ip"].(string)),
				TfoSubnetName:                   pointer.To(ipConfig["failover_test_subnet_name"].(string)),
				TfoPublicIPAddressId:            pointer.To(ipConfig["failover_test_public_ip_address_id"].(string)),
				IsPrimary:                       pointer.To(len(ipConfigs) == 1 || ipConfig["primary"].(bool)),
			})
		}
		return output
	}

	return output
}

func expandDiskEncryption(diskEncryptionInfoList []any) *replicationprotecteditems.DiskEncryptionInfo {
	if len(diskEncryptionInfoList) == 0 {
		return &replicationprotecteditems.DiskEncryptionInfo{}
	}
	diskEncryptionInfoMap := diskEncryptionInfoList[0].(map[string]any)

	dek := diskEncryptionInfoMap["disk_encryption_key"].([]any)[0].(map[string]any)
	diskEncryptionInfo := &replicationprotecteditems.DiskEncryptionInfo{
		DiskEncryptionKeyInfo: &replicationprotecteditems.DiskEncryptionKeyInfo{
			SecretIdentifier:      pointer.To(dek["secret_url"].(string)),
			KeyVaultResourceArmId: pointer.To(dek["vault_id"].(string)),
		},
	}

	if keyEncryptionKey := diskEncryptionInfoMap["key_encryption_key"].([]any); len(keyEncryptionKey) > 0 {
		kek := keyEncryptionKey[0].(map[string]any)
		diskEncryptionInfo.KeyEncryptionKeyInfo = &replicationprotecteditems.KeyEncryptionKeyInfo{
			KeyIdentifier:         pointer.To(kek["key_url"].(string)),
			KeyVaultResourceArmId: pointer.To(kek["vault_id"].(string)),
		}
	}

	return diskEncryptionInfo
}
