// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageGalleryCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	gallerySharingUpdateClient := meta.(*clients.Client).Compute.GallerySharingUpdateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSharedImageGalleryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, galleries.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_shared_image_gallery", id.ID())
		}
	}

	sharing, permission, err := expandSharedImageGallerySharing(d.Get("sharing").([]any))
	if err != nil {
		return fmt.Errorf("expanding `sharing`: %+v", err)
	}

	payload := galleries.Gallery{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleries.GalleryProperties{
			Description:    pointer.To(d.Get("description").(string)),
			SharingProfile: sharing,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	if permission == galleries.GallerySharingPermissionTypesCommunity {
		updatePayload := gallerysharingupdate.SharingUpdate{
			OperationType: gallerysharingupdate.SharingUpdateOperationTypesEnableCommunity,
		}
		if err = gallerySharingUpdateClient.GallerySharingProfileUpdateThenPoll(ctx, id, updatePayload); err != nil {
			return fmt.Errorf("enabling community sharing of %s: %+v", id, err)
		}
	}

	return resourceSharedImageGalleryRead(d, meta)
}

func expandSharedImageGallerySharing(input []any) (*galleries.SharingProfile, galleries.GallerySharingPermissionTypes, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, "", nil
	}

	v := input[0].(map[string]any)
	permission := galleries.GallerySharingPermissionTypes(v["permission"].(string))
	communityGallery := v["community_gallery"].([]any)

	if permission == galleries.GallerySharingPermissionTypesCommunity {
		if len(communityGallery) == 0 || communityGallery[0] == nil {
			return nil, permission, fmt.Errorf("`community_gallery` must be set when `permission` is set to `Community`")
		}
	}

	return &galleries.SharingProfile{
		Permissions:          pointer.To(permission),
		CommunityGalleryInfo: expandSharedImageGalleryCommunityGallery(communityGallery),
	}, permission, nil
}

func expandSharedImageGalleryCommunityGallery(input []any) *galleries.CommunityGalleryInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	return &galleries.CommunityGalleryInfo{
		Eula:             pointer.To(v["eula"].(string)),
		PublicNamePrefix: pointer.To(v["prefix"].(string)),
		PublisherContact: pointer.To(v["publisher_email"].(string)),
		PublisherUri:     pointer.To(v["publisher_uri"].(string)),
	}
}
