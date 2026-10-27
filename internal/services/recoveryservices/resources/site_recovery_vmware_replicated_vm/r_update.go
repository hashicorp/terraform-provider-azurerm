// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			var model SiteRecoveryReplicatedVmVMwareModel
			err := metadata.Decode(&model)
			if err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", metadata.ResourceData.Id(), err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: Model was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: Properties was nil", id)
			}
			if existing.Model.Properties.ProviderSpecificDetails == nil {
				return fmt.Errorf("retrieving %s: ProviderSpecificDetails was nil", id)
			}
			if _, ok := existing.Model.Properties.ProviderSpecificDetails.(replicationprotecteditems.InMageRcmReplicationDetails); !ok {
				return fmt.Errorf("retrieving %s: ProviderSpecificDetails was not InMageRcmProtectedItemDetails", id)
			}

			existingProps := *existing.Model.Properties
			existingDetails := existingProps.ProviderSpecificDetails.(replicationprotecteditems.InMageRcmReplicationDetails)

			vmNics := make([]replicationprotecteditems.InMageRcmNicInput, 0)
			if metadata.ResourceData.HasChange("network_interface") {
				vmNics = expandVMWareReplicatedVMNics(model.NetworkInterface)
			} else {
				if existingDetails.VMNics == nil {
					return errors.New("retrieving `network_interface`: VMNics was nil")
				} else {
					for _, respNic := range *existingDetails.VMNics {
						vmNics = append(vmNics, replicationprotecteditems.InMageRcmNicInput{
							IsPrimaryNic:          pointer.From(respNic.IsPrimaryNic),
							IsSelectedForFailover: respNic.IsSelectedForFailover,
							NicId:                 pointer.From(respNic.NicId),
							TargetStaticIPAddress: respNic.TargetIPAddress,
							TargetSubnetName:      respNic.TargetSubnetName,
							TestStaticIPAddress:   respNic.TestIPAddress,
							TestSubnetName:        respNic.TestSubnetName,
						})
					}
				}
			}

			updateInput := replicationprotecteditems.InMageRcmUpdateReplicationProtectedItemInput{
				VMNics: &vmNics,
			}

			if metadata.ResourceData.HasChange("license_type") {
				updateInput.LicenseType = pointer.ToEnum[replicationprotecteditems.LicenseType](model.LicenseType)
			} else if existingDetails.LicenseType != nil {
				updateInput.LicenseType = pointer.ToEnum[replicationprotecteditems.LicenseType](*existingDetails.LicenseType)
			}

			if metadata.ResourceData.HasChange("target_vm_name") {
				updateInput.TargetVMName = &model.TargetVmName
			} else {
				updateInput.TargetVMName = existingDetails.TargetVMName
			}

			if metadata.ResourceData.HasChange("target_resource_group_id") {
				updateInput.TargetResourceGroupId = &model.TargetResourceGroupId
			} else {
				updateInput.TargetResourceGroupId = existingDetails.TargetResourceGroupId
			}

			if metadata.ResourceData.HasChange("target_availability_set_id") {
				if model.TargetAvailabilitySetId != "" {
					updateInput.TargetAvailabilitySetId = &model.TargetAvailabilitySetId
				} else {
					updateInput.TargetAvailabilitySetId = nil
				}
			} else {
				updateInput.TargetAvailabilitySetId = existingDetails.TargetAvailabilitySetId
			}

			if metadata.ResourceData.HasChange("target_zone") {
				if model.TargetZone != "" {
					updateInput.TargetAvailabilityZone = &model.TargetZone
				} else {
					updateInput.TargetAvailabilityZone = nil
				}
			} else {
				updateInput.TargetAvailabilityZone = existingDetails.TargetAvailabilityZone
			}

			if metadata.ResourceData.HasChange("target_network_id") {
				updateInput.TargetNetworkId = &model.TargetNetworkId
			} else {
				updateInput.TargetNetworkId = existingDetails.TargetNetworkId
			}

			if metadata.ResourceData.HasChange("target_proximity_placement_group_id") {
				updateInput.TargetProximityPlacementGroupId = &model.TargetProximityPlacementGroupId
			} else {
				updateInput.TargetProximityPlacementGroupId = existingDetails.TargetProximityPlacementGroupId
			}

			if metadata.ResourceData.HasChange("target_boot_diagnostics_storage_account_id") {
				updateInput.TargetBootDiagnosticsStorageAccountId = &model.TargetBootDiagnosticsStorageAccountId
			} else {
				updateInput.TargetBootDiagnosticsStorageAccountId = existingDetails.TargetBootDiagnosticsStorageAccountId
			}

			props := replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
				ProviderSpecificDetails: updateInput,
			}

			if metadata.ResourceData.HasChange("target_availability_set_id") {
				props.RecoveryAvailabilitySetId = updateInput.TargetAvailabilitySetId
			} else {
				props.RecoveryAvailabilitySetId = existingDetails.TargetAvailabilitySetId
			}

			if metadata.ResourceData.HasChange("target_vm_name") {
				props.RecoveryAzureVMName = &model.TargetVmName
			} else {
				props.RecoveryAzureVMName = existingDetails.TargetVMName
			}

			if metadata.ResourceData.HasChange("target_network_id") {
				props.SelectedRecoveryAzureNetworkId = &model.TargetNetworkId
			} else {
				props.SelectedRecoveryAzureNetworkId = existingDetails.TargetNetworkId
			}

			if metadata.ResourceData.HasChange("target_vm_size") {
				props.RecoveryAzureVMSize = &model.TargetVmSize
			} else {
				props.RecoveryAzureVMSize = existingDetails.TargetVMSize
			}

			parameters := replicationprotecteditems.UpdateReplicationProtectedItemInput{
				Properties: &props,
			}

			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}
