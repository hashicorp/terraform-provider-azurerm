// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimages.NewGalleryImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_shared_image", id.ID())
		}
	}

	recommended, err := expandGalleryImageRecommended(d)
	if err != nil {
		return err
	}

	image := galleryimages.GalleryImage{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleryimages.GalleryImageProperties{
			Description:         pointer.To(d.Get("description").(string)),
			Disallowed:          expandGalleryImageDisallowed(d),
			Identifier:          expandGalleryImageIdentifier(d),
			PrivacyStatementUri: pointer.To(d.Get("privacy_statement_uri").(string)),
			ReleaseNoteUri:      pointer.To(d.Get("release_note_uri").(string)),
			Architecture:        pointer.ToEnum[galleryimages.Architecture](d.Get("architecture").(string)),
			OsType:              galleryimages.OperatingSystemTypes(d.Get("os_type").(string)),
			HyperVGeneration:    pointer.ToEnum[galleryimages.HyperVGeneration](d.Get("hyper_v_generation").(string)),
			PurchasePlan:        expandGalleryImagePurchasePlan(d.Get("purchase_plan").([]any)),
			Features:            expandSharedImageFeatures(d),
			Recommended:         recommended,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("end_of_life_date"); ok {
		endOfLifeDate, _ := time.Parse(time.RFC3339, v.(string))
		image.Properties.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if v, ok := d.GetOk("eula"); ok {
		image.Properties.Eula = pointer.To(v.(string))
	}

	if d.Get("specialized").(bool) {
		image.Properties.OsState = galleryimages.OperatingSystemStateTypesSpecialized
	} else {
		image.Properties.OsState = galleryimages.OperatingSystemStateTypesGeneralized
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, image, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceSharedImageRead(d, meta)
}
