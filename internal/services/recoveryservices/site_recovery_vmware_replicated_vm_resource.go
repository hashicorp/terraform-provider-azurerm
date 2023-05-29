package recoveryservices

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	vmwaremachines "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
	vmwarerunasaccounts "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
	validateResourceGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type installAccountType string

const lincreds installAccountType = "lincreds"
const v2arcmlab installAccountType = "v2arcmlab"
const wincreds installAccountType = "wincreds"

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
	CredentialType                        string                  `tfschema:"credential_type"`
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

var _ sdk.ResourceWithUpdate = VMWareReplicatedVmResource{}
var _ sdk.ResourceWithCustomizeDiff = VMWareReplicatedVmResource{}

func (s VMWareReplicatedVmResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationprotecteditems.ValidateVaultID,
		},

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

		"recovery_replication_policy_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     azure.ValidateResourceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"credential_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(lincreds),
				string(v2arcmlab),
				string(wincreds),
			}, false),
		},

		"license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(replicationprotecteditems.LicenseTypeNotSpecified),
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForLicenseType(), false),
		},

		"target_resource_group_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validateResourceGroup.ResourceGroupID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_vm_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"default_log_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageaccounts.ValidateStorageAccountID,
			AtLeastOneOf: []string{"managed_disk", "default_log_storage_account_id"},
		},

		"default_recovery_disk_type": { // only works in creation and will not be returned from service
			Type:         pluginsdk.TypeString,
			Optional:     true,
			AtLeastOneOf: []string{"managed_disk", "default_recovery_disk_type"},
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
		},

		"default_target_disk_encryption_set_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"managed_disk"},
			ValidateFunc:  diskencryptionsets.ValidateDiskEncryptionSetID,
		},

		"multi_vm_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_proximity_placement_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"target_vm_size": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_availability_set_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     availabilitysets.ValidateAvailabilitySetID,
			DiffSuppressFunc: suppress.CaseDifference,
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
		},

		"target_subnet_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"test_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"test_subnet_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_boot_diagnostics_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageaccounts.ValidateStorageAccountID,
		},

		// managed disk is enabled only if mobility service is already installed. (in most cases, it's not installed)
		"managed_disk": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem:     resourceSiteRecoveryVMWareReplicatedVMManagedDiskSchema(),
		},

		"network_interface": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem:     resourceSiteRecoveryVMWareReplicatedVMNetworkInterfaceSchema(),
		},
	}
}

