// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	vmwaremachines "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
	vmwarerunasaccounts "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type IncludedDiskModel struct {
	DiskId                    string `tfschema:"disk_id"`
	LogStorageAccountId       string `tfschema:"log_storage_account_id"`
	TargetDiskEncryptionSetId string `tfschema:"target_disk_encryption_set_id"`
	TargetDiskType            string `tfschema:"target_disk_type"`
}

type NetworkInterfaceModel struct {
	SourceMacAddress string `tfschema:"source_mac_address"`
	TargetStaticIp   string `tfschema:"target_static_ip"`
	TargetSubnetName string `tfschema:"target_subnet_name"`
	TestSubnetName   string `tfschema:"test_subnet_name"`
	IsPrimary        bool   `tfschema:"is_primary"`
}

type SiteRecoveryReplicatedVmVMwareModel struct {
	Name                                  string                  `tfschema:"name"`
	RecoveryVaultId                       string                  `tfschema:"recovery_vault_id"`
	SourceVmName                          string                  `tfschema:"source_vm_name"`
	ApplianceName                         string                  `tfschema:"appliance_name"`
	RecoveryReplicationPolicyId           string                  `tfschema:"recovery_replication_policy_id"`
	PhysicalServerCredentialName          string                  `tfschema:"physical_server_credential_name"`
	LicenseType                           string                  `tfschema:"license_type"`
	TargetResourceGroupId                 string                  `tfschema:"target_resource_group_id"`
	TargetVmName                          string                  `tfschema:"target_vm_name"`
	MultiVmGroupName                      string                  `tfschema:"multi_vm_group_name"`
	TargetProximityPlacementGroupId       string                  `tfschema:"target_proximity_placement_group_id"`
	TargetVmSize                          string                  `tfschema:"target_vm_size"`
	TargetAvailabilitySetId               string                  `tfschema:"target_availability_set_id"`
	TargetZone                            string                  `tfschema:"target_zone"`
	TargetNetworkId                       string                  `tfschema:"target_network_id"`
	TestNetworkId                         string                  `tfschema:"test_network_id"`
	TargetBootDiagnosticsStorageAccountId string                  `tfschema:"target_boot_diagnostics_storage_account_id"`
	DiskToInclude                         []IncludedDiskModel     `tfschema:"managed_disk"`
	NetworkInterface                      []NetworkInterfaceModel `tfschema:"network_interface"`
	DefaultLogStorageAccountId            string                  `tfschema:"default_log_storage_account_id"`
	DefaultRecoveryDiskType               string                  `tfschema:"default_recovery_disk_type"`
	DefaultTargetDiskEncryptionSetId      string                  `tfschema:"default_target_disk_encryption_set_id"`
}

type VMWareReplicatedVmResource struct{}

func (s VMWareReplicatedVmResource) ModelObject() interface{} {
	return &SiteRecoveryReplicatedVmVMwareModel{}
}

func (s VMWareReplicatedVmResource) ResourceType() string {
	return "azurerm_site_recovery_vmware_replicated_vm"
}

func (s VMWareReplicatedVmResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationprotecteditems.ValidateReplicationProtectedItemID
}

var (
	_ sdk.ResourceWithUpdate        = VMWareReplicatedVmResource{}
	_ sdk.ResourceWithCustomizeDiff = VMWareReplicatedVmResource{}
)

func (s VMWareReplicatedVmResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": commonschema.ResourceIDReferenceRequired(&vaults.VaultId{}),

		"source_vm_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"appliance_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_replication_policy_id": commonschema.ResourceIDReferenceRequired(&replicationpolicies.ReplicationPolicyId{}),

		"physical_server_credential_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(replicationprotecteditems.LicenseTypeNotSpecified),
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForLicenseType(), false),
		},

		"target_resource_group_id": commonschema.ResourceIDReferenceRequired(&commonids.ResourceGroupId{}),

		"target_vm_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"default_log_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"managed_disk", "default_log_storage_account_id"},
		},

		"default_recovery_disk_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"managed_disk", "default_recovery_disk_type"},
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
		},

		"default_target_disk_encryption_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateDiskEncryptionSetID,
		},

		"multi_vm_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_proximity_placement_group_id": commonschema.ResourceIDReferenceOptional(&proximityplacementgroups.ProximityPlacementGroupId{}),

		"target_vm_size": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_availability_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAvailabilitySetID,
			ConflictsWith: []string{
				"target_zone",
			},
		},

		"target_zone": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				"target_availability_set_id",
			},
		},

		"target_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
			RequiredWith: []string{"network_interface"},
		},

		"test_network_id": commonschema.ResourceIDReferenceOptional(&commonids.VirtualNetworkId{}),

		"target_boot_diagnostics_storage_account_id": commonschema.ResourceIDReferenceOptional(&commonids.StorageAccountId{}),

		// managed disk is enabled only if mobility service is already installed. (in most cases, it's not installed)
		"managed_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem:     resourceSiteRecoveryVMWareReplicatedVMManagedDiskSchema(),
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem:     resourceSiteRecoveryVMWareReplicatedVMNetworkInterfaceSchema(),
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

