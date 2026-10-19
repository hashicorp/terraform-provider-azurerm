// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient

			id, err := galleryapplicationversions.ParseApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			payload := galleryapplicationversions.GalleryApplicationVersionUpdate{}

			// lintignore:R019 // deliberate subset: only the fields feeding payload.Properties; tags are applied separately below
			if metadata.ResourceData.HasChanges(
				"enable_health_check",
				"end_of_life_date",
				"exclude_from_latest",
				"manage_actions",
				"source",
				"target_region",
			) {
				if payload.Properties == nil {
					payload.Properties = &galleryapplicationversions.GalleryApplicationVersionProperties{}
				}

				if metadata.ResourceData.HasChange("enable_health_check") {
					payload.Properties.PublishingProfile.EnableHealthCheck = pointer.To(state.EnableHealthCheck)
				}

				if metadata.ResourceData.HasChange("end_of_life_date") {
					endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
					payload.Properties.PublishingProfile.SetEndOfLifeDateAsTime(endOfLifeDate)
				}

				if metadata.ResourceData.HasChange("exclude_from_latest") {
					payload.Properties.PublishingProfile.ExcludeFromLatest = pointer.To(state.ExcludeFromLatest)
				}

				if metadata.ResourceData.HasChange("manage_actions") {
					payload.Properties.PublishingProfile.ManageActions = expandGalleryApplicationVersionManageAction(state.ManageAction)
				}

				if metadata.ResourceData.HasChange("source") {
					payload.Properties.PublishingProfile.Source = expandGalleryApplicationVersionSource(state.Source)
				}

				if metadata.ResourceData.HasChange("target_region") {
					payload.Properties.PublishingProfile.TargetRegions = expandGalleryApplicationVersionTargetRegion(state.TargetRegion)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
