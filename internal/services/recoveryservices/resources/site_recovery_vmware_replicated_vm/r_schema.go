// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the replicated VM.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": commonschema.ResourceIDReferenceRequired(&vaults.VaultId{}, "Recovery Services Vault"),

		"source_vm_name": {
			Description:  "The name of the source VM in VMWare.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"appliance_name": {
			Description:  "The name of VMWare appliance which handles the replication.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_replication_policy_id": commonschema.ResourceIDReferenceRequired(&replicationpolicies.ReplicationPolicyId{}, "Replication Policy"),

		"physical_server_credential_name": {
			Description:  "The name of the credential to access the source VM. More information about the credentials could be found [here](https://learn.microsoft.com/azure/site-recovery/deploy-vmware-azure-replication-appliance-modernized).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"license_type": {
			Description:  "The license type of the VM.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(replicationprotecteditems.LicenseTypeNotSpecified),
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForLicenseType(), false),
		},

		"target_resource_group_id": commonschema.ResourceIDReferenceRequired(&commonids.ResourceGroupId{}, "Resource Group"),

		"target_vm_name": {
			Description:  "Name of the VM that should be created when a failover is done.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"default_log_storage_account_id": {
			Description:  "The ID of the Storage Account that should be used for logging during replication.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"managed_disk", "default_log_storage_account_id"},
		},

		"default_recovery_disk_type": {
			Description:  "The type of Storage Account that should be used for recovery disks when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
			ExactlyOneOf: []string{"managed_disk", "default_recovery_disk_type"},
		},

		"default_target_disk_encryption_set_id": {
			Description:  "The ID of the default Disk Encryption Set that should be used for the disks when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateDiskEncryptionSetID,
		},

		"multi_vm_group_name": {
			Description:  "Name of group in which all machines will replicate together and have shared crash consistent and app-consistent recovery points when failed over.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_proximity_placement_group_id": commonschema.ResourceIDReferenceOptional(&proximityplacementgroups.ProximityPlacementGroupId{}, "Proximity Placement Group"),

		"target_vm_size": {
			Description:  "Size of the VM that should be created when a failover is done, such as `Standard_F2`. If it's not specified, it will automatically be set by detecting the source VM size.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_availability_set_id": {
			Description:  "The ID of Availability Set that the new VM should belong to when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAvailabilitySetID,
			ConflictsWith: []string{
				"target_zone",
			},
		},

		"target_zone": {
			Description:  "The Availability Zone where the Failover VM should exist.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				"target_availability_set_id",
			},
		},

		"target_network_id": {
			Description:  "The ID of network to use when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
			RequiredWith: []string{"network_interface"},
		},

		"test_network_id": commonschema.ResourceIDReferenceOptional(&commonids.VirtualNetworkId{}, "Virtual Network"),

		"target_boot_diagnostics_storage_account_id": commonschema.ResourceIDReferenceOptional(&commonids.StorageAccountId{}, "Storage Account"),

		// managed disk is enabled only if mobility service is already installed. (in most cases, it's not installed)
		"managed_disk": {
			Description: "It's available only if mobility service is already installed on the source VM.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem:        resourceSiteRecoveryVMWareReplicatedVMManagedDiskSchema(),
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem:     resourceSiteRecoveryVMWareReplicatedVMNetworkInterfaceSchema(),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// these fields are not returned by API, and only used in creation.
			diff := metadata.ResourceDiff

			_, newStorageAcc := diff.GetChange("default_log_storage_account_id")
			_, newDiskType := diff.GetChange("default_recovery_disk_type")
			_, newDes := diff.GetChange("default_target_disk_encryption_set_id")
			oldDisks, newDisks := diff.GetChange("managed_disk")
			for _, disk := range oldDisks.([]any) {
				disk := disk.(map[string]any)
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
				if len(newDisks.([]any)) != 0 {
					metadata.ResourceDiff.ForceNew("managed_disk")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
