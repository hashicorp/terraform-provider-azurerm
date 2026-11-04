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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_recovery_fabric_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_vm_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
			// user-specified segments are lower cased too.
			// tracked on https://github.com/Azure/azure-rest-api-specs/issues/24393
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_recovery_fabric_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"recovery_replication_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},

		"source_recovery_protection_container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_recovery_protection_container_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationprotectioncontainers.ValidateReplicationProtectionContainerID,
		},

		"target_resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateResourceGroupID,
		},

		"target_availability_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAvailabilitySetID,
			ConflictsWith: []string{
				"target_zone",
			},
		},

		"target_zone": commonschema.ZoneSingleOptionalForceNew(),

		"target_network_id": {
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"test_network_id": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"staging_storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"target_storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
		},

		"multi_vm_group_name": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"staging_storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"target_resource_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateResourceGroupID,
					},

					"target_disk_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
					},

					"target_replica_disk_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(replicationprotecteditems.PossibleValuesForDiskAccountType(), false),
					},

					"target_disk_encryption_set_id": {
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"target_boot_diagnostic_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"target_capacity_reservation_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
		},

		"target_virtual_machine_scale_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
		},

		"target_virtual_machine_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
