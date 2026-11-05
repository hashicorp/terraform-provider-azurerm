// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageGalleryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSharedImageGalleryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			uniqueName := ""
			if props.Identifier != nil && props.Identifier.UniqueName != nil {
				uniqueName = *props.Identifier.UniqueName
			}
			d.Set("unique_name", uniqueName)

			d.Set("sharing", flattenSharedImageGallerySharing(props.SharingProfile))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenSharedImageGallerySharing(input *galleries.SharingProfile) []any {
	if input == nil {
		return make([]any, 0)
	}

	permission := ""
	if v := input.Permissions; v != nil {
		permission = pointer.FromEnum(v)
	}

	return []any{
		map[string]any{
			"permission":        permission,
			"community_gallery": flattenSharedImageGalleryCommunityGallery(input.CommunityGalleryInfo),
		},
	}
}

func flattenSharedImageGalleryCommunityGallery(input *galleries.CommunityGalleryInfo) []any {
	if input == nil {
		return make([]any, 0)
	}

	eula := ""
	if input.Eula != nil {
		eula = pointer.From(input.Eula)
	}

	publicName := ""
	if input.PublicNames != nil {
		if v := pointer.From(input.PublicNames); len(v) > 0 {
			publicName = v[0]
		}
	}

	publicNamePrefix := ""
	if input.PublicNamePrefix != nil {
		publicNamePrefix = pointer.From(input.PublicNamePrefix)
	}

	publisherEmail := ""
	if input.PublisherContact != nil {
		publisherEmail = pointer.From(input.PublisherContact)
	}

	publisherUri := ""
	if input.PublisherUri != nil {
		publisherUri = pointer.From(input.PublisherUri)
	}

	return []any{
		map[string]any{
			"eula":            eula,
			"name":            publicName,
			"prefix":          publicNamePrefix,
			"publisher_email": publisherEmail,
			"publisher_uri":   publisherUri,
		},
	}
}
