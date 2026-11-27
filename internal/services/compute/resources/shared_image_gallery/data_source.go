// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Shared Image Gallery.",
		Read:        dataSourceSharedImageGalleryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Shared Image Gallery.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"description": {
				Description: "A description for the Shared Image Gallery.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"image_names": {
				Description: "A list of Shared Image names within this Shared Image Gallery.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"unique_name": {
				Description: "The unique name assigned to the Shared Image Gallery.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceSharedImageGalleryRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).GalleriesClient
	imagesClient := computeClient.From(meta).GalleryImagesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := commonids.NewSharedImageGalleryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, galleries.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	imagesResp, err := imagesClient.ListByGalleryComplete(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	imageNames := make([]string, 0)
	for _, image := range imagesResp.Items {
		if image.Name != nil {
			imageNames = append(imageNames, *image.Name)
		}
	}

	d.Set("image_names", imageNames)

	return nil
}
