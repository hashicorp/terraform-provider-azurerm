// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient
			runAsAccountId := metadata.Client.RecoveryServices.VMWareRunAsAccountsClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id.String(), err)
			}

			vaultId := replicationprotecteditems.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

			state := SiteRecoveryReplicatedVmVMwareModel{
				Name:            id.ReplicationProtectedItemName,
				RecoveryVaultId: vaultId.ID(),
			}

			if resp.Model != nil && resp.Model.Properties != nil {
				prop := *resp.Model.Properties

				state.SourceVmName = pointer.From(prop.FriendlyName)

				policyId := ""
				if respPolicyId := pointer.From(prop.PolicyId); respPolicyId != "" {
					parsedPolicyId, err := replicationpolicies.ParseReplicationPolicyIDInsensitively(respPolicyId)
					if err != nil {
						return fmt.Errorf("parse %q: %+v", respPolicyId, err)
					}
					policyId = parsedPolicyId.ID()
				}
				state.RecoveryReplicationPolicyId = policyId

				if inMageRcm, isV2A := prop.ProviderSpecificDetails.(replicationprotecteditems.InMageRcmReplicationDetails); isV2A {
					state.ApplianceName = pointer.From(inMageRcm.ProcessServerName)

					state.LicenseType = pointer.From(inMageRcm.LicenseType)

					state.MultiVmGroupName = pointer.From(inMageRcm.MultiVMGroupName)

					state.TargetProximityPlacementGroupId = pointer.From(inMageRcm.TargetProximityPlacementGroupId)

					state.TargetResourceGroupId = pointer.From(inMageRcm.TargetResourceGroupId)

					state.TargetVmName = pointer.From(inMageRcm.TargetVMName)

					state.TargetVmSize = pointer.From(inMageRcm.TargetVMSize)

					state.TargetAvailabilitySetId = pointer.From(inMageRcm.TargetAvailabilitySetId)

					state.TargetZone = pointer.From(inMageRcm.TargetAvailabilityZone)

					state.TargetNetworkId = pointer.From(inMageRcm.TargetNetworkId)

					state.TestNetworkId = pointer.From(inMageRcm.TestNetworkId)

					state.TargetBootDiagnosticsStorageAccountId = pointer.From(inMageRcm.TargetBootDiagnosticsStorageAccountId)

					credential := ""
					if inMageRcm.RunAsAccountId != nil && *inMageRcm.RunAsAccountId != "" {
						credential, err = fetchCredentialByRunAsAccountId(ctx, runAsAccountId, *inMageRcm.RunAsAccountId)
						if err != nil {
							return fmt.Errorf("retrieving credential by run as account id %q: %+v", *inMageRcm.RunAsAccountId, err)
						}
					}
					state.PhysicalServerCredentialName = credential

					if inMageRcm.ProtectedDisks != nil {
						diskOutputs := make([]IncludedDiskModel, 0)
						for _, diskRaw := range *inMageRcm.ProtectedDisks {
							diskModel := IncludedDiskModel{
								DiskId:         *diskRaw.DiskId,
								TargetDiskType: string(*diskRaw.DiskType),
							}
							if diskRaw.DiskEncryptionSetId != nil {
								diskModel.TargetDiskEncryptionSetId = *diskRaw.DiskEncryptionSetId
							}
							diskOutputs = append(diskOutputs, diskModel)
						}
						state.DiskToInclude = diskOutputs
					}

					networkInterfaceModel := NetworkInterfaceModel{}
					networkInterfaceModel.TargetStaticIp = pointer.From(inMageRcm.PrimaryNicIPAddress)

					if inMageRcm.VMNics != nil {
						nicsOutput := make([]NetworkInterfaceModel, 0)
						for _, nic := range *inMageRcm.VMNics {
							nicModel := NetworkInterfaceModel{}
							nicModel.TargetStaticIp = pointer.From(nic.TargetIPAddress)
							nicModel.TargetSubnetName = pointer.From(nic.TargetSubnetName)
							nicModel.TestSubnetName = pointer.From(nic.TestSubnetName)
							nicModel.SourceMacAddress = pointer.From(nic.NicId)
							nicModel.IsPrimary = pointer.From(nic.IsPrimaryNic) == "true"

							nicsOutput = append(nicsOutput, nicModel)
						}
						state.NetworkInterface = nicsOutput
					}
				}
			}

			var plan SiteRecoveryReplicatedVmVMwareModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			return metadata.Encode(&state)
		},
	}
}

func resourceSiteRecoveryVMWareReplicatedVMManagedDiskSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"disk_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"target_disk_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
			},

			"target_disk_encryption_set_id": commonschema.ResourceIDReferenceOptional(&commonids.DiskEncryptionSetId{}),

			"log_storage_account_id": commonschema.ResourceIDReferenceOptional(&commonids.StorageAccountId{}),
		},
	}
}

