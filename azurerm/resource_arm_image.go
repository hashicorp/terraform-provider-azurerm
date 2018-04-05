package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmImageCreateUpdate,
		Read:   resourceArmImageRead,
		Update: resourceArmImageCreateUpdate,
		Delete: resourceArmImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"source_virtual_machine_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"os_disk": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.Linux),
								string(compute.Windows),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"os_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.Generalized),
								string(compute.Specialized),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"managed_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},

						"blob_uri": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(compute.None),
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.CachingTypesNone),
								string(compute.CachingTypesReadOnly),
								string(compute.CachingTypesReadWrite),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
					},
				},
			},

			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"lun": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"managed_disk_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"blob_uri": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(compute.None),
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.CachingTypesNone),
								string(compute.CachingTypesReadOnly),
								string(compute.CachingTypesReadWrite),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmImageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).imageClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Image creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)
	properties := compute.ImageProperties{}

	osDisk, err := expandAzureRmImageOsDisk(d)
	if err != nil {
		return err
	}

	dataDisks, err := expandAzureRmImageDataDisks(d)
	if err != nil {
		return err
	}

	storageProfile := compute.ImageStorageProfile{
		OsDisk:    osDisk,
		DataDisks: &dataDisks,
	}

	sourceVM := compute.SubResource{}
	if v, ok := d.GetOk("source_virtual_machine_id"); ok {
		vmID := v.(string)
		sourceVM = compute.SubResource{
			ID: &vmID,
		}
	}

	//either source VM or storage profile can be specified, but not both
	if sourceVM.ID == nil {
		//if both sourceVM and storageProfile are empty, return an error
		if storageProfile.OsDisk == nil && len(*storageProfile.DataDisks) == 0 {
			return fmt.Errorf("[ERROR] Cannot create image when both source VM and storage profile are empty")
		}

		properties = compute.ImageProperties{
			StorageProfile: &storageProfile,
		}
	} else {
		//creating an image from source VM
		properties = compute.ImageProperties{
			SourceVirtualMachine: &sourceVM,
		}
	}

	createImage := compute.Image{
		Name:            &name,
		Location:        &location,
		Tags:            expandedTags,
		ImageProperties: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, createImage)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read AzureRM Image %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmImageRead(d, meta)
}

func resourceArmImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).imageClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["images"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on AzureRM Image %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	//either source VM or storage profile can be specified, but not both
	if resp.SourceVirtualMachine != nil {
		d.Set("source_virtual_machine_id", resp.SourceVirtualMachine.ID)
	} else if resp.StorageProfile != nil {
		if disk := resp.StorageProfile.OsDisk; disk != nil {
			if err := d.Set("os_disk", flattenAzureRmImageOSDisk(disk)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image OS Disk error: %+v", err)
			}
		}

		if disks := resp.StorageProfile.DataDisks; disks != nil {
			if err := d.Set("data_disk", flattenAzureRmImageDataDisks(disks)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image Data Disks error: %+v", err)
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).imageClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["images"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}

func flattenAzureRmImageOSDisk(osDisk *compute.ImageOSDisk) []interface{} {
	result := make(map[string]interface{})

	if disk := osDisk; disk != nil {
		if uri := osDisk.BlobURI; uri != nil {
			result["blob_uri"] = *uri
		}
		if diskSizeDB := osDisk.DiskSizeGB; diskSizeDB != nil {
			result["size_gb"] = *diskSizeDB
		}
		if disk := osDisk.ManagedDisk; disk != nil {
			result["managed_disk_id"] = *disk.ID
		}
		result["caching"] = string(osDisk.Caching)
		result["os_type"] = osDisk.OsType
		result["os_state"] = osDisk.OsState
	}

	return []interface{}{result}
}

func flattenAzureRmImageDataDisks(diskImages *[]compute.ImageDataDisk) []interface{} {
	result := make([]interface{}, 0)

	if images := diskImages; images != nil {
		for _, disk := range *images {
			l := make(map[string]interface{})
			if disk.BlobURI != nil {
				l["blob_uri"] = *disk.BlobURI
			}
			l["caching"] = string(disk.Caching)
			if disk.DiskSizeGB != nil {
				l["size_gb"] = *disk.DiskSizeGB
			}
			l["lun"] = *disk.Lun
			if disk.ManagedDisk != nil {
				l["managed_disk_id"] = *disk.ManagedDisk.ID
			}

			result = append(result, l)
		}
	}

	return result
}

func expandAzureRmImageOsDisk(d *schema.ResourceData) (*compute.ImageOSDisk, error) {
	osDisk := &compute.ImageOSDisk{}
	disks := d.Get("os_disk").(*schema.Set).List()

	if len(disks) > 0 {
		config := disks[0].(map[string]interface{})

		if v := config["os_type"].(string); v != "" {
			osType := compute.OperatingSystemTypes(v)
			osDisk.OsType = osType
		}

		if v := config["os_state"].(string); v != "" {
			osState := compute.OperatingSystemStateTypes(v)
			osDisk.OsState = osState
		}

		managedDiskID := config["managed_disk_id"].(string)
		if managedDiskID != "" {
			managedDisk := &compute.SubResource{
				ID: &managedDiskID,
			}
			osDisk.ManagedDisk = managedDisk
		}

		blobURI := config["blob_uri"].(string)
		osDisk.BlobURI = &blobURI

		if v := config["caching"].(string); v != "" {
			caching := compute.CachingTypes(v)
			osDisk.Caching = caching
		}

		if size := config["size_gb"]; size != 0 {
			diskSize := int32(size.(int))
			osDisk.DiskSizeGB = &diskSize
		}
	}

	return osDisk, nil
}

func expandAzureRmImageDataDisks(d *schema.ResourceData) ([]compute.ImageDataDisk, error) {

	disks := d.Get("data_disk").([]interface{})

	dataDisks := make([]compute.ImageDataDisk, 0, len(disks))
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]interface{})

		managedDiskID := d.Get("managed_disk_id").(string)
		blobURI := d.Get("blob_uri").(string)
		lun := int32(config["lun"].(int))

		dataDisk := compute.ImageDataDisk{
			Lun:     &lun,
			BlobURI: &blobURI,
		}

		if size := d.Get("size_gb"); size != 0 {
			diskSize := int32(size.(int))
			dataDisk.DiskSizeGB = &diskSize
		}

		if v := d.Get("caching").(string); v != "" {
			caching := compute.CachingTypes(v)
			dataDisk.Caching = caching
		}

		if managedDiskID != "" {
			managedDisk := &compute.SubResource{
				ID: &managedDiskID,
			}
			dataDisk.ManagedDisk = managedDisk
		}

		dataDisks = append(dataDisks, dataDisk)
	}

	return dataDisks, nil

}
