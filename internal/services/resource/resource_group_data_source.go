// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourceGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name":     commonschema.ResourceGroupNameForDataSource(),
			"location": commonschema.LocationComputed(),
			"tags":     commonschema.TagsDataSource(),
			"managed_by": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%q was not found", id)
		}

		return err
	}

	d.SetId(id.ID())

	d.Set("name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("managed_by", pointer.From(model.ManagedBy))

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
