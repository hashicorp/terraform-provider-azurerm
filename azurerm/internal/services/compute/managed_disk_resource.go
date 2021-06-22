package compute

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceManagedDisk() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedDiskCreate,
		Read:   resourceManagedDiskRead,
		Update: resourceManagedDiskUpdate,
		Delete: resourceManagedDiskDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedDiskID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaSingleZone(),

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.StandardLRS),
					string(compute.PremiumLRS),
					string(compute.StandardSSDLRS),
					string(compute.UltraSSDLRS),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"create_option": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Empty),
					string(compute.FromImage),
					string(compute.Import),
					string(compute.Restore),
				}, false),
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true, // Not supported by disk update
				ValidateFunc: azure.ValidateResourceID,
			},

			"image_reference_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Windows),
					string(compute.Linux),
				}, true),
			},

			"disk_size_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ManagedDiskSizeGB,
			},

			"disk_iops_read_write": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/8132
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"network_access_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.AllowAll),
					string(compute.AllowPrivate),
					string(compute.DenyAll),
				}, false),
			},
			"disk_access_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: make this case-sensitive once this bug in the Azure API has been fixed:
				//       https://github.com/Azure/azure-rest-api-specs/issues/14192
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceManagedDiskCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewManagedDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.DiskName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Managed Disk %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_disk", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := compute.DiskCreateOption(d.Get("create_option").(string))
	storageAccountType := d.Get("storage_account_type").(string)
	osType := d.Get("os_type").(string)

	t := d.Get("tags").(map[string]interface{})
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))
	skuName := compute.DiskStorageAccountTypes(storageAccountType)

	props := &compute.DiskProperties{
		CreationData: &compute.CreationData{
			CreateOption: createOption,
		},
		OsType: compute.OperatingSystemTypes(osType),
		Encryption: &compute.Encryption{
			Type: compute.EncryptionTypeEncryptionAtRestWithPlatformKey,
		},
	}

	if v := d.Get("disk_size_gb"); v != 0 {
		diskSize := int32(v.(int))
		props.DiskSizeGB = &diskSize
	}

	if storageAccountType == string(compute.UltraSSDLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			props.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			props.DiskMBpsReadWrite = &diskMBps
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") {
		return fmt.Errorf("[ERROR] disk_iops_read_write and disk_mbps_read_write are only available for UltraSSD disks")
	}

	if createOption == compute.Import {
		sourceUri := d.Get("source_uri").(string)
		if sourceUri == "" {
			return fmt.Errorf("`source_uri` must be specified when `create_option` is set to `Import`")
		}

		storageAccountId := d.Get("storage_account_id").(string)
		if storageAccountId == "" {
			return fmt.Errorf("`storage_account_id` must be specified when `create_option` is set to `Import`")
		}

		props.CreationData.StorageAccountID = utils.String(storageAccountId)
		props.CreationData.SourceURI = utils.String(sourceUri)
	}
	if createOption == compute.Copy || createOption == compute.Restore {
		sourceResourceId := d.Get("source_resource_id").(string)
		if sourceResourceId == "" {
			return fmt.Errorf("`source_resource_id` must be specified when `create_option` is set to `Copy` or `Restore`")
		}

		props.CreationData.SourceResourceID = utils.String(sourceResourceId)
	}
	if createOption == compute.FromImage {
		imageReferenceId := d.Get("image_reference_id").(string)
		if imageReferenceId == "" {
			return fmt.Errorf("`image_reference_id` must be specified when `create_option` is set to `Import`")
		}

		props.CreationData.ImageReference = &compute.ImageDiskReference{
			ID: utils.String(imageReferenceId),
		}
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		props.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
		props.Encryption = &compute.Encryption{
			Type:                compute.EncryptionTypeEncryptionAtRestWithCustomerKey,
			DiskEncryptionSetID: utils.String(diskEncryptionSetId),
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		props.NetworkAccessPolicy = compute.NetworkAccessPolicy(networkAccessPolicy)
	} else {
		props.NetworkAccessPolicy = compute.AllowAll
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case props.NetworkAccessPolicy == compute.AllowPrivate:
			props.DiskAccessID = utils.String(diskAccessID)
		case diskAccessID != "" && props.NetworkAccessPolicy != compute.AllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			props.DiskAccessID = nil
		}
	}

	if tier := d.Get("tier").(string); tier != "" {
		if storageAccountType != string(compute.PremiumZRS) && storageAccountType != string(compute.PremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		props.Tier = &tier
	}

	createDisk := compute.Disk{
		Name:           &name,
		Location:       &location,
		DiskProperties: props,
		Sku: &compute.DiskSku{
			Name: skuName,
		},
		Tags:  tags.Expand(t),
		Zones: zones,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, createDisk)
	if err != nil {
		return fmt.Errorf("Error creating/updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for create/update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Error reading Managed Disk %s (Resource Group %q): ID was nil", name, resourceGroup)
	}

	d.SetId(id.ID())

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Managed Disk update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	storageAccountType := d.Get("storage_account_type").(string)
	shouldShutDown := false

	disk, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(disk.Response) {
			return fmt.Errorf("Error Managed Disk %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	diskUpdate := compute.DiskUpdate{
		DiskUpdateProperties: &compute.DiskUpdateProperties{},
	}

	if d.HasChange("tier") {
		if storageAccountType != string(compute.PremiumZRS) && storageAccountType != string(compute.PremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		shouldShutDown = true
		tier := d.Get("tier").(string)
		diskUpdate.Tier = &tier
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		diskUpdate.Tags = tags.Expand(t)
	}

	if d.HasChange("storage_account_type") {
		shouldShutDown = true
		var skuName compute.DiskStorageAccountTypes
		for _, v := range compute.PossibleDiskStorageAccountTypesValues() {
			if strings.EqualFold(storageAccountType, string(v)) {
				skuName = v
			}
		}
		diskUpdate.Sku = &compute.DiskSku{
			Name: skuName,
		}
	}

	if strings.EqualFold(storageAccountType, string(compute.UltraSSDLRS)) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskIOPS := int64(v.(int))
			diskUpdate.DiskIOPSReadWrite = &diskIOPS
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskMBps := int64(v.(int))
			diskUpdate.DiskMBpsReadWrite = &diskMBps
		}
	} else if d.HasChange("disk_iops_read_write") || d.HasChange("disk_mbps_read_write") {
		return fmt.Errorf("[ERROR] disk_iops_read_write and disk_mbps_read_write are only available for UltraSSD disks")
	}

	if d.HasChange("os_type") {
		diskUpdate.DiskUpdateProperties.OsType = compute.OperatingSystemTypes(d.Get("os_type").(string))
	}

	if d.HasChange("disk_size_gb") {
		if old, new := d.GetChange("disk_size_gb"); new.(int) > old.(int) {
			shouldShutDown = true
			diskUpdate.DiskUpdateProperties.DiskSizeGB = utils.Int32(int32(new.(int)))
		} else {
			return fmt.Errorf("Error - New size must be greater than original size. Shrinking disks is not supported on Azure")
		}
	}

	if d.HasChange("disk_encryption_set_id") {
		shouldShutDown = true
		if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskUpdate.Encryption = &compute.Encryption{
				Type:                compute.EncryptionTypeEncryptionAtRestWithCustomerKey,
				DiskEncryptionSetID: utils.String(diskEncryptionSetId),
			}
		} else {
			return fmt.Errorf("Once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if networkAccessPolicy := d.Get("network_access_policy").(string); networkAccessPolicy != "" {
		diskUpdate.NetworkAccessPolicy = compute.NetworkAccessPolicy(networkAccessPolicy)
	} else {
		diskUpdate.NetworkAccessPolicy = compute.AllowAll
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		switch {
		case diskUpdate.NetworkAccessPolicy == compute.AllowPrivate:
			diskUpdate.DiskAccessID = utils.String(diskAccessID)
		case diskAccessID != "" && diskUpdate.NetworkAccessPolicy != compute.AllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			diskUpdate.DiskAccessID = nil
		}
	}

	// whilst we need to shut this down, if we're not attached to anything there's no point
	if shouldShutDown && disk.ManagedBy == nil {
		shouldShutDown = false
	}

	// if we are attached to a VM we bring down the VM as necessary for the operations which are not allowed while it's online
	if shouldShutDown {
		virtualMachine, err := parse.VirtualMachineID(*disk.ManagedBy)
		if err != nil {
			return fmt.Errorf("Error parsing VMID %q for disk attachment: %+v", *disk.ManagedBy, err)
		}
		// check instanceView State
		vmClient := meta.(*clients.Client).Compute.VMClient

		locks.ByName(name, virtualMachineResourceName)
		defer locks.UnlockByName(name, virtualMachineResourceName)

		instanceView, err := vmClient.InstanceView(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
		if err != nil {
			return fmt.Errorf("Error retrieving InstanceView for Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
		}

		shouldTurnBackOn := true
		shouldDeallocate := true

		if instanceView.Statuses != nil {
			for _, status := range *instanceView.Statuses {
				if status.Code == nil {
					continue
				}

				// could also be the provisioning state which we're not bothered with here
				state := strings.ToLower(*status.Code)
				if !strings.HasPrefix(state, "powerstate/") {
					continue
				}

				state = strings.TrimPrefix(state, "powerstate/")
				switch strings.ToLower(state) {
				case "deallocated":
				case "deallocating":
					shouldTurnBackOn = false
					shouldShutDown = false
					shouldDeallocate = false
				case "stopping":
				case "stopped":
					shouldShutDown = false
					shouldTurnBackOn = false
				}
			}
		}

		// Shutdown
		if shouldShutDown {
			log.Printf("[DEBUG] Shutting Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			forceShutdown := false
			future, err := vmClient.PowerOff(ctx, virtualMachine.ResourceGroup, virtualMachine.Name, utils.Bool(forceShutdown))
			if err != nil {
				return fmt.Errorf("Error sending Power Off to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for Power Off of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Shut Down Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// De-allocate
		if shouldDeallocate {
			log.Printf("[DEBUG] Deallocating Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			deAllocFuture, err := vmClient.Deallocate(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
			if err != nil {
				return fmt.Errorf("Error Deallocating to Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := deAllocFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for Deallocation of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Deallocated Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}

		// Update Disk
		updateFuture, err := client.Update(ctx, resourceGroup, name, diskUpdate)
		if err != nil {
			return fmt.Errorf("Error updating Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if err := updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for update of Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if shouldTurnBackOn {
			log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
			future, err := vmClient.Start(ctx, virtualMachine.ResourceGroup, virtualMachine.Name)
			if err != nil {
				return fmt.Errorf("Error starting Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for start of Virtual Machine %q (Resource Group %q): %+v", virtualMachine.Name, virtualMachine.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Started Virtual Machine %q (Resource Group %q)..", virtualMachine.Name, virtualMachine.ResourceGroup)
		}
	} else { // otherwise, just update it
		diskFuture, err := client.Update(ctx, resourceGroup, name, diskUpdate)
		if err != nil {
			return fmt.Errorf("Error expanding managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = diskFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for expand operation on managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceManagedDiskRead(d, meta)
}

func resourceManagedDiskRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Managed Disk %s (resource group %s): %s", id.DiskName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			d.Set("create_option", string(creationData.CreateOption))

			imageReferenceID := ""
			if creationData.ImageReference != nil && creationData.ImageReference.ID != nil {
				imageReferenceID = *creationData.ImageReference.ID
			}
			d.Set("image_reference_id", imageReferenceID)

			d.Set("source_resource_id", creationData.SourceResourceID)
			d.Set("source_uri", creationData.SourceURI)
			d.Set("storage_account_id", creationData.StorageAccountID)
		}

		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
		d.Set("os_type", props.OsType)
		d.Set("tier", props.Tier)

		if networkAccessPolicy := props.NetworkAccessPolicy; networkAccessPolicy != compute.AllowAll {
			d.Set("network_access_policy", props.NetworkAccessPolicy)
		}
		d.Set("disk_access_id", props.DiskAccessID)

		diskEncryptionSetId := ""
		if props.Encryption != nil && props.Encryption.DiskEncryptionSetID != nil {
			diskEncryptionSetId = *props.Encryption.DiskEncryptionSetID
		}
		d.Set("disk_encryption_set_id", diskEncryptionSetId)

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("Error setting `encryption_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceManagedDiskDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DiskName)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Managed Disk %q (Resource Group %q): %+v", id.DiskName, id.ResourceGroup, err)
	}

	return nil
}
