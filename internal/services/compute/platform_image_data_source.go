// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
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
	client := meta.(*clients.Client).Compute.VMImageClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	publisher := d.Get("publisher").(string)
	offer := d.Get("offer").(string)
	sku := d.Get("sku").(string)

	result, err := client.List(ctx, location, publisher, offer, sku, "", utils.Int32(int32(1000)), "name")
	if err != nil || result.Value == nil {
		return fmt.Errorf("reading Platform Images: %+v", err)
	}

	var image *compute.VirtualMachineImageResource
	if v, ok := d.GetOk("version"); ok {
		version := v.(string)
		for _, item := range *result.Value {
			if item.Name != nil && *item.Name == version {
				image = &item
				break
			}
		}
		if image == nil {
			return fmt.Errorf("could not find image (location %q / publisher %q / offer %q / sku %q / version % q): %+v", location, publisher, offer, sku, version, err)
		}
	} else {
		// get the latest image (the last value is the latest, apparently)
		// list can be empty if user hasn't licensed any matching images
		if len(*result.Value) == 0 {
			return fmt.Errorf("no images available to this user (location %q / publisher %q / offer %q / sku %q)", location, publisher, offer, sku)
		}
		image = &(*result.Value)[len(*result.Value)-1]
	}

	d.SetId(*image.ID)
	if location := image.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("publisher", publisher)
	d.Set("offer", offer)
	d.Set("sku", sku)
	d.Set("version", image.Name)

	return nil
}
