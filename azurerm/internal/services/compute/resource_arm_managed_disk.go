package compute

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedDiskCreateUpdate,
		Read:   resourceArmManagedDiskRead,
		Update: resourceArmManagedDiskCreateUpdate,
		Delete: resourceArmManagedDiskDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaSingleZone(),

			"storage_account_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.StandardLRS),
					string(compute.PremiumLRS),
					string(compute.StandardSSDLRS),
					string(compute.UltraSSDLRS),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"create_option": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Empty),
					string(compute.FromImage),
					string(compute.Import),
					string(compute.Restore),
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
					string(compute.Windows),
					string(compute.Linux),
				}, true),
			},

			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateDiskSizeGB,
			},

			"disk_iops_read_write": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.EncryptionAtRestWithPlatformKey),
					string(compute.EncryptionAtRestWithCustomerKey),
				}, false),
				Default: string(compute.EncryptionAtRestWithPlatformKey),
			},

			"managed_disk_encryption_set_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"tags": tags.Schema(),
		},
	}
}

func validateDiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 32767 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` can only be between 0 and 32767"))
	}
	return warnings, errors
}

func resourceArmManagedDiskCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Managed Disk %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_disk", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	createDisk := compute.Disk{
		Name:     &name,
		Location: &location,
		DiskProperties: &compute.DiskProperties{
			OsType: compute.OperatingSystemTypes(osType),
		},
		Sku: &compute.DiskSku{
			Name: compute.DiskStorageAccountTypes(storageAccountType),
		},
		Tags:  expandedTags,
		Zones: zones,
	}

	if v := d.Get("disk_size_gb"); v != 0 {
		diskSize := int32(v.(int))
		createDisk.DiskProperties.DiskSizeGB = &diskSize
	}

	if strings.EqualFold(storageAccountType, string(compute.UltraSSDLRS)) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			createDisk.DiskProperties.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int32(v.(int))
			createDisk.DiskProperties.DiskMBpsReadWrite = &diskMBps
		}
	} else {
		if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") {
			return fmt.Errorf("[ERROR] disk_iops_read_write and disk_mbps_read_write are only available for UltraSSD disks")
		}
	}

	createOption := d.Get("create_option").(string)
	createDisk.CreationData = &compute.CreationData{
		CreateOption: compute.DiskCreateOption(createOption),
	}

	if strings.EqualFold(createOption, string(compute.Import)) {
		if sourceUri := d.Get("source_uri").(string); sourceUri != "" {
			createDisk.CreationData.SourceURI = &sourceUri
		} else {
			return fmt.Errorf("[ERROR] source_uri must be specified when create_option is `%s`", compute.Import)
		}
	} else if strings.EqualFold(createOption, string(compute.Copy)) || strings.EqualFold(createOption, string(compute.Restore)) {
		if sourceResourceId := d.Get("source_resource_id").(string); sourceResourceId != "" {
			createDisk.CreationData.SourceResourceID = &sourceResourceId
		} else {
			return fmt.Errorf("[ERROR] source_resource_id must be specified when create_option is `%s`", compute.Copy)
		}
	} else if strings.EqualFold(createOption, string(compute.FromImage)) {
		if imageReferenceId := d.Get("image_reference_id").(string); imageReferenceId != "" {
			createDisk.CreationData.ImageReference = &compute.ImageDiskReference{
				ID: utils.String(imageReferenceId),
			}
		} else {
			return fmt.Errorf("[ERROR] image_reference_id must be specified when create_option is `%s`", compute.FromImage)
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		createDisk.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	encryption := compute.Encryption{}

	if v, ok := d.GetOk("encryption_type"); ok {
		encryption.Type = compute.EncryptionType(v.(string))
		if strings.EqualFold(v.(string), string(compute.EncryptionAtRestWithPlatformKey)) {
			if _, ok := d.GetOk("managed_disk_encryption_set_id"); ok {
				return fmt.Errorf("[Error] `managed_disk_encryption_set_id` should not be set when `encryption_type` is `%s`", compute.EncryptionAtRestWithPlatformKey)
			}
		} else if strings.EqualFold(v.(string), string(compute.EncryptionAtRestWithCustomerKey)) {
			if v, ok := d.GetOk("managed_disk_encryption_set_id"); ok {
				encryption.DiskEncryptionSetID = utils.String(v.(string))
			} else {
				return fmt.Errorf("[Error] `managed_disk_encryption_set_id` must be set when `encryption_type` is `%s`", compute.EncryptionAtRestWithCustomerKey)
			}
		}
	}

	createDisk.Encryption = &encryption

	future, err := client.CreateOrUpdate(ctx, resGroup, name, createDisk)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
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
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %s (resource group %s): %s", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			flattenAzureRmManagedDiskCreationData(d, creationData)
		}
		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("os_type", props.OsType)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)

		if encryption := props.Encryption; encryption != nil {
			d.Set("encryption_type", string(encryption.Type))
			d.Set("managed_disk_encryption_set_id", encryption.DiskEncryptionSetID)
		}

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("Error setting `encryption_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmManagedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return err
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return err
		}
	}

	return nil
}

func flattenAzureRmManagedDiskCreationData(d *schema.ResourceData, creationData *compute.CreationData) {
	d.Set("create_option", string(creationData.CreateOption))
	d.Set("source_resource_id", creationData.SourceResourceID)
	d.Set("source_uri", creationData.SourceURI)
	if ref := creationData.ImageReference; ref != nil {
		d.Set("image_reference_id", ref.ID)
	}
}
