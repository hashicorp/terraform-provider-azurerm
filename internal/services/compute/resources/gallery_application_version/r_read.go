// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			id, err := galleryapplicationversions.ParseApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, galleryapplicationversions.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &GalleryApplicationVersionModel{
				Name:                 id.VersionName,
				GalleryApplicationId: galleryapplications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.GalleryName, id.ApplicationName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					if props.PublishingProfile.EnableHealthCheck != nil {
						state.EnableHealthCheck = *props.PublishingProfile.EnableHealthCheck
					}

					if props.PublishingProfile.EndOfLifeDate != nil {
						d, err := props.PublishingProfile.GetEndOfLifeDateAsTime()
						if err != nil {
							return fmt.Errorf("parsing API response for `end_of_life_date`: %+v", err)
						}
						state.EndOfLifeDate = d.Format(time.RFC3339)
					}

					state.ExcludeFromLatest = pointer.From(props.PublishingProfile.ExcludeFromLatest)

					state.ConfigFile = ""
					state.PackageFile = ""
					if props.PublishingProfile.Settings != nil {
						state.ConfigFile = pointer.From(props.PublishingProfile.Settings.ConfigFileName)
						state.PackageFile = pointer.From(props.PublishingProfile.Settings.PackageFileName)
					}

					state.ManageAction = flattenGalleryApplicationVersionManageAction(props.PublishingProfile.ManageActions)
					state.Source = flattenGalleryApplicationVersionSource(props.PublishingProfile.Source)
					state.TargetRegion = flattenGalleryApplicationVersionTargetRegion(props.PublishingProfile.TargetRegions)
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
	return &galleryapplicationversions.ApplicationVersionId{}
}

type Source struct {
	MediaLink                string `tfschema:"media_link"`
	DefaultConfigurationLink string `tfschema:"default_configuration_link"`
}

type ManageAction struct {
	Install string `tfschema:"install"`
	Remove  string `tfschema:"remove"`
	Update  string `tfschema:"update"`
}

type TargetRegion struct {
	Name                 string `tfschema:"name"`
	RegionalReplicaCount int64  `tfschema:"regional_replica_count"`
	ExcludeFromLatest    bool   `tfschema:"exclude_from_latest"`
	StorageAccountType   string `tfschema:"storage_account_type"`
}
