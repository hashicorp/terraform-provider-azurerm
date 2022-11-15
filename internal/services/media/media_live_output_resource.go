package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveevents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveoutputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaLiveOutput() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaLiveOutputCreate,
		Read:   resourceMediaLiveOutputRead,
		Delete: resourceMediaLiveOutputDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := liveoutputs.ParseLiveOutputID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LiveOutputV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"live_event_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: liveevents.ValidateLiveEventID,
			},

			"archive_window_duration": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"asset_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
					"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
				),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"hls_fragments_per_ts_segment": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"manifest_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"output_snap_time_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceMediaLiveOutputCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveOutputsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	eventId, err := liveevents.ParseLiveEventID(d.Get("live_event_id").(string))
	if err != nil {
		return err
	}
	id := liveoutputs.NewLiveOutputID(subscriptionID, eventId.ResourceGroupName, eventId.AccountName, eventId.LiveEventName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.AccountName, id.LiveEventName, id.LiveOutputName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_live_event_output", id.ID())
		}
	}

	parameters := media.LiveOutput{
		LiveOutputProperties: &media.LiveOutputProperties{},
	}

	if archiveWindowLength, ok := d.GetOk("archive_window_duration"); ok {
		parameters.LiveOutputProperties.ArchiveWindowLength = utils.String(archiveWindowLength.(string))
	}

	if assetName, ok := d.GetOk("asset_name"); ok {
		parameters.LiveOutputProperties.AssetName = utils.String(assetName.(string))
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.LiveOutputProperties.Description = utils.String(description.(string))
	}

	if hlsFragmentsPerTsSegment, ok := d.GetOk("hls_fragments_per_ts_segment"); ok {
		parameters.LiveOutputProperties.Hls = &media.Hls{
			FragmentsPerTsSegment: utils.Int32(int32(hlsFragmentsPerTsSegment.(int))),
		}
	}

	if manifestName, ok := d.GetOk("manifest_name"); ok {
		parameters.LiveOutputProperties.ManifestName = utils.String(manifestName.(string))
	}

	if outputSnapTime, ok := d.GetOk("output_snap_time_in_seconds"); ok {
		parameters.LiveOutputProperties.OutputSnapTime = utils.Int64(int64(outputSnapTime.(int)))
	}

	future, err := client.Create(ctx, id.ResourceGroupName, id.AccountName, id.LiveEventName, id.LiveOutputName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaLiveOutputRead(d, meta)
}

func resourceMediaLiveOutputRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveOutputsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveoutputs.ParseLiveOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.AccountName, id.LiveEventName, id.LiveOutputName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LiveOutputName)
	d.Set("live_event_id", liveevents.NewLiveEventID(subscriptionID, id.ResourceGroupName, id.AccountName, id.LiveEventName).ID())

	if props := resp.LiveOutputProperties; props != nil {
		d.Set("archive_window_duration", props.ArchiveWindowLength)
		d.Set("asset_name", props.AssetName)
		d.Set("description", props.Description)

		var hlsFragmentsPerTsSegment int32
		if props.Hls != nil && props.Hls.FragmentsPerTsSegment != nil {
			hlsFragmentsPerTsSegment = *props.Hls.FragmentsPerTsSegment
		}
		d.Set("hls_fragments_per_ts_segment", hlsFragmentsPerTsSegment)
		d.Set("manifest_name", props.ManifestName)

		var outputSnapTime int64
		if props.OutputSnapTime != nil {
			outputSnapTime = *props.OutputSnapTime
		}
		d.Set("output_snap_time_in_seconds", outputSnapTime)
	}

	return nil
}

func resourceMediaLiveOutputDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveOutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveoutputs.ParseLiveOutputID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.AccountName, id.LiveEventName, id.LiveOutputName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to delete: %+v", id, err)
	}

	return nil
}
