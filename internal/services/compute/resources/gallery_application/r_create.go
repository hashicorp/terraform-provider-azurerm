// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state GalleryApplicationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Compute.GalleryApplicationsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			galleryId, err := commonids.ParseSharedImageGalleryID(state.GalleryId)
			if err != nil {
				return err
			}

			id := galleryapplications.NewApplicationID(subscriptionId, galleryId.ResourceGroupName, galleryId.GalleryName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := galleryapplications.GalleryApplication{
				Location: location.Normalize(state.Location),
				Properties: &galleryapplications.GalleryApplicationProperties{
					SupportedOSType: galleryapplications.OperatingSystemTypes(state.SupportedOSType),
				},
				Tags: pointer.To(state.Tags),
			}

			if state.Description != "" {
				payload.Properties.Description = pointer.To(state.Description)
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				payload.Properties.SetEndOfLifeDateAsTime(endOfLifeDate)
			}

			if state.Eula != "" {
				payload.Properties.Eula = pointer.To(state.Eula)
			}

			if state.PrivacyStatementURI != "" {
				payload.Properties.PrivacyStatementUri = pointer.To(state.PrivacyStatementURI)
			}

			if state.ReleaseNoteURI != "" {
				payload.Properties.ReleaseNoteUri = pointer.To(state.ReleaseNoteURI)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
		Timeout: 30 * time.Minute,
	}
}
