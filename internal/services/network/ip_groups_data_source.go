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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

// Find IDs and names of multiple IP Groups, filtered by name substring
func dataSourceIpGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {

	// Establish a client to handle i/o operations against the API
	client := meta.(*clients.Client).Network.IPGroupsClient

	// Create a context for the request and defer cancellation
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Get resource group name from data source
	resourceGroupName := d.Get("resource_group_name").(string)

	// Make the request to the API to download all IP groups in the resource group
	allGroups, err := client.ListByResourceGroup(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing IP groups: %+v", err)
	}

	// Establish lists of strings to append to, set equal to empty set to start
	// If no IP groups are found, an empty set will be returned
	names := []string{}
	ids := []string{}

	// Filter IDs list by substring
	for _, ipGroup := range allGroups.Values() {
		if ipGroup.Name != nil && strings.Contains(*ipGroup.Name, d.Get("name").(string)) {
			names = append(names, *ipGroup.Name)
			ids = append(ids, *ipGroup.ID)
		}
	}

	// Sort lists of strings alphabetically
	slices.Sort(names)
	slices.Sort(ids)

	// Set resource ID, required for Terraform state
	// Since this is a multi-resource data source, we need to create a unique ID
	// Using the ID of the resource group
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := commonids.NewResourceGroupID(subscriptionId, resourceGroupName)
	d.SetId(id.ID())

	// Set names
	err = d.Set("names", names)
	if err != nil {
		return fmt.Errorf("error setting names: %+v", err)
	}

	// Set IDs
	err = d.Set("ids", ids)
	if err != nil {
		return fmt.Errorf("error setting ids: %+v", err)
	}

	// Return nil error
	return nil

}
