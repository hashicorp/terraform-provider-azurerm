// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
)

func (c *Client) FindPrivateDnsZoneId(ctx context.Context, resourceGroupsClient *resourcegroups.ResourceGroupsClient, subscriptionId commonids.SubscriptionId, name string) (*privatezones.PrivateDnsZoneId, error) {
	opts := privatezones.DefaultListOperationOptions()
	results, err := c.PrivateZonesClient.ListComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing the Private DNS Zones within %s: %+v", subscriptionId, err)
	}

	for _, item := range results.Items {
		if item.Id == nil {
			continue
		}

		itemId := *item.Id
		parsed, err := privatezones.ParsePrivateDnsZoneIDInsensitively(itemId)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Private DNS Zone ID: %+v", itemId, err)
		}

		if parsed.PrivateDnsZoneName != name {
			continue
		}

		// however the Resource Group name isn't necessarily cased correctly, so now that we've found the Resource
		// Group name let's pull the canonical casing from the Resource Groups API
		resourceGroupId := commonids.NewResourceGroupID(parsed.SubscriptionId, parsed.ResourceGroupName)
		resp, err := resourceGroupsClient.Get(ctx, resourceGroupId)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", resourceGroupId, err)
		}
		if model := resp.Model; model != nil && model.Name != nil {
			parsed.ResourceGroupName = *model.Name
		}
		return parsed, nil
	}

	return nil, fmt.Errorf("No Private DNS Zones found with name: %q", name)
}
