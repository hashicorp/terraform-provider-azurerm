package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedDiskCreate,
		Read:   resourceArmManagedDiskRead,
		Update: resourceArmManagedDiskCreate,
		Delete: resourceArmManagedDiskDelete,
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

			"storage_account_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disk.PremiumLRS),
					string(disk.StandardLRS),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"create_option": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disk.Copy),
					string(disk.Empty),
					string(disk.FromImage),
					string(disk.Import),
				}, true),
			},

			"source_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_reference_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(disk.Windows),
					string(disk.Linux),
				}, true),
			},

			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateDiskSizeGB,
			},

			"tags": tagsSchema(),
		},
	}
}

func validateDiskSizeGB(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 4095 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` can only be between 1 and 4095"))
	}
	return
}

func resourceArmManagedDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	diskClient := client.diskClient

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	createDisk := disk.Model{
		Name:     &name,
		Location: &location,
		Tags:     expandedTags,
	}

	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)

	createDisk.Properties = &disk.Properties{
		AccountType: disk.StorageAccountTypes(storageAccountType),
		OsType:      disk.OperatingSystemTypes(osType),
	}

	if v := d.Get("disk_size_gb"); v != 0 {
		diskSize := int32(v.(int))
		createDisk.Properties.DiskSizeGB = &diskSize
	}

	createOption := d.Get("create_option").(string)
	createDisk.CreationData = &disk.CreationData{
		CreateOption: disk.CreateOption(createOption),
	}

	if strings.EqualFold(createOption, string(disk.Import)) {
		if sourceUri := d.Get("source_uri").(string); sourceUri != "" {
			createDisk.CreationData.SourceURI = &sourceUri
		} else {
			return fmt.Errorf("[ERROR] source_uri must be specified when create_option is `%s`", disk.Import)
		}
	} else if strings.EqualFold(createOption, string(disk.Copy)) {
		if sourceResourceId := d.Get("source_resource_id").(string); sourceResourceId != "" {
			createDisk.CreationData.SourceResourceID = &sourceResourceId
		} else {
			return fmt.Errorf("[ERROR] source_resource_id must be specified when create_option is `%s`", disk.Copy)
		}
	} else if strings.EqualFold(createOption, string(disk.FromImage)) {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			createDisk.CreationData.ImageReference = &disk.ImageDiskReference{
				ID: utils.String(imageReferenceId),
			}
		} else {
			return fmt.Errorf("[ERROR] image_reference_id must be specified when create_option is `%s`", disk.FromImage)
		}
	}

	_, diskErr := diskClient.CreateOrUpdate(resGroup, name, createDisk, make(chan struct{}))
	err := <-diskErr
	if err != nil {
		return err
	}

	read, err := diskClient.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read Managed Disk %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmManagedDiskRead(d, meta)
}

func resourceArmManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	diskClient := meta.(*ArmClient).diskClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	resp, err := diskClient.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %s (resource group %s): %s", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if resp.Properties != nil {
		flattenAzureRmManagedDiskProperties(d, resp.Properties)
	}

	if resp.CreationData != nil {
		flattenAzureRmManagedDiskCreationData(d, resp.CreationData)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmManagedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	diskClient := meta.(*ArmClient).diskClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	deleteResp, deleteErr := diskClient.Delete(resGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-deleteErr
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return err
		}
	}

	return nil
}

func flattenAzureRmManagedDiskProperties(d *schema.ResourceData, properties *disk.Properties) {
	d.Set("storage_account_type", string(properties.AccountType))
	if properties.DiskSizeGB != nil {
		d.Set("disk_size_gb", *properties.DiskSizeGB)
	}
	if properties.OsType != "" {
		d.Set("os_type", string(properties.OsType))
	}
}

func flattenAzureRmManagedDiskCreationData(d *schema.ResourceData, creationData *disk.CreationData) {
	d.Set("create_option", string(creationData.CreateOption))
	if ref := creationData.ImageReference; ref != nil {
		d.Set("image_reference_id", *ref.ID)
	}
	if id := creationData.SourceResourceID; id != nil {
		d.Set("source_resource_id", *id)
	}
	if creationData.SourceURI != nil {
		d.Set("source_uri", *creationData.SourceURI)
	}
}
