// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimages.ParseGalleryImageID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	payload := existing.Model

	if payload == nil {
		return fmt.Errorf("model is nil for %s", id)
	}

	if d.HasChange("disk_types_not_allowed") {
		payload.Properties.Disallowed = expandGalleryImageDisallowed(d)
	}

	if d.HasChange("end_of_life_date") {
		endOfLifeDate, _ := time.Parse(time.RFC3339, d.Get("end_of_life_date").(string))
		payload.Properties.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("eula") {
		payload.Properties.Description = pointer.To(d.Get("eula").(string))
	}

	if d.HasChange("specialized") {
		if d.Get("specialized").(bool) {
			payload.Properties.OsState = galleryimages.OperatingSystemStateTypesSpecialized
		} else {
			payload.Properties.OsState = galleryimages.OperatingSystemStateTypesGeneralized
		}
	}

	if d.HasChange("release_note_uri") {
		payload.Properties.ReleaseNoteUri = pointer.To(d.Get("release_note_uri").(string))
	}

	if d.HasChanges(
		"max_recommended_vcpu_count",
		"min_recommended_vcpu_count",
		"max_recommended_memory_in_gb",
		"min_recommended_memory_in_gb",
	) {
		recommended, err := expandGalleryImageRecommended(d)
		if err != nil {
			return err
		}
		payload.Properties.Recommended = recommended
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageRead(d, meta)
}
