package compute

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
	ID       *parse.VirtualMachineScaleSetId
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
			if upgradeMode == compute.Automatic {
				if err := metadata.upgradeInstancesForAutomaticUpgradePolicy(ctx); err != nil {
					return err
				}
			}

			if upgradeMode == compute.Manual {
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

	log.Printf("[DEBUG] Updating %s Virtual Machine Scale Set %q (Resource Group %q)..", metadata.OSType, id.Name, id.ResourceGroup)
	future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
	if err != nil {
		return fmt.Errorf("Error updating L%sinux Virtual Machine Scale Set %q (Resource Group %q): %+v", metadata.OSType, id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for update of %s Virtual Machine Scale Set %q (Resource Group %q)..", metadata.OSType, id.Name, id.ResourceGroup)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of %s Virtual Machine Scale Set %q (Resource Group %q): %+v", metadata.OSType, id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Updated %s Virtual Machine Scale Set %q (Resource Group %q).", metadata.OSType, id.Name, id.ResourceGroup)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForAutomaticUpgradePolicy(ctx context.Context) error {
	client := metadata.Client.VMScaleSetClient
	rollingUpgradesClient := metadata.Client.VMScaleSetRollingUpgradesClient
	id := metadata.ID

	log.Printf("[DEBUG] Updating instances for %s Virtual Machine Scale Set %q (Resource Group %q)..", metadata.OSType, id.Name, id.ResourceGroup)
	future, err := rollingUpgradesClient.StartOSUpgrade(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error updating instances for %s Virtual Machine Scale Set %q (Resource Group %q): %+v", metadata.OSType, id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for update of instances for %s Virtual Machine Scale Set %q (Resource Group %q)..", metadata.OSType, id.Name, id.ResourceGroup)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of instances for %s Virtual Machine Scale Set %q (Resource Group %q): %+v", metadata.OSType, id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Updated instances for %s Virtual Machine Scale Set %q (Resource Group %q).", metadata.OSType, id.Name, id.ResourceGroup)

	return nil
}

func (metadata virtualMachineScaleSetUpdateMetaData) upgradeInstancesForManualUpgradePolicy(ctx context.Context) error {
	client := metadata.Client.VMScaleSetClient
	id := metadata.ID

	log.Printf("[DEBUG] Rolling the VM Instances for %s Virtual Machine Scale Set %q (Resource Group %q)..", metadata.OSType, id.Name, id.ResourceGroup)
	instancesClient := metadata.Client.VMScaleSetVMsClient
	instances, err := instancesClient.ListComplete(ctx, id.ResourceGroup, id.Name, "", "", "")
	if err != nil {
		return fmt.Errorf("Error listing VM Instances for %s Virtual Machine Scale Set %q (Resource Group %q): %+v", metadata.OSType, id.Name, id.ResourceGroup, err)
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
			return fmt.Errorf("Error enumerating instances: %s", err)
		}
	}

	// TODO: there's a performance enhancement to do batches here, but this is fine for a first pass
	for _, instanceId := range instanceIdsToRoll {
		instanceIds := []string{instanceId}

		log.Printf("[DEBUG] Updating Instance %q to the Latest Configuration..", instanceId)
		ids := compute.VirtualMachineScaleSetVMInstanceRequiredIDs{
			InstanceIds: &instanceIds,
		}
		future, err := client.UpdateInstances(ctx, id.ResourceGroup, id.Name, ids)
		if err != nil {
			return fmt.Errorf("Error updating Instance %q (%s VM Scale Set %q / Resource Group %q) to the Latest Configuration: %+v", instanceId, metadata.OSType, id.Name, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Instance %q (%s VM Scale Set %q / Resource Group %q) to the Latest Configuration: %+v", instanceId, metadata.OSType, id.Name, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated Instance %q to the Latest Configuration.", instanceId)

		// TODO: does this want to be a separate, user-configurable toggle?
		log.Printf("[DEBUG] Reimaging Instance %q..", instanceId)
		reimageInput := &compute.VirtualMachineScaleSetReimageParameters{
			InstanceIds: &instanceIds,
		}
		reimageFuture, err := client.Reimage(ctx, id.ResourceGroup, id.Name, reimageInput)
		if err != nil {
			return fmt.Errorf("Error reimaging Instance %q (%s VM Scale Set %q / Resource Group %q): %+v", instanceId, metadata.OSType, id.Name, id.ResourceGroup, err)
		}

		if err = reimageFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for reimage of Instance %q (%s VM Scale Set %q / Resource Group %q): %+v", instanceId, metadata.OSType, id.Name, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Reimaged Instance %q..", instanceId)
	}

	log.Printf("[DEBUG] Rolled the VM Instances for %s Virtual Machine Scale Set %q (Resource Group %q).", metadata.OSType, id.Name, id.ResourceGroup)
	return nil
}
