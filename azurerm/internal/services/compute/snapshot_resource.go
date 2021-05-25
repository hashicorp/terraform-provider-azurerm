package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSnapshotCreateUpdate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotCreateUpdate,
		Delete: resourceSnapshotDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SnapshotName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"create_option": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Copy),
					string(compute.Import),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"tags": tags.Schema(),
		},
	}
}

func resourceSnapshotCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := d.Get("create_option").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_snapshot", *existing.ID)
		}
	}

	properties := compute.Snapshot{
		Location: utils.String(location),
		SnapshotProperties: &compute.SnapshotProperties{
			CreationData: &compute.CreationData{
				CreateOption: compute.DiskCreateOption(createOption),
			},
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("source_uri"); ok {
		properties.SnapshotProperties.CreationData.SourceURI = utils.String(v.(string))
	}

	if v, ok := d.GetOk("source_resource_id"); ok {
		properties.SnapshotProperties.CreationData.SourceResourceID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		properties.SnapshotProperties.CreationData.StorageAccountID = utils.String(v.(string))
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB > 0 {
		properties.SnapshotProperties.DiskSizeGB = utils.Int32(int32(diskSizeGB))
	}

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings := v.([]interface{})
		settings := encryptionSettings[0].(map[string]interface{})
		properties.EncryptionSettingsCollection = expandManagedDiskEncryptionSettings(settings)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["snapshots"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Snapshot %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Snapshot %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SnapshotProperties; props != nil {
		if data := props.CreationData; data != nil {
			d.Set("create_option", string(data.CreateOption))

			if accountId := data.StorageAccountID; accountId != nil {
				d.Set("storage_account_id", accountId)
			}
		}

		if props.DiskSizeGB != nil {
			d.Set("disk_size_gb", int(*props.DiskSizeGB))
		}

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("Error setting `encryption_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSnapshotDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["snapshots"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Snapshot: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error deleting Snapshot: %+v", err)
	}

	return nil
}
