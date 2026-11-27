// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	resourceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
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

func dataSourceResourceGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := resourceClient.From(meta).ResourceGroupsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := commonids.NewResourceGroupID(account.From(meta).SubscriptionId, d.Get("name").(string))
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
