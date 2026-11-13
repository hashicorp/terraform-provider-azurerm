// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
)

var (
	serviceTagsCache = map[string]*virtualnetworks.ServiceTagsListResult{}
	serviceTagsMu    sync.Mutex
)

func serviceTagsCacheKey(locationName string) string {
	return strings.ToLower(locationName)
}

func (c *Client) CachedServiceTagsList(ctx context.Context, locationId virtualnetworks.LocationId) (virtualnetworks.ServiceTagsListOperationResponse, error) {
	cacheKey := serviceTagsCacheKey(locationId.LocationName)

	serviceTagsMu.Lock()
	defer serviceTagsMu.Unlock()

	if cached, ok := serviceTagsCache[cacheKey]; ok {
		return virtualnetworks.ServiceTagsListOperationResponse{
			Model: cached,
		}, nil
	}

	resp, err := c.VirtualNetworks.ServiceTagsList(ctx, locationId)
	if err != nil {
		return resp, fmt.Errorf("listing network service tags for %s: %+v", locationId, err)
	}

	if resp.Model != nil {
		serviceTagsCache[cacheKey] = resp.Model
	}

	return resp, nil
}
