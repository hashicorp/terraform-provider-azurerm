package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachineDataDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineDataDiskAttachmentCreateUpdate,
		Read:   resourceArmVirtualMachineDataDiskAttachmentRead,
		Update: resourceArmVirtualMachineDataDiskAttachmentCreateUpdate,
		Delete: resourceArmVirtualMachineDataDiskAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"managed_disk_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lun": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"create_option": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(compute.DiskCreateOptionTypesAttach),
				ValidateFunc:     validation.StringInSlice([]string{}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"caching": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.CachingTypesNone),
					string(compute.CachingTypesReadOnly),
					string(compute.CachingTypesReadWrite),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
		},
	}
}

func resourceArmVirtualMachineDataDiskAttachmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	virtualMachineId := d.Get("virtual_machine_id").(string)
	parsedVirtualMachineId, err := parseAzureResourceID(virtualMachineId)
	if err != nil {
		return fmt.Errorf("Error parsing Virtual Machine ID %q: %+v", virtualMachineId, err)
	}

	resourceGroup := parsedVirtualMachineId.ResourceGroup
	virtualMachineName := parsedVirtualMachineId.Path["virtualMachines"]

	azureRMLockByName(virtualMachineName, virtualMachineResourceName)
	defer azureRMUnlockByName(virtualMachineName, virtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, resourceGroup, virtualMachineName, compute.InstanceView)
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			return fmt.Errorf("Virtual Machine %q (Resource Group %q) was not found", virtualMachineName, resourceGroup)
		}

		return fmt.Errorf("Error loading Virtual Machine %q (Resource Group %q): %+v", virtualMachineName, resourceGroup, err)
	}

	managedDiskId := d.Get("managed_disk_id").(string)
	managedDisk, err := retrieveDataDiskAttachmentManagedDisk(meta, managedDiskId)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Disk %q: %+v", managedDiskId, err)
	}

	if managedDisk.Sku == nil {
		return fmt.Errorf("Error: unable to determine Storage Account Type for Managed Disk %q: %+v", managedDiskId, err)
	}

	name := *managedDisk.Name
	lun := int32(d.Get("lun").(int))
	caching := d.Get("caching").(string)
	createOption := compute.DiskCreateOptionTypes(d.Get("create_option").(string))

	expandedDisk := compute.DataDisk{
		Name:         utils.String(name),
		Caching:      compute.CachingTypes(caching),
		CreateOption: createOption,
		Lun:          utils.Int32(lun),
		ManagedDisk: &compute.ManagedDiskParameters{
			ID:                 utils.String(managedDiskId),
			StorageAccountType: managedDisk.Sku.Name,
		},
	}

	disks := *virtualMachine.StorageProfile.DataDisks
	if d.IsNewResource() {
		disks = append(disks, expandedDisk)
	} else {
		// iterate over the disks and swap it out in-place
		existingIndex := -1
		for i, disk := range disks {
			if *disk.Name == name {
				existingIndex = i
				break
			}
		}

		if existingIndex == -1 {
			return fmt.Errorf("Unable to find Disk %q attached to Virtual Machine %q (Resource Group %q)", name, virtualMachineName, resourceGroup)
		}

		disks[existingIndex] = expandedDisk
	}

	virtualMachine.StorageProfile.DataDisks = &disks

	// if there's too many disks we get a 409 back with:
	//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualMachineName, virtualMachine)
	if err != nil {
		return fmt.Errorf("Error updating Virtual Machine %q (Resource Group %q) with Disk %q: %+v", virtualMachineName, resourceGroup, name, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for Virtual Machine %q (Resource Group %q) to finish updating Disk %q: %+v", virtualMachineName, resourceGroup, name, err)
	}

	d.SetId(fmt.Sprintf("%s/dataDisks/%s", virtualMachineId, name))

	return resourceArmVirtualMachineDataDiskAttachmentRead(d, meta)
}

func resourceArmVirtualMachineDataDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	virtualMachineName := id.Path["virtualMachines"]
	name := id.Path["dataDisks"]

	virtualMachine, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			return fmt.Errorf("Virtual Machine %q (Resource Group %q) was not found", virtualMachineName, resourceGroup)
		}

		return fmt.Errorf("Error loading Virtual Machine %q (Resource Group %q): %+v", virtualMachineName, resourceGroup, err)
	}

	var disk *compute.DataDisk
	if profile := virtualMachine.StorageProfile; profile != nil {
		if dataDisks := profile.DataDisks; dataDisks != nil {
			for _, dataDisk := range *dataDisks {
				// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.Equals`
				if *dataDisk.Name == name {
					disk = &dataDisk
					break
				}
			}
		}
	}

	if disk == nil {
		log.Printf("[DEBUG] Data Disk %q was not found on Virtual Machine %q (Resource Group %q) - removing from state", name, virtualMachineName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("virtual_machine_id", virtualMachine.ID)
	d.Set("caching", string(disk.Caching))
	d.Set("create_option", string(disk.CreateOption))

	if managedDisk := disk.ManagedDisk; managedDisk != nil {
		d.Set("managed_disk_id", managedDisk.ID)
	}

	if lun := disk.Lun; lun != nil {
		d.Set("lun", int(*lun))
	}

	return nil
}

func resourceArmVirtualMachineDataDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	virtualMachineName := id.Path["virtualMachines"]
	name := id.Path["dataDisks"]

	azureRMLockByName(virtualMachineName, virtualMachineResourceName)
	defer azureRMUnlockByName(virtualMachineName, virtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			return fmt.Errorf("Virtual Machine %q (Resource Group %q) was not found", virtualMachineName, resourceGroup)
		}

		return fmt.Errorf("Error loading Virtual Machine %q (Resource Group %q): %+v", virtualMachineName, resourceGroup, err)
	}

	dataDisks := make([]compute.DataDisk, 0)
	for _, dataDisk := range *virtualMachine.StorageProfile.DataDisks {
		// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.Equals`
		if *dataDisk.Name != name {
			dataDisks = append(dataDisks, dataDisk)
		}
	}

	virtualMachine.StorageProfile.DataDisks = &dataDisks

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualMachineName, virtualMachine)
	if err != nil {
		return fmt.Errorf("Error removing Disk %q from Virtual Machine %q (Resource Group %q): %+v", name, virtualMachineName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for Disk %q to be removed from Virtual Machine %q (Resource Group %q): %+v", name, virtualMachineName, resourceGroup, err)
	}

	return nil
}

func retrieveDataDiskAttachmentManagedDisk(meta interface{}, id string) (*compute.Disk, error) {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	parsedId, err := parseAzureResourceID(id)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Managed Disk ID %q: %+v", id, err)
	}
	resourceGroup := parsedId.ResourceGroup
	name := parsedId.Path["disks"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("Error Managed Disk %q (Resource Group %q) was not found!", name, resourceGroup)
		}

		return nil, fmt.Errorf("Error making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return &resp, nil
}