func resourceSiteRecoveryVMWareReplicatedVMManagedDiskSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"disk_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsNotEmpty,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"target_disk_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
			},

			"target_disk_encryption_set_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     diskencryptionsets.ValidateDiskEncryptionSetID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"log_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceSiteRecoveryVMWareReplicatedVMNetworkInterfaceSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_mac_address": { // if it was left blank, we can use the only one NIC id as the source mac address
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"target_static_ip": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
			_, newDisks := diff.GetChange("managed_disk")
			for _, disk := range newDisks.(*schema.Set).List() {
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
				// if user has specified `managed_disk`, it's force new.
				// or it acts as an optional field.
				if len(newDisks.(*schema.Set).List()) != 1 {
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
		// waiting for fully protected cost very long time.
		Timeout: 300 * time.Minute,
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

			if model.CredentialType == "" {
				return fmt.Errorf("`credential_type` must be specified in creation")
			}

			vaultId, err := replicationprotecteditems.ParseVaultID(model.RecoveryVaultId)
			if err != nil {
				return fmt.Errorf("parsing vault id %q: %+v", model.RecoveryVaultId, err)
			}

			containerId, err := fetchSiteRecoveryReplicatedVmVMWareContainerId(ctx, containerClient, vaultId.ID())
			if err != nil {
				return fmt.Errorf("fetch Replication Container from vault %q: %+v", vaultId, err)
			}

			parsedContainerId, err := replicationprotecteditems.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parse %s: %+v", containerId, err)
			}

			id := replicationprotecteditems.NewReplicationProtectedItemID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, model.Name)
			fabricId := replicationfabrics.NewReplicationFabricID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName)

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

			runAsAccountId, err := fetchRunAsAccountsIdBySite(ctx, runAsAccountsClient, siteID, model.CredentialType, model.ApplianceName)
			if err != nil {
				return fmt.Errorf("fetch run as account id %s: %+v", model.CredentialType, err)
			}

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing site recovery vmware replicated vm %q: %+v", id, err)
				}
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

			if model.TargetVmSize != "" {
				providerSpecificDetail.TargetVMSize = &model.TargetVmSize
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

			if model.TargetAvailabilitySetId != "" {
				providerSpecificDetail.TargetAvailabilitySetId = pointer.To(model.TargetAvailabilitySetId)
			}

			if model.TargetZone != "" {
				providerSpecificDetail.TargetAvailabilityZone = pointer.To(model.TargetZone)
			}

			if model.TargetBootDiagnosticsStorageAccountId != "" {
				providerSpecificDetail.TargetBootDiagnosticsStorageAccountId = pointer.To(model.TargetBootDiagnosticsStorageAccountId)
			}

			if model.MultiVmGroupName != "" {
				providerSpecificDetail.MultiVMGroupName = &model.MultiVmGroupName
			}

			if model.TargetProximityPlacementGroupId != "" {
				providerSpecificDetail.TargetProximityPlacementGroupId = &model.TargetProximityPlacementGroupId
			}

			if model.ApplianceName != "" {
				providerSpecificDetail.ProcessServerId = processServerId
			}

			// split test network and subnet into isolated parameters
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
			metadata.SetID(id) // once the PUT request returned successfully, an item has been created, even if it may fail in the poll process.

			err = poller.Poller.PollUntilDone()
			if err != nil {
				return fmt.Errorf("polling %q: %+v", id, err)
			}

			err = resourceSiteRecoveryReplicatedVmVMWareClassicUpdateInternal(ctx, metadata)
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
			return resourceSiteRecoveryReplicatedVmVMWareClassicUpdateInternal(ctx, metadata)
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
					ReplicationProviderInput: &siterecovery.DisableProtectionProviderSpecificInput{
						InstanceType: siterecovery.InstanceTypeDisableProtectionProviderSpecificInput,
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

func resourceSiteRecoveryReplicatedVmVMWareClassicUpdateInternal(ctx context.Context, metadata sdk.ResourceMetaData) error {
	client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient
	containerClient := metadata.Client.RecoveryServices.ProtectionContainerClient
	// We are only allowed to update the configuration once the VM is fully protected
	state, err := waitForVmwareReplicationToBeHealthy(ctx, metadata)
	if err != nil {
		return err
	}

	var model SiteRecoveryReplicatedVmVMwareModel
	err = metadata.Decode(&model)
	if err != nil {
		return fmt.Errorf("decoding %+v", err)
	}

	containerId, err := fetchSiteRecoveryReplicatedVmVMWareContainerId(ctx, containerClient, model.RecoveryVaultId)
	if err != nil {
		return fmt.Errorf("fetch Replication Container from vault %q: %+v", model.RecoveryVaultId, err)
	}
	parsedContainerId, err := replicationprotectioncontainers.ParseReplicationProtectionContainerID(containerId)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", containerId, err)
	}

	id := replicationprotecteditems.NewReplicationProtectedItemID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, model.Name)

	existing, err := client.Get(ctx, id)
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

	var targetAvailabilitySetID *string
	if model.TargetAvailabilitySetId != "" {
		targetAvailabilitySetID = &model.TargetAvailabilitySetId
	} else {
		targetAvailabilitySetID = nil
	}

	var targetAvailabilityZone *string
	if model.TargetZone != "" {
		targetAvailabilityZone = &model.TargetZone
	} else {
		targetAvailabilityZone = nil
	}

	vmNics := make([]replicationprotecteditems.InMageRcmNicInput, 0)
	if metadata.ResourceData.HasChange("network_interface") {
		for _, nic := range model.NetworkInterface {
			vmNic := replicationprotecteditems.InMageRcmNicInput{
				TargetSubnetName: &nic.TargetSubnetName,
			}
			if nic.SourceMacAddress != "" {
				vmNic.NicId = nic.SourceMacAddress
			} else {
				if len(model.NetworkInterface) > 1 {
					return fmt.Errorf("when `source_mac_address` is not set, there must be exactly one `network_interface` block")
				}
				if state == nil || state.Properties == nil || state.Properties.ProviderSpecificDetails == nil {
					return fmt.Errorf("failed to get nic id from state")
				}

				if detail, ok := state.Properties.ProviderSpecificDetails.(replicationprotecteditems.InMageRcmReplicationDetails); ok {
					if detail.VMNics != nil && len(*detail.VMNics) == 1 && (*detail.VMNics)[0].NicId != nil {
						vmNic.NicId = *(*detail.VMNics)[0].NicId
					} else {
						return fmt.Errorf("when `source_mac_address` is not set, there must be exactly one Network Adapter on the source VM")
					}
				} else {
					return fmt.Errorf("unexpected provider specific details type: %T", state.Properties.ProviderSpecificDetails)
				}
			}

			if nic.TargetStaticIp != "" {
				vmNic.TargetStaticIPAddress = &nic.TargetStaticIp
			}
			if nic.IsPrimary {
				vmNic.IsPrimaryNic = strconv.FormatBool(true)
				vmNic.IsSelectedForFailover = utils.String("true")
			} else {
				vmNic.IsPrimaryNic = strconv.FormatBool(false)
				vmNic.IsSelectedForFailover = utils.String("false")
			}
			vmNics = append(vmNics, vmNic)
		}

		if model.TargetNetworkId == "" && len(vmNics) > 0 {
			return fmt.Errorf("`target_network_id` must be set when a `network_interface` is configured")
		}
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
	} else {
		if existingDetails.LicenseType != nil {
			updateInput.LicenseType = pointer.To(replicationprotecteditems.LicenseType(*existingDetails.LicenseType))
		}
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
		updateInput.TargetAvailabilitySetId = targetAvailabilitySetID
	} else {
		updateInput.TargetAvailabilitySetId = existingDetails.TargetAvailabilitySetId
	}

	if metadata.ResourceData.HasChange("target_zone") {
		updateInput.TargetAvailabilityZone = targetAvailabilityZone
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

	if metadata.ResourceData.HasChange("target_availability_set_id") {
		props.RecoveryAvailabilitySetId = targetAvailabilitySetID
	} else {
		props.RecoveryAvailabilitySetId = existingDetails.TargetAvailabilitySetId
	}

	if metadata.ResourceData.HasChange("target_vm_size") {
		props.RecoveryAzureVMSize = &model.TargetVmSize
	} else {
		props.RecoveryAzureVMSize = existingDetails.TargetVMSize
	}

	parameters := replicationprotecteditems.UpdateReplicationProtectedItemInput{
		Properties: &props,
	}

	err = client.UpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	return nil
}

func waitForVmwareReplicationToBeHealthy(ctx context.Context, metadata sdk.ResourceMetaData) (*replicationprotecteditems.ReplicationProtectedItem, error) {
	stateConf := &pluginsdk.StateChangeConf{
		Target:       []string{"Protected", "normal"},
		Refresh:      waitForVmwareClassicReplicationToBeHealthyRefreshFunc(ctx, metadata),
		PollInterval: time.Minute,
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for site recovery to replicate vm: %+v", err)
	}

	protectedItem, ok := result.(replicationprotecteditems.ReplicationProtectedItem)
	if ok {
		return &protectedItem, nil
	} else {
		return nil, fmt.Errorf("waiting for site recovery return incompatible type")
	}
}

func waitForVmwareClassicReplicationToBeHealthyRefreshFunc(ctx context.Context, metadata sdk.ResourceMetaData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
		if err != nil {
			return nil, "", err
		}

		client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, "", fmt.Errorf("making Read request on site recovery replicated vm Vmware Classic %s : %+v", id.String(), err)
		}

		if resp.Model == nil {
			return nil, "", fmt.Errorf("Missing Model in response when making Read request on site recovery replicated vm Vmware Classic %s  %+v", id.String(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("Missing Properties in response when making Read request on site recovery replicated vm Vmware Classic %s  %+v", id.String(), err)
		}

		if resp.Model.Properties.ProviderSpecificDetails == nil {
			return nil, "", fmt.Errorf("missing Properties.ProviderSpecificDetails in response when making Read request on site recovery replicated vm Vmware Classic %s : %+v", id.String(), err)
		}

		// Find first disk that is not fully replicated yet
		if resp.Model.Properties.ProtectionState == nil {
			return nil, "", fmt.Errorf("missing ProtectionState in response when making Read request on site recovery replicated vm Vmware Classic %s : %+v", id.String(), err)
		}
		return *resp.Model, *resp.Model.Properties.ProtectionState, nil
	}
}

func fetchSiteRecoveryReplicatedVmVMWareContainerId(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, vaultId string) (string, error) {
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

func fetchRunAsAccountsIdBySite(ctx context.Context, runAsAccountClient *vmwarerunasaccounts.RunAsAccountsClient, siteId string, accountType string, applianceName string) (string, error) {
	parsedSiteId, err := vmwarerunasaccounts.ParseVMwareSiteIDInsensitively(siteId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", siteId, err)
	}

	hackedClinet := azuresdkhacks.RunAsAccountsClient{Client: runAsAccountClient.Client}
	resp, err := hackedClinet.GetAllRunAsAccountsInSiteComplete(ctx, *parsedSiteId)
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
		if strings.EqualFold(*account.Properties.DisplayName, accountType) && strings.EqualFold(*account.Properties.ApplianceName, applianceName) {
			return *account.Id, nil
		}
	}

	return "", fmt.Errorf("retiring %q: run as account %s not found", siteId, accountType)
}

func fetchProcessServerIdByName(ctx context.Context, fabricClient *replicationfabrics.ReplicationFabricsClient, fabricId replicationfabrics.ReplicationFabricId, processServerName string) (string, error) {
	resp, err := fabricClient.Get(ctx, fabricId, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		return "", err
	}

	if resp.Model == nil {
		return "", fmt.Errorf("retiring %q: Model is nil", fabricId)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retiring %q: Properties is nil", fabricId)
	}

	if resp.Model.Properties.CustomDetails == nil {
		return "", fmt.Errorf("retiring %q: CustomDetails is nil", fabricId)
	}

	if detail, ok := resp.Model.Properties.CustomDetails.(replicationfabrics.InMageRcmFabricSpecificDetails); ok {
		if detail.ProcessServers == nil || len(*detail.ProcessServers) < 1 {
			return "", fmt.Errorf("retiring %q: count of Process Servers is 0", fabricId)
		}
		for _, server := range *detail.ProcessServers {
			if strings.EqualFold(*server.Name, processServerName) {
				return *server.Id, nil
			}
		}
		return "", fmt.Errorf("retiring %q: process server %s not found", fabricId, processServerName)
	}
	return "", fmt.Errorf("retiring %q: Detail Type mismatch", fabricId)
}

func fetchDiscoveryMachineIdBySite(ctx context.Context, machinesClient *vmwaremachines.MachinesClient, siteId string, machineName string) (string, error) {
	parsedSiteId, err := vmwaremachines.ParseVMwareSiteIDInsensitively(siteId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", siteId, err)
	}

	resp, err := getAllVMWareMachinesInSite(ctx, machinesClient, *parsedSiteId, vmwaremachines.DefaultGetAllMachinesInSiteOperationOptions())
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
				return handleAzureSdkForGoBug2824(*machine.Id), nil
			}
		}
	}

	return "", fmt.Errorf("retiring %q: machine %s not found", siteId, machineName)
}

