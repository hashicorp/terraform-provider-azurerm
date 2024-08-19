// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceIpGroups() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIpGroupsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceIpGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(subscriptionId, d.Get("resource_group_name").(string))

	resp, err := client.ListByResourceGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("listing IP groups: %+v", err)
	}

	names := []string{}
	ids := []string{}

	if model := resp.Model; model != nil {
		for _, group := range *model {
			if group.Name != nil && strings.Contains(*group.Name, d.Get("name").(string)) {
				names = append(names, *group.Name)
				ids = append(ids, *group.Id)
			}
		}
	}

	slices.Sort(names)
	slices.Sort(ids)

	d.SetId(id.ID())

	err = d.Set("names", names)
	if err != nil {
		return fmt.Errorf("error setting names: %+v", err)
	}

	err = d.Set("ids", ids)
	if err != nil {
		return fmt.Errorf("error setting ids: %+v", err)
	}

	return nil
}
