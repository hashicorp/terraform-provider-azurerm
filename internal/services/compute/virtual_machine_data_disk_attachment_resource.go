// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineDataDiskAttachment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Read:   resourceVirtualMachineDataDiskAttachmentRead,
		Update: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Delete: resourceVirtualMachineDataDiskAttachmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataDiskID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Compute.VirtualMachinesClient
			id, err := parse.DataDiskID(d.Id())
			if err != nil {
				return nil, err
			}

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(virtualMachine.HttpResponse) {
					return nil, fmt.Errorf("%s was not found therefore Data Disk Attachment cannot be imported", virtualMachineId)
				}

				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var disk *virtualmachines.DataDisk
			if model := virtualMachine.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if *dataDisk.Name == id.Name {
									disk = &dataDisk
									break
								}
							}
						}
					}
				}
			}

			if disk == nil {
				return nil, fmt.Errorf("data disk %s was not found", *id)
			}

			if disk.CreateOption != virtualmachines.DiskCreateOptionTypesAttach && disk.CreateOption != virtualmachines.DiskCreateOptionTypesEmpty {
				return nil, fmt.Errorf("the value of `create_option` for the imported `azurerm_virtual_machine_data_disk_attachment` instance must be `Attach` or `Empty`, whereas now is %s", disk.CreateOption)
			}

			return []*pluginsdk.ResourceData{d}, nil
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
				ValidateFunc:     commonids.ValidateManagedDiskID,
			},

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualMachineID,
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
					string(virtualmachines.CachingTypesNone),
					string(virtualmachines.CachingTypesReadOnly),
					string(virtualmachines.CachingTypesReadWrite),
				}, false),
			},

			"create_option": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(virtualmachines.DiskCreateOptionTypesAttach),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualmachines.DiskCreateOptionTypesAttach),
					string(virtualmachines.DiskCreateOptionTypesEmpty),
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
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedVirtualMachineId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(parsedVirtualMachineId.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(parsedVirtualMachineId.VirtualMachineName, VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, *parsedVirtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("%s was not found", parsedVirtualMachineId)
		}
		return fmt.Errorf("retrieving %s: %+v", parsedVirtualMachineId, err)
	}

	if virtualMachine.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", parsedVirtualMachineId)
	}
	if virtualMachine.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", parsedVirtualMachineId)
	}
	if virtualMachine.Model.Properties.StorageProfile == nil {
		return fmt.Errorf("retrieving %s: `storageprofile` was nil", parsedVirtualMachineId)
	}

	managedDiskId := d.Get("managed_disk_id").(string)
	managedDisk, err := retrieveDataDiskAttachmentManagedDisk(d, meta, managedDiskId)
	if err != nil {
		return fmt.Errorf("retrieving Managed Disk %q: %+v", managedDiskId, err)
	}

	if managedDisk.Sku == nil {
		return fmt.Errorf("unable to determine Storage Account Type for Managed Disk %q: %+v", managedDiskId, err)
	}

	name := *managedDisk.Name
	resourceId := fmt.Sprintf("%s/dataDisks/%s", parsedVirtualMachineId.ID(), name)
	lun := int32(d.Get("lun").(int))
	caching := d.Get("caching").(string)
	createOption := virtualmachines.DiskCreateOptionTypes(d.Get("create_option").(string))
	writeAcceleratorEnabled := d.Get("write_accelerator_enabled").(bool)

	expandedDisk := virtualmachines.DataDisk{
		Name:         pointer.To(name),
		Caching:      pointer.To(virtualmachines.CachingTypes(caching)),
		CreateOption: createOption,
		Lun:          int64(lun),
		ManagedDisk: &virtualmachines.ManagedDiskParameters{
			Id:                 pointer.To(managedDiskId),
			StorageAccountType: pointer.To(virtualmachines.StorageAccountTypes(*managedDisk.Sku.Name)),
		},
		WriteAcceleratorEnabled: pointer.To(writeAcceleratorEnabled),
	}

	// there are ways to provision a VM without a StorageProfile and/or DataDisks
	if virtualMachine.Model.Properties.StorageProfile == nil {
		virtualMachine.Model.Properties.StorageProfile = &virtualmachines.StorageProfile{}
	}

	if virtualMachine.Model.Properties.StorageProfile.DataDisks == nil {
		virtualMachine.Model.Properties.StorageProfile.DataDisks = pointer.To(make([]virtualmachines.DataDisk, 0))
	}

	disks := *virtualMachine.Model.Properties.StorageProfile.DataDisks

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
			return fmt.Errorf("unable to find Disk %q attached to Virtual Machine %q ", name, parsedVirtualMachineId.String())
		}

		disks[existingIndex] = expandedDisk
	}

	virtualMachine.Model.Properties.StorageProfile.DataDisks = &disks

	// fixes #2485
	virtualMachine.Model.Identity = nil
	// fixes #1600
	virtualMachine.Model.Resources = nil
	// fixes #24145
	virtualMachine.Model.Properties.ApplicationProfile = nil

	// if there's too many disks we get a 409 back with:
	//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
	// which we're intentionally not wrapping, since the errors good.
	if err := client.CreateOrUpdateThenPoll(ctx, *parsedVirtualMachineId, *virtualMachine.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s with Disk %q: %+v", parsedVirtualMachineId, name, err)
	}

	d.SetId(resourceId)
	return resourceVirtualMachineDataDiskAttachmentRead(d, meta)
}

func resourceVirtualMachineDataDiskAttachmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)

	virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			log.Printf("[DEBUG] %s was not found therefore Data Disk Attachment cannot exist - removing from state", virtualMachineId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	var disk *virtualmachines.DataDisk
	if model := virtualMachine.Model; model != nil {
		if props := model.Properties; props != nil {
			if profile := props.StorageProfile; profile != nil {
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
		}
	}

	if disk == nil {
		log.Printf("[DEBUG] %s was not found on Virtual Machine %q  - removing from state", id, id.VirtualMachineName)
		d.SetId("")
		return nil
	}

	d.Set("virtual_machine_id", virtualMachineId.ID())
	d.Set("caching", string(pointer.From(disk.Caching)))
	d.Set("create_option", string(disk.CreateOption))
	d.Set("write_accelerator_enabled", disk.WriteAcceleratorEnabled)

	if managedDisk := disk.ManagedDisk; managedDisk != nil {
		d.Set("managed_disk_id", managedDisk.Id)
	}

	d.Set("lun", int(disk.Lun))

	return nil
}

func resourceVirtualMachineDataDiskAttachmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataDiskID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("%s was not found", virtualMachineId)
		}

		return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
	}

	if virtualMachine.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", virtualMachineId)
	}
	if virtualMachine.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", virtualMachineId)
	}
	if virtualMachine.Model.Properties.StorageProfile == nil {
		return fmt.Errorf("retrieving %s: `storageprofile` was nil", virtualMachineId)
	}

	dataDisks := make([]virtualmachines.DataDisk, 0)
	for _, dataDisk := range *virtualMachine.Model.Properties.StorageProfile.DataDisks {
		// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
		if *dataDisk.Name != id.Name {
			dataDisks = append(dataDisks, dataDisk)
		}
	}

	virtualMachine.Model.Properties.StorageProfile.DataDisks = &dataDisks

	// fixes #2485
	virtualMachine.Model.Identity = nil
	// fixes #1600
	virtualMachine.Model.Resources = nil
	// fixes #24145
	virtualMachine.Model.Properties.ApplicationProfile = nil

	if err := client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *virtualMachine.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("removing %s from Virtual Machine %q : %+v", id, id.VirtualMachineName, err)
	}

	return nil
}

func retrieveDataDiskAttachmentManagedDisk(d *pluginsdk.ResourceData, meta interface{}, id string) (*disks.Disk, error) {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedId, err := commonids.ParseManagedDiskID(id)
	if err != nil {
		return nil, fmt.Errorf("parsing Managed Disk ID %q: %+v", id, err)
	}

	resp, err := client.Get(ctx, *parsedId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("managed disk %q was not found", parsedId.String())
		}

		return nil, fmt.Errorf("making Read request on Azure Managed Disk %q : %+v", parsedId.String(), err)
	}

	return resp.Model, nil
}
