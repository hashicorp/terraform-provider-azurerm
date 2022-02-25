package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSnapshotCreateUpdate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotCreateUpdate,
		Delete: resourceSnapshotDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SnapshotID(id)
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
					string(compute.DiskCreateOptionCopy),
					string(compute.DiskCreateOptionImport),
				}, !features.ThreePointOh()),
				DiffSuppressFunc: suppress.CaseDifferenceV2Only,
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := d.Get("create_option").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_snapshot", id.ID())
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, properties)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Snapshot %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Snapshot %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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
			return fmt.Errorf("setting `encryption_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSnapshotDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Snapshot: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("deleting Snapshot: %+v", err)
	}

	return nil
}
