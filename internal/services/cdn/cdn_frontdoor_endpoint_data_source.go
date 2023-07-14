// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorEndpointRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorEndpointName,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceCdnFrontDoorEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AfdEndpointName)
	d.Set("profile_name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.AFDEndpointProperties; props != nil {
		d.Set("enabled", flattenEnabledBool(props.EnabledState))
		d.Set("host_name", props.HostName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
