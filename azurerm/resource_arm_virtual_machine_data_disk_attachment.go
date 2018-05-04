package azurerm

import (
	"fmt"

	"log"
	"strings"

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"create_option": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"lun": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"vhd_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"managed_disk_id", "managed_disk_type"},
			},

			"managed_disk_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ConflictsWith:    []string{"vhd_uri"},
			},

			"managed_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.PremiumLRS),
					string(compute.StandardLRS),
				}, true),
				ConflictsWith: []string{"vhd_uri"},
			},

			"caching": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateDiskSizeGB,
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
		return err
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

	name := d.Get("name").(string)
	createOption := d.Get("create_option").(string)
	lun := int32(d.Get("lun").(int))

	expandedDisk := compute.DataDisk{
		Name:         utils.String(name),
		Lun:          utils.Int32(lun),
		CreateOption: compute.DiskCreateOptionTypes(createOption),
	}
	if v, ok := d.GetOk("vhd_uri"); ok {
		expandedDisk.Vhd = &compute.VirtualHardDisk{
			URI: utils.String(v.(string)),
		}
	} else {
		storageAccountType := d.Get("managed_disk_type").(string)
		expandedDisk.ManagedDisk = &compute.ManagedDiskParameters{
			StorageAccountType: compute.StorageAccountTypes(storageAccountType),
		}

		if v, ok := d.GetOk("managed_disk_id"); ok {
			expandedDisk.ManagedDisk.ID = utils.String(v.(string))
		}
	}

	if v, ok := d.GetOk("caching"); ok {
		expandedDisk.Caching = compute.CachingTypes(v.(string))
	}

	if v, ok := d.GetOk("disk_size_gb"); ok {
		expandedDisk.DiskSizeGB = utils.Int32(int32(v.(int)))
	}

	disks := *virtualMachine.StorageProfile.DataDisks
	if !d.IsNewResource() {
		// find the existing disk, remove it from the array
		dataDisks := make([]compute.DataDisk, 0)
		for _, dataDisk := range disks {
			if !strings.EqualFold(*dataDisk.Name, *expandedDisk.Name) {
				disks = append(disks, dataDisk)
			}
		}
		disks = dataDisks
	}

	disks = append(disks, expandedDisk)
	virtualMachine.StorageProfile.DataDisks = &disks

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualMachineName, virtualMachine)
	if err != nil {
		return fmt.Errorf("Error updating Virtual Machine %q (Resource Group %q) with Disk %q: %+v", virtualMachineName, resourceGroup, *expandedDisk.Name, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for Virtual Machine %q (Resource Group %q) to finish updating Disk %q: %+v", virtualMachineName, resourceGroup, *expandedDisk.Name, err)
	}

	d.SetId(fmt.Sprintf("%s/dataDisks/%s", virtualMachineId, *expandedDisk.Name))

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
				if strings.EqualFold(*dataDisk.Name, name) {
					disk = &dataDisk
					break
				}
			}
		}
	}

	if disk == nil {
		log.Printf("[DEBUG] Disk %q was not found on Virtual Machine %q (Resource Group %q) - removing from state", name, virtualMachineName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("virtual_machine_id", virtualMachine.ID)

	if vhd := disk.Vhd; vhd != nil {
		d.Set("vhd_uri", vhd.URI)
	}

	if managedDisk := disk.ManagedDisk; managedDisk != nil {
		d.Set("managed_disk_id", managedDisk.ID)
		d.Set("managed_disk_type", string(managedDisk.StorageAccountType))
	}

	d.Set("create_option", string(disk.CreateOption))
	d.Set("caching", string(disk.Caching))
	if diskSizeGb := disk.DiskSizeGB; diskSizeGb != nil {
		d.Set("disk_size_gb", int(*diskSizeGb))
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
		if !strings.EqualFold(*dataDisk.Name, name) {
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
