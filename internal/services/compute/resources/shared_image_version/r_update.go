// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageVersionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, galleryimageversions.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	payload := existing.Model
	if payload == nil {
		return fmt.Errorf("model is nil for %s", id)
	}

	if payload.Properties == nil {
		return fmt.Errorf("properties is nil for %s", id)
	}

	if payload.Properties.PublishingProfile == nil {
		payload.Properties.PublishingProfile = &galleryimageversions.GalleryArtifactPublishingProfileBase{}
	}

	if d.HasChange("target_region") {
		targetRegions, err := expandSharedImageVersionTargetRegions(d)
		if err != nil {
			return err
		}

		payload.Properties.PublishingProfile.TargetRegions = targetRegions
	}

	if d.HasChange("end_of_life_date") {
		endOfLifeDate, _ := time.Parse(time.RFC3339, d.Get("end_of_life_date").(string))
		payload.Properties.PublishingProfile.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if d.HasChange("exclude_from_latest") {
		payload.Properties.PublishingProfile.ExcludeFromLatest = pointer.To(d.Get("exclude_from_latest").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageVersionRead(d, meta)
}
