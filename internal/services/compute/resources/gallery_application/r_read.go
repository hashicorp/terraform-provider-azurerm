// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient
			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%q was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &GalleryApplicationModel{
				Name:      id.ApplicationName,
				GalleryId: commonids.NewSharedImageGalleryID(id.SubscriptionId, id.ResourceGroupName, id.GalleryName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					if v := props.Description; v != nil {
						state.Description = *props.Description
					}

					if v := props.EndOfLifeDate; v != nil {
						d, err := props.GetEndOfLifeDateAsTime()
						if err != nil {
							return fmt.Errorf("parsing `end_of_life_date` from API Response: %+v", err)
						}
						if d != nil {
							state.EndOfLifeDate = d.Format(time.RFC3339)
						}
					}

					if v := props.Eula; v != nil {
						state.Eula = *props.Eula
					}

					if v := props.PrivacyStatementUri; v != nil {
						state.PrivacyStatementURI = *props.PrivacyStatementUri
					}

					if v := props.ReleaseNoteUri; v != nil {
						state.ReleaseNoteURI = *props.ReleaseNoteUri
					}

					state.SupportedOSType = string(props.SupportedOSType)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r Resource) Identity() resourceids.ResourceId {
	return &galleryapplications.ApplicationId{}
}
