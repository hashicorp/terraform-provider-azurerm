// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePlatformImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePlatformImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"publisher": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"offer": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourcePlatformImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachineimages.NewSkuID(subscriptionId, location.Normalize(d.Get("location").(string)), d.Get("publisher").(string), d.Get("offer").(string), d.Get("sku").(string))

	result, err := client.List(ctx, id, virtualmachineimages.DefaultListOperationOptions())
	if err != nil || result.Model == nil {
		return fmt.Errorf("retrieving: %+v", err)
	}

	var image *virtualmachineimages.VirtualMachineImageResource
	if v, ok := d.GetOk("version"); ok {
		version := v.(string)
		for _, item := range *result.Model {
			if item.Name == version {
				image = &item
				break
			}
		}
		if image == nil {
			return fmt.Errorf("could not find image for %s: %+v", id, err)
		}
	} else {
		// get the latest image (the last value is the latest, apparently)
		// list can be empty if user hasn't licensed any matching images
		if len(*result.Model) == 0 {
			return fmt.Errorf("no images available to this user for %s", id)
		}
		image = &(*result.Model)[len(*result.Model)-1]
	}

	imageId, err := virtualmachineimages.ParseSkuVersionIDInsensitively(*image.Id)
	if err != nil {
		return err
	}
	d.SetId(imageId.ID())

	d.Set("location", location.Normalize(image.Location))
	d.Set("publisher", id.PublisherName)
	d.Set("offer", id.OfferName)
	d.Set("sku", id.SkuName)
	d.Set("version", image.Name)

	return nil
}
