// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
)

type cacheItem struct {
	ID      managedhsms.ManagedHSMId
	BaseURI string
}

type localCache struct {
	mux        sync.Locker
	nameToItem map[string]cacheItem
}

var defaultCache = &localCache{
	mux:        &sync.Mutex{},
	nameToItem: map[string]cacheItem{},
}

func cacheKey(name string) string {
	return strings.ToLower(name)
}

func AddToCache(id managedhsms.ManagedHSMId, baseURI string) {
	defaultCache.add(id, baseURI)
}

func (l *localCache) add(id managedhsms.ManagedHSMId, baseURI string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.nameToItem[cacheKey(id.ManagedHSMName)] = cacheItem{
		ID:      id,
		BaseURI: baseURI,
	}
}

func (l *localCache) get(name string) (cacheItem, bool) {
	l.mux.Lock()
	defer l.mux.Unlock()

	item, ok := l.nameToItem[cacheKey(name)]
	return item, ok
}

func RemoveFromCache(name string) {
	defaultCache.remove(name)
}

func (l *localCache) remove(name string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	delete(l.nameToItem, cacheKey(name))
}

func (c *Client) BaseUriForManagedHSM(ctx context.Context, id managedhsms.ManagedHSMId) (*string, error) {
	item, ok := defaultCache.get(id.ManagedHSMName)
	if ok {
		return &item.BaseURI, nil
	}

	resp, err := c.ManagedHsmClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("managedHSM %s was not found", id)
		}
		return nil, fmt.Errorf("retrieving managedHSM %s: %+v", id, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.HsmUri != nil {
			vaultUri = *model.Properties.HsmUri
		}
	}
	if vaultUri == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", id)
	}

	defaultCache.add(id, vaultUri)
	return &vaultUri, nil
}

func (c *Client) ManagedHSMIDFromBaseUri(ctx context.Context, subscriptionId commonids.SubscriptionId, uri string) (*managedhsms.ManagedHSMId, error) {
	name, err := parseNameFromBaseUrl(uri)
	if err != nil {
		return nil, err
	}

	item, ok := defaultCache.get(*name)
	if ok {
		return &item.ID, nil
	}
	// fetch all managedhsms
	opts := managedhsms.DefaultListBySubscriptionOperationOptions()
	results, err := c.ManagedHsmClient.ListBySubscriptionComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing the managed HSM within %s: %+v", subscriptionId, err)
	}
	for _, item := range results.Items {
		if item.Id == nil || item.Properties.HsmUri == nil {
			continue
		}

		// Populate the managed HSM into the cache
		managedHSMID, err := managedhsms.ParseManagedHSMIDInsensitively(*item.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a managed HSM ID: %+v", *item.Id, err)
		}
		hsmUri := *item.Properties.HsmUri
		defaultCache.add(*managedHSMID, hsmUri)
	}

	// Now that the cache has been repopulated, check if we have the managed HSM or not
	if v, ok := defaultCache.get(*name); ok {
		return &v.ID, nil
	}
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) ManagedHSMExists(ctx context.Context, id managedhsms.ManagedHSMId) (bool, error) {
	_, ok := defaultCache.get(id.ManagedHSMName)
	if ok {
		return true, nil
	}

	resp, err := c.ManagedHsmClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving managedHSM %s: %+v", id, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.HsmUri != nil {
			vaultUri = *model.Properties.HsmUri
		}
	}
	if vaultUri == "" {
		return false, fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", id)
	}

	defaultCache.add(id, vaultUri)
	return true, nil
}

func parseNameFromBaseUrl(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	// https://the-hsm.managedhsm.azure.net
	// https://the-hsm.managedhsm.microsoftazure.de
	// https://the-hsm.managedhsm.usgovcloudapi.net
	// https://the-hsm.managedhsm.cloudapi.microsoft
	// https://the-hsm.managedhsm.azure.cn

	segments := strings.Split(uri.Host, ".")
	if len(segments) < 3 || segments[1] != "managedhsm" {
		return nil, fmt.Errorf("expected a URI in the format `the-managedhsm-name.managedhsm.**` but got %q", uri.Host)
	}
	return &segments[0], nil
}