// workaround for https://github.com/hashicorp/go-azure-sdk/issues/492
func getAllVMWareMachinesInSite(ctx context.Context, c *vmwaremachines.MachinesClient, id vmwaremachines.VMwareSiteId, options vmwaremachines.GetAllMachinesInSiteOperationOptions) (result vmwaremachines.GetAllMachinesInSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/machines", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	return warpedExecutePaged(ctx, req)
}

func warpedExecutePaged(ctx context.Context, req *client.Request) (result vmwaremachines.GetAllMachinesInSiteOperationResponse, err error) {
	resp, err := req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values   *[]vmwaremachines.VMwareMachine `json:"value"`
		NextLink *string                         `json:"nextLink"`
	}

	if err = resp.Unmarshal(&values); err != nil {
		return
	}
	result.Model = values.Values

	if values.NextLink != nil {
		nextReq := req
		u, err := url.Parse(*values.NextLink)
		if err != nil {
			return result, err
		}
		nextReq.URL = u
		nextResp, err := warpedExecutePaged(ctx, nextReq)
		if err != nil {
			return result, err
		}
		if nextResp.Model != nil {
			result.Model = pointer.To(append(*result.Model, *nextResp.Model...))
		}
	}

	return
}

func fetchVmwareSiteIdByFabric(ctx context.Context, fabricClient *replicationfabrics.ReplicationFabricsClient, fabricId replicationfabrics.ReplicationFabricId) (string, error) {
	resp, err := fabricClient.Get(ctx, fabricId, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		return "", err
	}
	if resp.Model == nil {
		return "", fmt.Errorf("retiring %q: Model is nil", fabricId)
	}
	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retiring %q: Properties is nil", fabricId)
	}
	if v, ok := resp.Model.Properties.CustomDetails.(replicationfabrics.InMageRcmFabricSpecificDetails); ok {
		if v.VMwareSiteId == nil {
			return "", fmt.Errorf("retiring %q: VMwareSiteId is nil", fabricId)
		}
		return *v.VMwareSiteId, nil
	}
	return "", fmt.Errorf("retiring %q: Detail Type mismatch", fabricId)
}
