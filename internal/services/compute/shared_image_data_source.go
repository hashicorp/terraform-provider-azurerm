// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSharedImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSharedImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"architecture": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"specialized": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identifier": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"publisher": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"offer": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"eula": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"privacy_statement_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"purchase_plan": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"publisher": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"product": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"release_note_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceSharedImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimages.NewGalleryImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			architecture := ""
			if props.Architecture != nil {
				architecture = string(*props.Architecture)
			}
			d.Set("architecture", architecture)
			d.Set("description", props.Description)
			d.Set("eula", props.Eula)
			hyperVGeneration := ""
			if props.HyperVGeneration != nil {
				hyperVGeneration = string(*props.HyperVGeneration)
			}
			d.Set("hyper_v_generation", hyperVGeneration)
			if err := d.Set("identifier", flattenGalleryImageDataSourceIdentifier(props.Identifier)); err != nil {
				return fmt.Errorf("setting `identifier`: %+v", err)
			}
			d.Set("os_type", string(props.OsType))
			d.Set("privacy_statement_uri", props.PrivacyStatementUri)
			if err := d.Set("purchase_plan", flattenGalleryImageDataSourcePurchasePlan(props.PurchasePlan)); err != nil {
				return fmt.Errorf("setting `purchase_plan`: %+v", err)
			}
			d.Set("release_note_uri", props.ReleaseNoteUri)
			d.Set("specialized", props.OsState == galleryimages.OperatingSystemStateTypesSpecialized)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenGalleryImageDataSourceIdentifier(input galleryimages.GalleryImageIdentifier) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"offer":     input.Offer,
			"publisher": input.Publisher,
			"sku":       input.Sku,
		},
	}
}

func flattenGalleryImageDataSourcePurchasePlan(input *galleryimages.ImagePurchasePlan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	product := ""
	if input.Product != nil {
		product = *input.Product
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"publisher": publisher,
			"product":   product,
		},
	}
}