func (s VMWareReplicatedVmResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k VMWareReplicatedVmResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// these fields are not returned by API, and only used in creation.
			diff := metadata.ResourceDiff

			_, newStorageAcc := diff.GetChange("default_log_storage_account_id")
			_, newDiskType := diff.GetChange("default_recovery_disk_type")
			_, newDes := diff.GetChange("default_target_disk_encryption_set_id")
			oldDisks, newDisks := diff.GetChange("managed_disk")
			for _, disk := range oldDisks.([]interface{}) {
				disk := disk.(map[string]interface{})
				if newStorageAcc.(string) != "" && disk["log_storage_account_id"] != newStorageAcc.(string) {
					metadata.ResourceDiff.ForceNew("default_log_storage_account_id")
				}
				if newDiskType.(string) != "" && disk["target_disk_type"] != newDiskType.(string) {
					metadata.ResourceDiff.ForceNew("default_recovery_disk_type")
				}
				if newDes.(string) != "" && disk["target_disk_encryption_set_id"] != newDes.(string) {
					metadata.ResourceDiff.ForceNew("default_target_disk_encryption_set_id")
				}
			}

			if diff.HasChanges("managed_disk") {
				// if user has specified `managed_disk`, it forces new.
				// or it acts as an optional field.
				if len(newDisks.([]interface{})) != 0 {
					metadata.ResourceDiff.ForceNew("managed_disk")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (s VMWareReplicatedVmResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing site recovery vmware replicated vm %q: %+v", id, err)
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

			if existing.Model != nil {
				return tf.ImportAsExistsError("azurerm_site_recovery_vmware_replicated_vm", *existing.Model.Id)
			}

			providerSpecificDetail := replicationprotecteditems.InMageRcmEnableProtectionInput{
				LicenseType:              pointer.To(replicationprotecteditems.LicenseType(model.LicenseType)),
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

			poller, err := client.Create(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}
			// once the PUT request returned successfully, an item has been created, even if it may fail in the poll process.
			metadata.SetID(id)

			err = poller.Poller.PollUntilDone(ctx)
			if err != nil {
				return fmt.Errorf("polling %q: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Pending"},
				Target:  []string{"Protected"},
				Refresh: func() (result interface{}, state string, err error) {
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

			err = client.UpdateThenPoll(ctx, id, updateInput)
			if err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			return nil
		},
	}
}

func (s VMWareReplicatedVmResource) Update() sdk.ResourceFunc {
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
					return fmt.Errorf("retrieving `network_interface`: VMNics was nil.")
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
				updateInput.LicenseType = pointer.To(replicationprotecteditems.LicenseType(model.LicenseType))
			} else if existingDetails.LicenseType != nil {
				updateInput.LicenseType = pointer.To(replicationprotecteditems.LicenseType(*existingDetails.LicenseType))
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

			err = client.UpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}

func (s VMWareReplicatedVmResource) Read() sdk.ResourceFunc {
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

func (s VMWareReplicatedVmResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			disableProtectionReason := replicationprotecteditems.DisableProtectionReasonNotSpecified

			disableProtectionInput := replicationprotecteditems.DisableProtectionInput{
				Properties: replicationprotecteditems.DisableProtectionInputProperties{
					DisableProtectionReason: &disableProtectionReason,
					// It's a workaround for https://github.com/hashicorp/pandora/issues/1864
					ReplicationProviderInput: replicationprotecteditems.BaseDisableProtectionProviderSpecificInputImpl{
						InstanceType: string(siterecovery.InstanceTypeDisableProtectionProviderSpecificInput),
					},
				},
			}

			err = client.DeleteThenPoll(ctx, *id, disableProtectionInput)
			if err != nil {
				return fmt.Errorf("deleting %s : %+v", id.String(), err)
			}

			return nil
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

func fetchRunAsAccountsIdBySite(ctx context.Context, runAsAccountClient *vmwarerunasaccounts.RunAsAccountsClient, siteId string, displayName string, applianceName string) (string, error) {
	parsedSiteId, err := vmwarerunasaccounts.ParseVMwareSiteIDInsensitively(siteId)
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

func fetchDiscoveryMachineIdBySite(ctx context.Context, machinesClient *vmwaremachines.MachinesClient, siteId string, machineName string) (string, error) {
	parsedSiteId, err := vmwaremachines.ParseVMwareSiteIDInsensitively(siteId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", siteId, err)
	}

	hackedClient := azuresdkhacks.MachinesClient{Client: machinesClient.Client}
	resp, err := hackedClient.GetAllVMWareMachinesInSite(ctx, *parsedSiteId, vmwaremachines.DefaultGetAllMachinesInSiteOperationOptions())
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

func fetchCredentialByRunAsAccountId(ctx context.Context, client *vmwarerunasaccounts.RunAsAccountsClient, id string) (string, error) {
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

func expandVMWareReplicatedVMNics(input []NetworkInterfaceModel) []replicationprotecteditems.InMageRcmNicInput {
	vmNics := make([]replicationprotecteditems.InMageRcmNicInput, 0)
	for _, nic := range input {
		vmNic := replicationprotecteditems.InMageRcmNicInput{
			NicId:            nic.SourceMacAddress,
			TargetSubnetName: &nic.TargetSubnetName,
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
