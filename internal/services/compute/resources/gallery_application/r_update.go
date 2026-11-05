// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient

			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			payload := *existing.Model
			// lintignore:R019 // deliberate subset: only the fields feeding payload.Properties; tags are applied separately below
			if metadata.ResourceData.HasChanges(
				"description",
				"end_of_life_date",
				"eula",
				"privacy_statement_uri",
				"release_note_uri",
			) {
				if payload.Properties == nil {
					payload.Properties = &galleryapplications.GalleryApplicationProperties{}
				}

				if metadata.ResourceData.HasChange("description") {
					payload.Properties.Description = pointer.To(state.Description)
				}

				if metadata.ResourceData.HasChange("end_of_life_date") {
					endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
					payload.Properties.SetEndOfLifeDateAsTime(endOfLifeDate)
				}

				if metadata.ResourceData.HasChange("eula") {
					payload.Properties.Eula = pointer.To(state.Eula)
				}

				if metadata.ResourceData.HasChange("privacy_statement_uri") {
					payload.Properties.PrivacyStatementUri = pointer.To(state.PrivacyStatementURI)
				}

				if metadata.ResourceData.HasChange("release_note_uri") {
					payload.Properties.ReleaseNoteUri = pointer.To(state.ReleaseNoteURI)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
