package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-01-02-preview/containerservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSnapshots() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSnapshotsCreateUpdate,
		Read:   resourceSnapshotsRead,
		Update: resourceSnapshotsCreateUpdate,
		Delete: resourceSnapshotsDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SnapshotsID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"source_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"snapshot_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.SnapshotTypeNodePool),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSnapshotsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSnapshotsID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SnapshotName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_snapshots", id.ID())
		}
	}

	properties := containerservice.SnapshotProperties{
		CreationData: &containerservice.CreationData{},
	}

	if v, ok := d.GetOk("source_resource_id"); ok {
		properties.CreationData.SourceResourceID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_type"); ok {
		properties.SnapshotType = containerservice.SnapshotType(v.(string))
	}

	parameters := containerservice.Snapshot{
		Location:           &location,
		SnapshotProperties: &properties,
		Tags:               tags.Expand(t),
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SnapshotName, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSnapshotsRead(d, meta)
}

func resourceSnapshotsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.SnapshotsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SnapshotName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Snapshots %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Snapshots %q: %+v", id.SnapshotName, err)
	}

	d.Set("name", id.SnapshotName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SnapshotProperties; props != nil {
		if data := props.CreationData; data != nil {
			d.Set("source_resource_id", data.SourceResourceID)
		}

		d.Set("snapshot_type", string(props.SnapshotType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSnapshotsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.SnapshotsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotsID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.SnapshotName); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
