// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetrollingupgrades"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
)

type virtualMachineScaleSetUpdateMetaData struct {
	// is "automaticOSUpgrade" enable in the upgradeProfile block
	AutomaticOSUpgradeIsEnabled bool

	// can we reimage instances when `upgrade_mode` is set to `Manual`? this is a feature toggle
	CanReimageOnManualUpgrade bool

	// can we roll instances if we need too? this is a feature toggle
	CanRollInstancesWhenRequired bool

	// do we need to roll the instances in this scale set?
	UpdateInstances bool

	Client   *client.Client
	Existing virtualmachinescalesets.VirtualMachineScaleSet
	ID       *virtualmachinescalesets.VirtualMachineScaleSetId
	OSType   virtualmachinescalesets.OperatingSystemTypes
}

func (metadata virtualMachineScaleSetUpdateMetaData) performUpdate(ctx context.Context, update virtualmachinescalesets.VirtualMachineScaleSetUpdate) error {
	if metadata.AutomaticOSUpgradeIsEnabled {
		// Virtual Machine Scale Sets with Automatic OS Upgrade enabled must have all VM instances upgraded to same
		// Platform Image. Upgrade all VM instances to latest Virtual Machine Scale Set model while property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' is false and then update property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' to true

		update.Properties.UpgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade = pointer.To(false)
	}

	if err := metadata.updateVmss(ctx, update); err != nil {
		return err
	}

	// if we update the SKU, we also need to subsequently roll the instances using the `UpdateInstances` API
	if metadata.UpdateInstances {
		userWantsToRollInstances := metadata.CanRollInstancesWhenRequired
		upgradeMode := metadata.Existing.Properties.UpgradePolicy.Mode

		if userWantsToRollInstances {
			// If the updated image version is not "latest" and upgrade mode is automatic then azure will roll the instances automatically.
			// Calling upgradeInstancesForAutomaticUpgradePolicy() in this case will cause an error.
			if *upgradeMode == virtualmachinescalesets.UpgradeModeAutomatic && isUsingLatestImage(update) {
				if err := metadata.upgradeInstancesForAutomaticUpgradePolicy(ctx); err != nil {
					return err
				}
			}

			if *upgradeMode == virtualmachinescalesets.UpgradeModeManual {
				if err := metadata.upgradeInstancesForManualUpgradePolicy(ctx); err != nil {
					return err
				}
			}
		}
	}

	if metadata.AutomaticOSUpgradeIsEnabled {
		// Virtual Machine Scale Sets with Automatic OS Upgrade enabled must have all VM instances upgraded to same
		// Platform Image. Upgrade all VM instances to latest Virtual Machine Scale Set model while property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' is false and then update property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' to true

		// finally set this to true
		update.Properties.UpgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade = pointer.To(true)

		// then update the VM
		if err := metadata.updateVmss(ctx, update); err != nil {
			return err
		}
	}

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) updateVmss(ctx context.Context, update virtualmachinescalesets.VirtualMachineScaleSetUpdate) error {
	client := metadata.Client.VirtualMachineScaleSetsClient
	id := metadata.ID

	log.Printf("[DEBUG] Updating %s %s", metadata.OSType, id)
	if err := client.UpdateThenPoll(ctx, *id, update, virtualmachinescalesets.DefaultUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s %s: %+v", metadata.OSType, id, err)
	}
	log.Printf("[DEBUG] Updated %s %s", metadata.OSType, id)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForAutomaticUpgradePolicy(ctx context.Context) error {
	rollingUpgradesClient := metadata.Client.VirtualMachineScaleSetRollingUpgradesClient
	id := metadata.ID
	virtualMachineScaleSetId := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName)

	log.Printf("[DEBUG] Updating instances for %s %s", metadata.OSType, id)
	if err := rollingUpgradesClient.StartOSUpgradeThenPoll(ctx, virtualMachineScaleSetId); err != nil {
		return fmt.Errorf("updating instances for %s %s: %+v", metadata.OSType, id, err)
	}
	log.Printf("[DEBUG] Updated instances for %s %s.", metadata.OSType, id)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForManualUpgradePolicy(ctx context.Context) error {
	client := metadata.Client.VirtualMachineScaleSetsClient
	id := metadata.ID

	log.Printf("[DEBUG] Rolling the VM Instances for %s %s", metadata.OSType, id)
	instancesClient := metadata.Client.VirtualMachineScaleSetVMsClient
	virtualMachineScaleSetId := virtualmachinescalesetvms.NewVirtualMachineScaleSetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName)
	instances, err := instancesClient.ListComplete(ctx, virtualMachineScaleSetId, virtualmachinescalesetvms.DefaultListOperationOptions())
	if err != nil {
		return fmt.Errorf("listing VM Instances for %s %s: %+v", metadata.OSType, id, err)
	}

	log.Printf("[DEBUG] Determining instances to roll..")
	instanceIdsToRoll := make([]string, 0)
	for _, item := range instances.Items {
		props := item.Properties
		if props != nil && item.InstanceId != nil {
			latestModel := props.LatestModelApplied
			if !*latestModel {
				instanceIdsToRoll = append(instanceIdsToRoll, *item.InstanceId)
			}
		}
	}

	// TODO: there's a performance enhancement to do batches here, but this is fine for a first pass
	for _, instanceId := range instanceIdsToRoll {
		instanceIds := []string{instanceId}

		log.Printf("[DEBUG] Updating Instance %q to the Latest Configuration..", instanceId)
		ids := virtualmachinescalesets.VirtualMachineScaleSetVMInstanceRequiredIDs{
			InstanceIds: instanceIds,
		}
		if err := client.UpdateInstancesThenPoll(ctx, *id, ids); err != nil {
			return fmt.Errorf("updating Instance %q (%s %s) to the Latest Configuration: %+v", instanceId, metadata.OSType, id, err)
		}
		log.Printf("[DEBUG] Updated Instance %q to the Latest Configuration.", instanceId)

		if metadata.CanReimageOnManualUpgrade {
			log.Printf("[DEBUG] Reimaging Instance %q..", instanceId)
			reImageInput := virtualmachinescalesets.VirtualMachineScaleSetReimageParameters{
				InstanceIds: &instanceIds,
			}
			if err := client.ReimageThenPoll(ctx, *id, reImageInput); err != nil {
				return fmt.Errorf("reimaging Instance %q (%s %s): %+v", instanceId, metadata.OSType, id, err)
			}
			log.Printf("[DEBUG] Reimaged Instance %q..", instanceId)
		}
	}

	log.Printf("[DEBUG] Rolled the VM Instances for %s %s.", metadata.OSType, id)
	return nil
}

func isUsingLatestImage(update virtualmachinescalesets.VirtualMachineScaleSetUpdate) bool {
	if update.Properties.VirtualMachineProfile.StorageProfile == nil ||
		update.Properties.VirtualMachineProfile.StorageProfile.ImageReference == nil ||
		update.Properties.VirtualMachineProfile.StorageProfile.ImageReference.Version == nil {
		return false
	}
	if strings.EqualFold(*update.Properties.VirtualMachineProfile.StorageProfile.ImageReference.Version, "latest") {
		return true
	}
	return false
}
