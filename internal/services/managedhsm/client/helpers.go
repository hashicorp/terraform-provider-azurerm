// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

var (
	cache    = map[string]managedHSMDetail{}
	keysmith = &sync.RWMutex{}
	lock     = map[string]*sync.RWMutex{}
)

type managedHSMDetail struct {
	managedHSMId     managedhsms.ManagedHSMId
	dataPlaneBaseUri string
}

func (c *Client) AddToCache(managedHsmId managedhsms.ManagedHSMId, dataPlaneUri string) {
	cacheKey := c.cacheKeyForManagedHSM(managedHsmId.ManagedHSMName)
	keysmith.Lock()
	cache[cacheKey] = managedHSMDetail{
		managedHSMId:     managedHsmId,
		dataPlaneBaseUri: dataPlaneUri,
	}
	keysmith.Unlock()
}

func (c *Client) BaseUriForManagedHSM(ctx context.Context, managedHsmId managedhsms.ManagedHSMId) (*string, error) {
	cacheKey := c.cacheKeyForManagedHSM(managedHsmId.ManagedHSMName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := cache[cacheKey]; ok {
		return &v.dataPlaneBaseUri, nil
	}

	resp, err := c.ManagedHsmClient.Get(ctx, managedHsmId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", managedHsmId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", managedHsmId, err)
	}

	dataPlaneUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.HsmUri != nil {
			dataPlaneUri = *model.Properties.HsmUri
		}
	}
	if dataPlaneUri == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.HsmUri` was nil", managedHsmId)
	}

	c.AddToCache(managedHsmId, dataPlaneUri)
	return &dataPlaneUri, nil
}

func (c *Client) Purge(managedHSMId managedhsms.ManagedHSMId) {
	cacheKey := c.cacheKeyForManagedHSM(managedHSMId.ManagedHSMName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	delete(cache, cacheKey)
	lock[cacheKey].Unlock()
}

func (c *Client) ManagedHSMIDFromBaseUrl(ctx context.Context, subscriptionId commonids.SubscriptionId, managedHsmBaseUrl string, domainSuffix *string) (*managedhsms.ManagedHSMId, error) {
	endpoint, err := parse.ManagedHSMEndpoint(managedHsmBaseUrl, domainSuffix)
	if err != nil {
		return nil, err
	}

	cacheKey := c.cacheKeyForManagedHSM(endpoint.ManagedHSMName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	// Check the cache to determine if we have an entry for this Managed HSM
	if v, ok := cache[cacheKey]; ok {
		return &v.managedHSMId, nil
	}

	// Populate the cache if not found
	if err = c.populateCache(ctx, subscriptionId); err != nil {
		return nil, err
	}

	// Now that the cache has been repopulated, check if we have the Managed HSM or not
	if v, ok := cache[cacheKey]; ok {
		return &v.managedHSMId, nil
	}

	// We haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) cacheKeyForManagedHSM(name string) string {
	return strings.ToLower(name)
}

func (c *Client) populateCache(ctx context.Context, subscriptionId commonids.SubscriptionId) error {
	// Pull out the list of Managed HSMs available within the Subscription to re-populate the cache
	//
	// Whilst we've historically used the Resources API to query the single Managed HSM in question
	// this endpoint has caching related issues - and whilst the ResourceGraph API has been suggested
	// as an alternative that fixes this, we've seen similar caching issues there.
	// Therefore, we're falling back on querying all the Managed HSMs within the specified Subscription, which
	// comes from the `KeyVault` Resource Provider rather than the `Resources` Resource Provider - which
	// is an approach we've used previously, but now with better caching.
	//
	// Whilst querying ALL Managed HSMs within a Subscription IS excessive where only a single Managed HSM
	// is used - having the cache populated (one-time, per Provider launch) should alleviate problems
	// in Terraform Configurations defining a large number of Managed HSM related items.
	//
	// Note that we will only populate the cache on demand where HSM resources are actually being managed, to avoid
	// the overhead in configurations where HSM is not used.
	//
	// Finally, it's worth noting that we intentionally List ALL the Managed HSMs within a Subscription
	// to be able to cache ALL of them - prior to looking up the specific Managed HSM we're interested
	// in from the freshly populated cache.
	// This fixes an issue in the previous implementation where the Cache was being repeatedly semi-populated
	// until the specified Managed HSM was found, at which point we skipped populating the cache, which
	// affected both the `Resources` API implementation:
	// https://github.com/hashicorp/terraform-provider-azurerm/blob/3e88e5e74e12577d785f10298281b1b3c172254f/internal/services/keyvault/client/helpers.go#L133-L173
	// and the `ListBySubscription` endpoint:
	// https://github.com/hashicorp/terraform-provider-azurerm/blob/a5e728dc62e832e74d7bb0f40a79af0ae5a79e1e/azurerm/helpers/azure/key_vault.go#L42-L89

	opts := managedhsms.DefaultListBySubscriptionOperationOptions()
	results, err := c.ManagedHsmClient.ListBySubscriptionComplete(ctx, subscriptionId, opts)
	if err != nil {
		return fmt.Errorf("listing the Managed HSMs within %s: %+v", subscriptionId, err)
	}
	for _, item := range results.Items {
		if item.Id == nil || item.Properties.HsmUri == nil {
			continue
		}

		// Populate the Managed HSM into the cache
		managedHsm, err := managedhsms.ParseManagedHSMIDInsensitively(*item.Id)
		if err != nil {
			return fmt.Errorf("parsing %q as a Managed HSM ID: %+v", *item.Id, err)
		}
		dataPlaneUri := *item.Properties.HsmUri
		c.AddToCache(*managedHsm, dataPlaneUri)
	}

	return nil
}
