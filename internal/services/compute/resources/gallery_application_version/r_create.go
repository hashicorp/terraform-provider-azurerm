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
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			galleryApplicationId, err := galleryapplications.ParseApplicationID(state.GalleryApplicationId)
			if err != nil {
				return err
			}

			id := galleryapplicationversions.NewApplicationVersionID(subscriptionId, galleryApplicationId.ResourceGroupName, galleryApplicationId.GalleryName, galleryApplicationId.ApplicationName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, galleryapplicationversions.DefaultGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := galleryapplicationversions.GalleryApplicationVersion{
				Location: location.Normalize(state.Location),
				Properties: &galleryapplicationversions.GalleryApplicationVersionProperties{
					PublishingProfile: galleryapplicationversions.GalleryApplicationVersionPublishingProfile{
						EnableHealthCheck: pointer.To(state.EnableHealthCheck),
						ExcludeFromLatest: pointer.To(state.ExcludeFromLatest),
						ManageActions:     expandGalleryApplicationVersionManageAction(state.ManageAction),
						Source:            expandGalleryApplicationVersionSource(state.Source),
						TargetRegions:     expandGalleryApplicationVersionTargetRegion(state.TargetRegion),
					},
					SafetyProfile: &galleryapplicationversions.GalleryArtifactSafetyProfileBase{
						AllowDeletionOfReplicatedLocations: pointer.To(true),
					},
				},
				Tags: pointer.To(state.Tags),
			}

			if state.ConfigFile != "" {
				if payload.Properties.PublishingProfile.Settings == nil {
					payload.Properties.PublishingProfile.Settings = &galleryapplicationversions.UserArtifactSettings{}
				}

				payload.Properties.PublishingProfile.Settings.ConfigFileName = &state.ConfigFile
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				payload.Properties.PublishingProfile.SetEndOfLifeDateAsTime(endOfLifeDate)
			}

			if state.PackageFile != "" {
				if payload.Properties.PublishingProfile.Settings == nil {
					payload.Properties.PublishingProfile.Settings = &galleryapplicationversions.UserArtifactSettings{}
				}

				payload.Properties.PublishingProfile.Settings.PackageFileName = &state.PackageFile
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

func expandGalleryApplicationVersionManageAction(input []ManageAction) *galleryapplicationversions.UserArtifactManage {
	if len(input) == 0 {
		return &galleryapplicationversions.UserArtifactManage{}
	}
	v := input[0]
	return &galleryapplicationversions.UserArtifactManage{
		Install: v.Install,
		Remove:  v.Remove,
		Update:  pointer.To(v.Update),
	}
}

func expandGalleryApplicationVersionSource(input []Source) galleryapplicationversions.UserArtifactSource {
	if len(input) == 0 {
		return galleryapplicationversions.UserArtifactSource{}
	}
	v := input[0]
	return galleryapplicationversions.UserArtifactSource{
		MediaLink:                v.MediaLink,
		DefaultConfigurationLink: pointer.To(v.DefaultConfigurationLink),
	}
}

func expandGalleryApplicationVersionTargetRegion(input []TargetRegion) *[]galleryapplicationversions.TargetRegion {
	results := make([]galleryapplicationversions.TargetRegion, 0)
	for _, item := range input {
		targetRegion := galleryapplicationversions.TargetRegion{
			Name:                 location.Normalize(item.Name),
			RegionalReplicaCount: pointer.To(item.RegionalReplicaCount),
			StorageAccountType:   pointer.ToEnum[galleryapplicationversions.StorageAccountType](item.StorageAccountType),
		}

		if item.ExcludeFromLatest {
			targetRegion.ExcludeFromLatest = &item.ExcludeFromLatest
		}

		results = append(results, targetRegion)
	}

	return &results
}