func resourceSiteRecoveryVMWareReplicatedVMNetworkInterfaceSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_mac_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$`), "The `source_mac_address` must be in format `00:00:00:00:00:00`"),
			},

			"target_static_ip": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"target_subnet_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"test_subnet_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"is_primary": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func fetchSiteRecoveryContainerId(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, vaultId string) (string, error) {
	vId, err := replicationprotectioncontainers.ParseVaultID(vaultId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", vaultId, err)
	}

	resp, err := containerClient.ListComplete(ctx, *vId)
	if err != nil {
		return "", err
	}

	if len(resp.Items) != 1 {
		return "", fmt.Errorf("there should be only 1 protection container in Recovery Vault, get: %v", len(resp.Items))
	}

	parsedID, err := replicationprotectioncontainers.ParseReplicationProtectionContainerIDInsensitively(*resp.Items[0].Id)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", *resp.Items[0].Id, err)
	}

	return parsedID.ID(), nil
}

func fetchRunAsAccountsIdBySite(ctx context.Context, runAsAccountClient *runasaccounts.RunAsAccountsClient, siteId string, displayName string, applianceName string) (string, error) {
	parsedSiteId, err := runasaccounts.ParseVMwareSiteIDInsensitively(siteId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", siteId, err)
	}

	hackedClient := azuresdkhacks.RunAsAccountsClient{Client: runAsAccountClient.Client}
	// GET on Site is not working, tracked on https://github.com/Azure/azure-rest-api-specs/issues/24711
	resp, err := hackedClient.GetAllRunAsAccountsInSiteComplete(ctx, *parsedSiteId)
	if err != nil {
		return "", err
	}

	if len(resp.Items) == 0 {
		return "", fmt.Errorf("retire run as account from %s, get 0 item", siteId)
	}

	for _, account := range resp.Items {
		if account.Properties == nil {
			continue
		}
		if account.Properties.DisplayName == nil {
			continue
		}
		if account.Properties.ApplianceName == nil {
			continue
		}
		if strings.EqualFold(*account.Properties.DisplayName, displayName) && strings.EqualFold(*account.Properties.ApplianceName, applianceName) {
			return *account.Id, nil
		}
	}

	return "", fmt.Errorf("retrieving %q: run as account %s not found", siteId, displayName)
}

func fetchProcessServerIdByName(ctx context.Context, fabricClient *replicationfabrics.ReplicationFabricsClient, fabricId replicationfabrics.ReplicationFabricId, processServerName string) (string, error) {
	resp, err := fabricClient.Get(ctx, fabricId, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		return "", err
	}

	if resp.Model == nil {
		return "", fmt.Errorf("retrieving %q: Model is nil", fabricId)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving %q: Properties is nil", fabricId)
	}

	if resp.Model.Properties.CustomDetails == nil {
		return "", fmt.Errorf("retrieving %q: CustomDetails is nil", fabricId)
	}

	if detail, ok := resp.Model.Properties.CustomDetails.(replicationfabrics.InMageRcmFabricSpecificDetails); ok {
		if detail.ProcessServers == nil {
			return "", fmt.Errorf("retrieving %q: ProcessServers is nil", fabricId)
		}
		for _, server := range *detail.ProcessServers {
			if strings.EqualFold(*server.Name, processServerName) {
				return *server.Id, nil
			}
		}
		return "", fmt.Errorf("retrieving %q: process server %s not found", fabricId, processServerName)
	}
	return "", fmt.Errorf("retrieving %q: Detail Type mismatch", fabricId)
}

func fetchDiscoveryMachineIdBySite(ctx context.Context, machinesClient *machines.MachinesClient, siteId string, machineName string) (string, error) {
	parsedSiteId, err := machines.ParseVMwareSiteIDInsensitively(siteId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", siteId, err)
	}

	hackedClient := azuresdkhacks.MachinesClient{Client: machinesClient.Client}
	resp, err := hackedClient.GetAllVMWareMachinesInSite(ctx, *parsedSiteId, machines.DefaultGetAllMachinesInSiteOperationOptions())
	if err != nil {
		return "", err
	}

	if model := resp.Model; model != nil {
		for _, machine := range *model {
			if machine.Properties == nil {
				continue
			}
			if machine.Properties.DisplayName == nil {
				continue
			}
			if strings.EqualFold(*machine.Properties.DisplayName, machineName) {
				parsedMachineId, err := commonids.ParseVMwareSiteMachineIDInsensitively(*machine.Id)
				if err != nil {
					return "", fmt.Errorf("parse %s: %+v", *machine.Id, err)
				}
				return parsedMachineId.ID(), nil
			}
		}
	}

	return "", fmt.Errorf("retrieving %q: machine %s not found", siteId, machineName)
}

func fetchCredentialByRunAsAccountId(ctx context.Context, client *runasaccounts.RunAsAccountsClient, id string) (string, error) {
	parsedRunAsAccountId, err := commonids.ParseVMwareSiteRunAsAccountIDInsensitively(id)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", id, err)
	}

	resp, err := client.GetRunAsAccount(ctx, *parsedRunAsAccountId)
	if err != nil {
		return "", err
	}

	if resp.Model == nil {
		return "", fmt.Errorf("retrieving %q: Model was nil", id)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving %q: Properties was nil", id)
	}

	if resp.Model.Properties.DisplayName == nil {
		return "", fmt.Errorf("retrieving %q: DisplayName was nil", id)
	}

	return *resp.Model.Properties.DisplayName, nil
}

func fetchVmwareSiteIdByFabric(ctx context.Context, fabricClient *replicationfabrics.ReplicationFabricsClient, fabricId replicationfabrics.ReplicationFabricId) (string, error) {
	resp, err := fabricClient.Get(ctx, fabricId, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		return "", err
	}
	if resp.Model == nil {
		return "", fmt.Errorf("retrieving %q: Model is nil", fabricId)
	}
	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving %q: Properties is nil", fabricId)
	}
	if v, ok := resp.Model.Properties.CustomDetails.(replicationfabrics.InMageRcmFabricSpecificDetails); ok {
		if v.VMwareSiteId == nil {
			return "", fmt.Errorf("retrieving %q: VMwareSiteId is nil", fabricId)
		}
		return *v.VMwareSiteId, nil
	}
	return "", fmt.Errorf("retrieving %q: Detail Type mismatch", fabricId)
}
