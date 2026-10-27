// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	//nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient
			fabricClient := metadata.Client.RecoveryServices.FabricClient
			containerClient := metadata.Client.RecoveryServices.ProtectionContainerClient
			vmwareMachinesClient := metadata.Client.RecoveryServices.VMWareMachinesClient
			runAsAccountsClient := metadata.Client.RecoveryServices.VMWareRunAsAccountsClient

			var model SiteRecoveryReplicatedVmVMwareModel
			err := metadata.Decode(&model)
			if err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			vaultId, err := replicationprotecteditems.ParseVaultID(model.RecoveryVaultId)
			if err != nil {
				return fmt.Errorf("parsing vault id %q: %+v", model.RecoveryVaultId, err)
			}

			containerId, err := fetchSiteRecoveryContainerId(ctx, containerClient, vaultId.ID())
			if err != nil {
				return fmt.Errorf("fetch Replication Container from vault %q: %+v", vaultId, err)
			}

			parsedContainerId, err := replicationprotecteditems.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parse %s: %+v", containerId, err)
			}

			id := replicationprotecteditems.NewReplicationProtectedItemID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, model.Name)
			fabricId := replicationfabrics.NewReplicationFabricID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing site recovery vmware replicated vm %q: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_site_recovery_vmware_replicated_vm", *existing.Model.Id)
				}
			}

			processServerId, err := fetchProcessServerIdByName(ctx, fabricClient, fabricId, model.ApplianceName)
			if err != nil {
				return fmt.Errorf("fetch process server id: %+v", err)
			}

			siteID, err := fetchVmwareSiteIdByFabric(ctx, fabricClient, fabricId)
			if err != nil {
				return fmt.Errorf("fetch VMWare site id: %+v", err)
			}

			discoveryMachineId, err := fetchDiscoveryMachineIdBySite(ctx, vmwareMachinesClient, siteID, model.SourceVmName)
			if err != nil {
				return fmt.Errorf("fetch discovery machine id %s: %+v", model.SourceVmName, err)
			}

			runAsAccountId, err := fetchRunAsAccountsIdBySite(ctx, runAsAccountsClient, siteID, model.PhysicalServerCredentialName, model.ApplianceName)
			if err != nil {
				return fmt.Errorf("fetch run as account id %s: %+v", model.PhysicalServerCredentialName, err)
			}

			providerSpecificDetail := replicationprotecteditems.InMageRcmEnableProtectionInput{
				LicenseType:              pointer.ToEnum[replicationprotecteditems.LicenseType](model.LicenseType),
				TargetVMName:             &model.TargetVmName,
				TargetResourceGroupId:    model.TargetResourceGroupId,
				FabricDiscoveryMachineId: discoveryMachineId,
				RunAsAccountId:           &runAsAccountId,
			}

			diskDefaultValueSet := false
			diskDefaultValue := replicationprotecteditems.InMageRcmDisksDefaultInput{}

			if model.DefaultRecoveryDiskType != "" {
				diskDefaultValueSet = true
				diskDefaultValue.DiskType = replicationprotecteditems.DiskAccountType(model.DefaultRecoveryDiskType)
			}

			if model.DefaultLogStorageAccountId != "" {
				diskDefaultValueSet = true
				diskDefaultValue.LogStorageAccountId = model.DefaultLogStorageAccountId
			}

			if model.DefaultTargetDiskEncryptionSetId != "" {
				diskDefaultValueSet = true
				diskDefaultValue.DiskEncryptionSetId = &model.DefaultTargetDiskEncryptionSetId
			}

			if diskDefaultValueSet {
				providerSpecificDetail.DisksDefault = &diskDefaultValue
			}

			if model.MultiVmGroupName != "" {
				providerSpecificDetail.MultiVMGroupName = &model.MultiVmGroupName
			}

			if model.ApplianceName != "" {
				providerSpecificDetail.ProcessServerId = processServerId
			}

			if model.TargetVmSize != "" {
				providerSpecificDetail.TargetVMSize = &model.TargetVmSize
			}

			if model.TargetAvailabilitySetId != "" {
				providerSpecificDetail.TargetAvailabilitySetId = pointer.To(model.TargetAvailabilitySetId)
			}

			if model.TargetZone != "" {
				providerSpecificDetail.TargetAvailabilityZone = pointer.To(model.TargetZone)
			}

			if model.TargetBootDiagnosticsStorageAccountId != "" {
				providerSpecificDetail.TargetBootDiagnosticsStorageAccountId = pointer.To(model.TargetBootDiagnosticsStorageAccountId)
			}

			if model.TargetProximityPlacementGroupId != "" {
				providerSpecificDetail.TargetProximityPlacementGroupId = &model.TargetProximityPlacementGroupId
			}

			if model.TargetNetworkId != "" {
				providerSpecificDetail.TargetNetworkId = &model.TargetNetworkId
			}

			if model.TestNetworkId != "" {
				providerSpecificDetail.TestNetworkId = &model.TestNetworkId
			}

			diskToIncludeOutput := make([]replicationprotecteditems.InMageRcmDiskInput, 0)
			for _, diskRaw := range model.DiskToInclude {
				diskToIncludeOutput = append(diskToIncludeOutput, replicationprotecteditems.InMageRcmDiskInput{
					DiskEncryptionSetId: &diskRaw.TargetDiskEncryptionSetId,
					DiskId:              diskRaw.DiskId,
					DiskType:            replicationprotecteditems.DiskAccountType(diskRaw.TargetDiskType),
					LogStorageAccountId: diskRaw.LogStorageAccountId,
				})
			}

			if len(diskToIncludeOutput) > 0 {
				providerSpecificDetail.DisksToInclude = &diskToIncludeOutput
			}

			parameters := replicationprotecteditems.EnableProtectionInput{
				Properties: &replicationprotecteditems.EnableProtectionInputProperties{
					PolicyId:                &model.RecoveryReplicationPolicyId,
					ProviderSpecificDetails: providerSpecificDetail,
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}
			metadata.SetID(id)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Pending"},
				Target:  []string{"Protected"},
				Refresh: func() (result any, state string, err error) {
					resp, err := client.Get(ctx, id)
					if err != nil {
						return nil, "error", fmt.Errorf("retrieving %s: %+v", id, err)
					}

					protectionState := ""
					if model := resp.Model; model != nil && model.Properties != nil && resp.Model.Properties.ProtectionState != nil {
						protectionState = *model.Properties.ProtectionState
					}

					if strings.EqualFold(protectionState, "Protected") {
						return resp, protectionState, nil
					}

					// The `protectionState` has pretty much enums, and will changes in the duration of replicate.
					// While failed ones and canceled ones have common pattern.
					if strings.HasSuffix(protectionState, "Failed") || strings.HasPrefix(protectionState, "Cancel") {
						return resp, protectionState, fmt.Errorf("replicate failed or canceled")
					}

					return resp, "Pending", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully protected: %s", id, err)
			}

			// it needs an additional update to set the network interface.
			updateInput := replicationprotecteditems.UpdateReplicationProtectedItemInput{
				Properties: &replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
					RecoveryAvailabilitySetId:      providerSpecificDetail.TargetAvailabilitySetId,
					RecoveryAzureVMName:            providerSpecificDetail.TargetVMName,
					RecoveryAzureVMSize:            providerSpecificDetail.TargetVMSize,
					SelectedRecoveryAzureNetworkId: providerSpecificDetail.TargetNetworkId,
					SelectedTfoAzureNetworkId:      providerSpecificDetail.TestNetworkId,
					ProviderSpecificDetails: replicationprotecteditems.InMageRcmUpdateReplicationProtectedItemInput{
						LicenseType:                           providerSpecificDetail.LicenseType,
						TargetAvailabilityZone:                providerSpecificDetail.TargetAvailabilityZone,
						TargetAvailabilitySetId:               providerSpecificDetail.TargetAvailabilitySetId,
						TargetBootDiagnosticsStorageAccountId: providerSpecificDetail.TargetBootDiagnosticsStorageAccountId,
						TargetNetworkId:                       providerSpecificDetail.TargetNetworkId,
						TargetProximityPlacementGroupId:       providerSpecificDetail.TargetProximityPlacementGroupId,
						TargetResourceGroupId:                 &providerSpecificDetail.TargetResourceGroupId,
						TargetVMName:                          providerSpecificDetail.TargetVMName,
						TargetVMSize:                          providerSpecificDetail.TargetVMSize,
						TestNetworkId:                         providerSpecificDetail.TestNetworkId,
						VMNics:                                pointer.To(expandVMWareReplicatedVMNics(model.NetworkInterface)),
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, id, updateInput); err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			return nil
		},
	}
}

func expandVMWareReplicatedVMNics(input []NetworkInterfaceModel) []replicationprotecteditems.InMageRcmNicInput {
	vmNics := make([]replicationprotecteditems.InMageRcmNicInput, 0)
	for _, nic := range input {
		vmNic := replicationprotecteditems.InMageRcmNicInput{
			NicId:            nic.SourceMacAddress,
			TargetSubnetName: &nic.TargetSubnetName,
			TestSubnetName:   &nic.TestSubnetName,
		}
		if nic.TargetStaticIp != "" {
			vmNic.TargetStaticIPAddress = &nic.TargetStaticIp
		}
		if nic.IsPrimary {
			vmNic.IsPrimaryNic = strconv.FormatBool(true)
			vmNic.IsSelectedForFailover = pointer.To("true")
		} else {
			vmNic.IsPrimaryNic = strconv.FormatBool(false)
			vmNic.IsSelectedForFailover = pointer.To("false")
		}
		vmNics = append(vmNics, vmNic)
	}
	return vmNics
}
