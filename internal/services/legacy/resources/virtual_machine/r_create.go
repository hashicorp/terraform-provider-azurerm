// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachines.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_machine", id.ID())
			}
		}
	}

	osDisk, err := expandAzureRmVirtualMachineOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile := virtualmachines.StorageProfile{
		OsDisk: osDisk,
	}

	if _, ok := d.GetOk("storage_image_reference"); ok {
		imageRef, err2 := expandAzureRmVirtualMachineImageReference(d)
		if err2 != nil {
			return err2
		}
		storageProfile.ImageReference = imageRef
	}

	if _, ok := d.GetOk("storage_data_disk"); ok {
		dataDisks, err2 := expandAzureRmVirtualMachineDataDisk(d)
		if err2 != nil {
			return err2
		}
		storageProfile.DataDisks = &dataDisks
	}

	networkProfile := expandAzureRmVirtualMachineNetworkProfile(d)
	vmSize := d.Get("vm_size").(string)
	properties := virtualmachines.VirtualMachineProperties{
		NetworkProfile: &networkProfile,
		HardwareProfile: &virtualmachines.HardwareProfile{
			VMSize: pointer.ToEnum[virtualmachines.VirtualMachineSizeTypes](vmSize),
		},
		StorageProfile: &storageProfile,
	}

	if v, ok := d.GetOk("license_type"); ok {
		properties.LicenseType = pointer.To(v.(string))
	}

	if _, ok := d.GetOk("boot_diagnostics"); ok {
		diagnosticsProfile := expandAzureRmVirtualMachineDiagnosticsProfile(d)
		if diagnosticsProfile != nil {
			properties.DiagnosticsProfile = diagnosticsProfile
		}
	}
	if _, ok := d.GetOk("additional_capabilities"); ok {
		properties.AdditionalCapabilities = expandAzureRmVirtualMachineAdditionalCapabilities(d)
	}

	if _, ok := d.GetOk("os_profile"); ok {
		osProfile, err2 := expandAzureRmVirtualMachineOsProfile(d)
		if err2 != nil {
			return err2
		}
		properties.OsProfile = osProfile
	}

	if v, ok := d.GetOk("availability_set_id"); ok {
		availSet := virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}

		properties.AvailabilitySet = &availSet
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		properties.ProximityPlacementGroup = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	vm := virtualmachines.VirtualMachine{
		Name:       &id.VirtualMachineName,
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &properties,
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
		Zones:      helpers.ExpandZones(d.Get("zones").([]any)),
	}

	if _, ok := d.GetOk("identity"); ok {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		vm.Identity = identityExpanded
	}

	if _, ok := d.GetOk("plan"); ok {
		vm.Plan = expandAzureRmVirtualMachinePlan(d)
	}

	locks.ByName(id.VirtualMachineName, compute.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, compute.VirtualMachineResourceName)

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, vm, virtualmachines.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, vm, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}
	read, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return err
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if read.Model.Id == nil {
		return fmt.Errorf("retrieving %s: `id` was nil", id)
	}

	ipAddress, err := determineVirtualMachineIPAddress(ctx, meta, read.Model.Properties)
	if err != nil {
		return fmt.Errorf("determining IP Address for %s: %+v", id, err)
	}

	provisionerType := "ssh"
	if props := read.Model.Properties; props != nil {
		if profile := props.OsProfile; profile != nil {
			if profile.WindowsConfiguration != nil {
				provisionerType = "winrm"
			}
		}
	}
	d.SetConnInfo(map[string]string{
		"type": provisionerType,
		"host": ipAddress,
	})

	return resourceVirtualMachineRead(d, meta)
}
