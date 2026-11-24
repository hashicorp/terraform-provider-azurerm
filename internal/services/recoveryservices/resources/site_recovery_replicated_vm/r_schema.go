// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func siteRecoveryReplicatedVmSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the replication for the replicated VM.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Description:  "The name of the vault that should be updated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_recovery_fabric_name": {
			Description:  "Name of fabric that should contain this replication.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_vm_id": {
			Description:  "ID of the VM to replicate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
			// user-specified segments are lower cased too.
			// tracked on https://github.com/Azure/azure-rest-api-specs/issues/24393
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_recovery_fabric_id": {
			Description:  "ID of fabric where the VM replication should be handled when a failover is done.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"recovery_replication_policy_id": {
			Description:  "ID of the policy to use for this replicated VM.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},

		"source_recovery_protection_container_name": {
			Description:  "Name of the protection container to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_recovery_protection_container_id": {
			Description:  "ID of protection container where the VM replication should be created when a failover is done.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationprotectioncontainers.ValidateReplicationProtectionContainerID,
		},

		"target_resource_group_id": {
			Description:  "ID of Resource Group where the VM should be created when a failover is done.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateResourceGroupID,
		},

		"target_availability_set_id": {
			Description:  "ID of Availability Set that the new VM should belong to when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAvailabilitySetID,
			ConflictsWith: []string{
				"target_zone",
			},
		},

		"target_zone": commonschema.ZoneSingleOptionalForceNew(),

		"target_network_id": {
			Description:  "Network to use when a failover is done (recommended to set if any network_interface is configured for failover).",
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"test_network_id": {
			Description:  "Network to use when a test failover is done.",
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"target_edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"unmanaged_disk": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			ForceNew:   true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_uri": {
						Description:  "ID of disk that should be replicated.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"staging_storage_account_id": {
						Description:  "Storage Account that should be used for caching.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"target_storage_account_id": {
						Description:  "Storage Account disk should belong to when a failover is done.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
		},

		"multi_vm_group_name": {
			Description:  "Name of group in which all machines will replicate together and have shared crash consistent and app-consistent recovery points when failed over.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"managed_disk": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Set:        resourceSiteRecoveryReplicatedVMDiskHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_id": {
						Description:  "ID of disk that should be replicated.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"staging_storage_account_id": {
						Description:  "Storage Account that should be used for caching.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"target_resource_group_id": {
						Description:  "Resource Group disk should belong to when a failover is done.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateResourceGroupID,
					},

					"target_disk_type": {
						Description:  "What type should the disk be when a failover is done.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
					},

					"target_replica_disk_type": {
						Description:  "What type should the disk be that holds the replication data.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
					},

					"target_disk_encryption_set_id": {
						Description:  "The Disk Encryption Set that the Managed Disk will be associated with.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					},

					"target_disk_encryption": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						MaxItems:   1,
						Elem:       diskEncryptionResource(),
					},
				},
			},
		},

		"target_proximity_placement_group_id": {
			Description:  "ID of Proximity Placement Group the new VM should belong to when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"target_boot_diagnostic_storage_account_id": {
			Description:  "ID of the Storage Account which the new VM should used for boot diagnostic when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"target_capacity_reservation_group_id": {
			Description:  "ID of the Capacity Reservation Group where the new VM should belong to when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
		},

		"target_virtual_machine_scale_set_id": {
			Description:  "ID of the Virtual Machine Scale Set which the new VM should belong to when a failover is done.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
		},

		"target_virtual_machine_size": {
			Description: "The size the Virtual Machine should have.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C if not specified, this gets set to the vm_size of the virtual machine
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_interface": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem:       networkInterfaceResource(),
		},
	}
}

func resourceSiteRecoveryReplicatedVMCustomizeDiff(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
	rawConfig := diff.GetRawConfig()
	if !rawConfig.IsKnown() || rawConfig.IsNull() {
		return nil
	}

	networkInterfaces := rawConfig.GetAttr("network_interface")
	if networkInterfaces.IsKnown() && !networkInterfaces.IsNull() {
		for _, networkInterface := range networkInterfaces.AsValueSlice() {
			if !networkInterface.IsKnown() || networkInterface.IsNull() {
				return nil
			}

			ipConfigurations := networkInterface.GetAttr("ip_configuration")
			if !ipConfigurations.IsKnown() || ipConfigurations.IsNull() || ipConfigurations.LengthInt() <= 1 {
				continue
			}

			primaryCount := 0
			for _, ipConfiguration := range ipConfigurations.AsValueSlice() {
				if !ipConfiguration.IsKnown() || ipConfiguration.IsNull() {
					return nil
				}

				name := ipConfiguration.GetAttr("name")
				if !name.IsKnown() {
					return nil
				}
				if name.IsNull() || name.AsString() == "" {
					return fmt.Errorf("each `network_interface.ip_configuration` block must specify `name` when multiple blocks are configured")
				}

				primary := ipConfiguration.GetAttr("primary")
				if !primary.IsKnown() {
					return nil
				}
				if !primary.IsNull() && primary.True() {
					primaryCount++
				}
			}

			if primaryCount != 1 {
				return fmt.Errorf("`network_interface.ip_configuration` must contain exactly one block with `primary` set to `true` when multiple blocks are configured")
			}
		}
	}

	if diff.Id() != "" {
		managedDiskCustomDiff := pluginsdk.CustomDiffComputedSetCannotRemove(
			"managed_disk",
			resourceSiteRecoveryReplicatedVMDiskHash,
			[]string{"staging_storage_account_id", "target_resource_group_id", "target_disk_encryption_set_id"},
			[]string{"target_disk_type", "target_replica_disk_type"},
			[]string{"staging_storage_account_id", "target_resource_group_id", "target_disk_encryption_set_id"},
		)
		if err := managedDiskCustomDiff(ctx, diff, v); err != nil {
			return err
		}
	}

	return nil
}
