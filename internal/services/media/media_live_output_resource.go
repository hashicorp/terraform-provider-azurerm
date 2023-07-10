// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveevents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveoutputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
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

		DeprecationMessage: azureMediaRetirementMessage,

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
				ValidateFunc: validate.ISO8601DurationBetween("PT1M", "PT25H"),
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

			"rewind_window_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.ISO8601DurationBetween("PT1M", "PT25H"),
			},
		},
	}
}

func resourceMediaLiveOutputCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveOutputs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	eventId, err := liveevents.ParseLiveEventID(d.Get("live_event_id").(string))
	if err != nil {
		return err
	}
	id := liveoutputs.NewLiveOutputID(eventId.SubscriptionId, eventId.ResourceGroupName, eventId.MediaServiceName, eventId.LiveEventName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_live_event_output", id.ID())
		}
	}

	parameters := liveoutputs.LiveOutput{
		Properties: &liveoutputs.LiveOutputProperties{},
	}

	if archiveWindowLength, ok := d.GetOk("archive_window_duration"); ok {
		parameters.Properties.ArchiveWindowLength = archiveWindowLength.(string)
	}

	if assetName, ok := d.GetOk("asset_name"); ok {
		parameters.Properties.AssetName = assetName.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.Properties.Description = utils.String(description.(string))
	}

	if hlsFragmentsPerTsSegment, ok := d.GetOk("hls_fragments_per_ts_segment"); ok {
		parameters.Properties.Hls = &liveoutputs.Hls{
			FragmentsPerTsSegment: pointer.To(int64(hlsFragmentsPerTsSegment.(int))),
		}
	}

	if manifestName, ok := d.GetOk("manifest_name"); ok {
		parameters.Properties.ManifestName = utils.String(manifestName.(string))
	}

	if outputSnapTime, ok := d.GetOk("output_snap_time_in_seconds"); ok {
		parameters.Properties.OutputSnapTime = utils.Int64(int64(outputSnapTime.(int)))
	}

	if rewindWindowLength, ok := d.GetOk("rewind_window_duration"); ok {
		parameters.Properties.RewindWindowLength = utils.String(rewindWindowLength.(string))
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaLiveOutputRead(d, meta)
}

func resourceMediaLiveOutputRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveOutputs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveoutputs.ParseLiveOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LiveOutputName)
	d.Set("live_event_id", liveevents.NewLiveEventID(id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.LiveEventName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("archive_window_duration", props.ArchiveWindowLength)
			d.Set("asset_name", props.AssetName)
			d.Set("description", props.Description)

			var hlsFragmentsPerTsSegment int64
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
			d.Set("rewind_window_duration", props.RewindWindowLength)
		}
	}

	return nil
}

func resourceMediaLiveOutputDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveOutputs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveoutputs.ParseLiveOutputID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
