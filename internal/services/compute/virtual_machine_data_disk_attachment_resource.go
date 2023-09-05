// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceVirtualMachineDataDiskAttachment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Read:   resourceVirtualMachineDataDiskAttachmentRead,
		Update: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Delete: resourceVirtualMachineDataDiskAttachmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataDiskID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     disks.ValidateDiskID,
			},

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualMachineID,
			},

			"lun": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"caching": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.CachingTypesNone),
					string(compute.CachingTypesReadOnly),
					string(compute.CachingTypesReadWrite),
				}, false),
			},

			"create_option": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(compute.DiskCreateOptionTypesAttach),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.DiskCreateOptionTypesAttach),
					string(compute.DiskCreateOptionTypesEmpty),
				}, false),
			},

			"write_accelerator_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceVirtualMachineDataDiskAttachmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedVirtualMachineId, err := parse.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Virtual Machine ID %q: %+v", parsedVirtualMachineId.ID(), err)
	}

	locks.ByName(parsedVirtualMachineId.Name, VirtualMachineResourceName)
	defer locks.UnlockByName(parsedVirtualMachineId.Name, VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, parsedVirtualMachineId.ResourceGroup, parsedVirtualMachineId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			return fmt.Errorf("Virtual Machine %q  was not found", parsedVirtualMachineId.String())
		}

		return fmt.Errorf("loading Virtual Machine %q : %+v", parsedVirtualMachineId.String(), err)
	}

	managedDiskId := d.Get("managed_disk_id").(string)
	managedDisk, err := retrieveDataDiskAttachmentManagedDisk(d, meta, managedDiskId)
	if err != nil {
		return fmt.Errorf("retrieving Managed Disk %q: %+v", managedDiskId, err)
	}

	if managedDisk.Sku == nil {
		return fmt.Errorf("Error: unable to determine Storage Account Type for Managed Disk %q: %+v", managedDiskId, err)
	}

	name := *managedDisk.Name
	resourceId := fmt.Sprintf("%s/dataDisks/%s", parsedVirtualMachineId.ID(), name)
	lun := int32(d.Get("lun").(int))
	caching := d.Get("caching").(string)
	createOption := compute.DiskCreateOptionTypes(d.Get("create_option").(string))
	writeAcceleratorEnabled := d.Get("write_accelerator_enabled").(bool)

	expandedDisk := compute.DataDisk{
		Name:         utils.String(name),
		Caching:      compute.CachingTypes(caching),
		CreateOption: createOption,
		Lun:          utils.Int32(lun),
		ManagedDisk: &compute.ManagedDiskParameters{
			ID:                 utils.String(managedDiskId),
			StorageAccountType: compute.StorageAccountTypes(string(*managedDisk.Sku.Name)),
		},
		WriteAcceleratorEnabled: utils.Bool(writeAcceleratorEnabled),
	}

	disks := *virtualMachine.StorageProfile.DataDisks

	existingIndex := -1
	for i, disk := range disks {
		if *disk.Name == name {
			existingIndex = i
			break
		}
	}

	if d.IsNewResource() {
		if existingIndex != -1 {
			return tf.ImportAsExistsError("azurerm_virtual_machine_data_disk_attachment", resourceId)
		}

		disks = append(disks, expandedDisk)
	} else {
		if existingIndex == -1 {
			return fmt.Errorf("Unable to find Disk %q attached to Virtual Machine %q ", name, parsedVirtualMachineId.String())
		}

		disks[existingIndex] = expandedDisk
	}

	virtualMachine.StorageProfile.DataDisks = &disks

	// fixes #2485
	virtualMachine.Identity = nil
	// fixes #1600
	virtualMachine.Resources = nil

	// if there's too many disks we get a 409 back with:
	//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
	// which we're intentionally not wrapping, since the errors good.
	future, err := client.CreateOrUpdate(ctx, parsedVirtualMachineId.ResourceGroup, parsedVirtualMachineId.Name, virtualMachine)
	if err != nil {
		return fmt.Errorf("updating Virtual Machine %q  with Disk %q: %+v", parsedVirtualMachineId.String(), name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Virtual Machine %q to finish updating Disk %q: %+v", parsedVirtualMachineId.String(), name, err)
	}

	d.SetId(resourceId)
	return resourceVirtualMachineDataDiskAttachmentRead(d, meta)
}

func resourceVirtualMachineDataDiskAttachmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	virtualMachine, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			log.Printf("[DEBUG] Virtual Machine %q was not found (Resource Group %q) therefore Data Disk Attachment cannot exist - removing from state", id.VirtualMachineName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading Virtual Machine %q : %+v", id.String(), err)
	}

	var disk *compute.DataDisk
	if profile := virtualMachine.StorageProfile; profile != nil {
		if dataDisks := profile.DataDisks; dataDisks != nil {
			for _, dataDisk := range *dataDisks {
				// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
				if *dataDisk.Name == id.Name {
					disk = &dataDisk
					break
				}
			}
		}
	}

	if disk == nil {
		log.Printf("[DEBUG] Data Disk %q was not found on Virtual Machine %q  - removing from state", id.Name, id.String())
		d.SetId("")
		return nil
	}

	d.Set("virtual_machine_id", virtualMachine.ID)
	d.Set("caching", string(disk.Caching))
	d.Set("create_option", string(disk.CreateOption))
	d.Set("write_accelerator_enabled", disk.WriteAcceleratorEnabled)

	if managedDisk := disk.ManagedDisk; managedDisk != nil {
		d.Set("managed_disk_id", managedDisk.ID)
	}

	if lun := disk.Lun; lun != nil {
		d.Set("lun", int(*lun))
	}

	return nil
}

func resourceVirtualMachineDataDiskAttachmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			return fmt.Errorf("Virtual Machine %q was not found", id.String())
		}

		return fmt.Errorf("loading Virtual Machine %q : %+v", id.String(), err)
	}

	dataDisks := make([]compute.DataDisk, 0)
	for _, dataDisk := range *virtualMachine.StorageProfile.DataDisks {
		// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
		if *dataDisk.Name != id.Name {
			dataDisks = append(dataDisks, dataDisk)
		}
	}

	virtualMachine.StorageProfile.DataDisks = &dataDisks

	// fixes #2485
	virtualMachine.Identity = nil
	// fixes #1600
	virtualMachine.Resources = nil

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualMachineName, virtualMachine)
	if err != nil {
		return fmt.Errorf("removing Disk %q from Virtual Machine %q : %+v", id.Name, id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Disk %q to be removed from Virtual Machine %q : %+v", id.Name, id.String(), err)
	}

	return nil
}

func retrieveDataDiskAttachmentManagedDisk(d *pluginsdk.ResourceData, meta interface{}, id string) (*disks.Disk, error) {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedId, err := disks.ParseDiskID(id)
	if err != nil {
		return nil, fmt.Errorf("parsing Managed Disk ID %q: %+v", id, err)
	}

	resp, err := client.Get(ctx, *parsedId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("Managed Disk %q  was not found!", parsedId.String())
		}

		return nil, fmt.Errorf("making Read request on Azure Managed Disk %q : %+v", parsedId.String(), err)
	}

	return resp.Model, nil
}
