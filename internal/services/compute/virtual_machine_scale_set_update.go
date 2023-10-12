// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

type virtualMachineScaleSetUpdateMetaData struct {
	// is "automaticOSUpgrade" enable in the upgradeProfile block
	AutomaticOSUpgradeIsEnabled bool

	// can we roll instances if we need too? this is a feature toggle
	CanRollInstancesWhenRequired bool

	// do we need to roll the instances in this scale set?
	UpdateInstances bool

	Client   *client.Client
	Existing compute.VirtualMachineScaleSet
	ID       *commonids.VirtualMachineScaleSetId
	OSType   compute.OperatingSystemTypes
}

func (metadata virtualMachineScaleSetUpdateMetaData) performUpdate(ctx context.Context, update compute.VirtualMachineScaleSetUpdate) error {
	if metadata.AutomaticOSUpgradeIsEnabled {
		// Virtual Machine Scale Sets with Automatic OS Upgrade enabled must have all VM instances upgraded to same
		// Platform Image. Upgrade all VM instances to latest Virtual Machine Scale Set model while property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' is false and then update property
		// 'upgradePolicy.automaticOSUpgradePolicy.enableAutomaticOSUpgrade' to true

		update.VirtualMachineScaleSetUpdateProperties.UpgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade = utils.Bool(false)
	}

	if err := metadata.updateVmss(ctx, update); err != nil {
		return err
	}

	// if we update the SKU, we also need to subsequently roll the instances using the `UpdateInstances` API
	if metadata.UpdateInstances {
		userWantsToRollInstances := metadata.CanRollInstancesWhenRequired
		upgradeMode := metadata.Existing.VirtualMachineScaleSetProperties.UpgradePolicy.Mode

		if userWantsToRollInstances {
			// If the updated image version is not "latest" and upgrade mode is automatic then azure will roll the instances automatically.
			// Calling upgradeInstancesForAutomaticUpgradePolicy() in this case will cause an error.
			if upgradeMode == compute.UpgradeModeAutomatic && isUsingLatestImage(update) {
				if err := metadata.upgradeInstancesForAutomaticUpgradePolicy(ctx); err != nil {
					return err
				}
			}

			if upgradeMode == compute.UpgradeModeManual {
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
		update.VirtualMachineScaleSetUpdateProperties.UpgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade = utils.Bool(true)

		// then update the VM
		if err := metadata.updateVmss(ctx, update); err != nil {
			return err
		}
	}

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) updateVmss(ctx context.Context, update compute.VirtualMachineScaleSetUpdate) error {
	client := metadata.Client.VMScaleSetClient
	id := metadata.ID

	log.Printf("[DEBUG] Updating %s %s", metadata.OSType, id)
	future, err := client.Update(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, update)
	if err != nil {
		return fmt.Errorf("updating %s %s: %+v", metadata.OSType, id, err)
	}

	log.Printf("[DEBUG] Waiting for update of %s %s", metadata.OSType, id)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s %s: %+v", metadata.OSType, id, err)
	}
	log.Printf("[DEBUG] Updated %s %s", metadata.OSType, id)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForAutomaticUpgradePolicy(ctx context.Context) error {
	client := metadata.Client.VMScaleSetClient
	rollingUpgradesClient := metadata.Client.VMScaleSetRollingUpgradesClient
	id := metadata.ID

	log.Printf("[DEBUG] Updating instances for %s %s", metadata.OSType, id)
	future, err := rollingUpgradesClient.StartOSUpgrade(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName)
	if err != nil {
		return fmt.Errorf("updating instances for %s %s: %+v", metadata.OSType, id, err)
	}

	log.Printf("[DEBUG] Waiting for update of instances for %s %s", metadata.OSType, id)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of instances for %s %s: %+v", metadata.OSType, id, err)
	}
	log.Printf("[DEBUG] Updated instances for %s %s.", metadata.OSType, id)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForManualUpgradePolicy(ctx context.Context) error {
	client := metadata.Client.VMScaleSetClient
	id := metadata.ID

	log.Printf("[DEBUG] Rolling the VM Instances for %s %s", metadata.OSType, id)
	instancesClient := metadata.Client.VMScaleSetVMsClient
	instances, err := instancesClient.ListComplete(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, "", "", "")
	if err != nil {
		return fmt.Errorf("listing VM Instances for %s %s: %+v", metadata.OSType, id, err)
	}

	log.Printf("[DEBUG] Determining instances to roll..")
	instanceIdsToRoll := make([]string, 0)
	for instances.NotDone() {
		instance := instances.Value()
		props := instance.VirtualMachineScaleSetVMProperties
		if props != nil && instance.InstanceID != nil {
			latestModel := props.LatestModelApplied
			if latestModel != nil || !*latestModel {
				instanceIdsToRoll = append(instanceIdsToRoll, *instance.InstanceID)
			}
		}

		if err := instances.NextWithContext(ctx); err != nil {
			return fmt.Errorf("enumerating instances: %s", err)
		}
	}

	// TODO: there's a performance enhancement to do batches here, but this is fine for a first pass
	for _, instanceId := range instanceIdsToRoll {
		instanceIds := []string{instanceId}

		log.Printf("[DEBUG] Updating Instance %q to the Latest Configuration..", instanceId)
		ids := compute.VirtualMachineScaleSetVMInstanceRequiredIDs{
			InstanceIds: &instanceIds,
		}
		future, err := client.UpdateInstances(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, ids)
		if err != nil {
			return fmt.Errorf("updating Instance %q (%s %s) to the Latest Configuration: %+v", instanceId, metadata.OSType, id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of Instance %q (%s %s) to the Latest Configuration: %+v", instanceId, metadata.OSType, id, err)
		}
		log.Printf("[DEBUG] Updated Instance %q to the Latest Configuration.", instanceId)

		// TODO: does this want to be a separate, user-configurable toggle?
		log.Printf("[DEBUG] Reimaging Instance %q..", instanceId)
		reimageInput := &compute.VirtualMachineScaleSetReimageParameters{
			InstanceIds: &instanceIds,
		}
		reimageFuture, err := client.Reimage(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, reimageInput)
		if err != nil {
			return fmt.Errorf("reimaging Instance %q (%s %s): %+v", instanceId, metadata.OSType, id, err)
		}

		if err = reimageFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for reimage of Instance %q (%s %s): %+v", instanceId, metadata.OSType, id, err)
		}
		log.Printf("[DEBUG] Reimaged Instance %q..", instanceId)
	}

	log.Printf("[DEBUG] Rolled the VM Instances for %s %s.", metadata.OSType, id)
	return nil
}

func isUsingLatestImage(update compute.VirtualMachineScaleSetUpdate) bool {
	if update.VirtualMachineProfile.StorageProfile == nil ||
		update.VirtualMachineProfile.StorageProfile.ImageReference == nil ||
		update.VirtualMachineProfile.StorageProfile.ImageReference.Version == nil {
		return false
	}
	if strings.EqualFold(*update.VirtualMachineProfile.StorageProfile.ImageReference.Version, "latest") {
		return true
	}
	return false
}
