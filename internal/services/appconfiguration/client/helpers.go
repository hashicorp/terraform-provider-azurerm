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
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
)

var (
	configurationStoreCache = map[string]ConfigurationStoreDetails{}
	keysmith                = &sync.RWMutex{}
	lock                    = map[string]*sync.RWMutex{}
)

type ConfigurationStoreDetails struct {
	configurationStoreId string
	dataPlaneEndpoint    string
}

func (c *Client) AddToCache(configurationStoreId configurationstores.ConfigurationStoreId, dataPlaneEndpoint string) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	configurationStoreCache[cacheKey] = ConfigurationStoreDetails{
		configurationStoreId: configurationStoreId.ID(),
		dataPlaneEndpoint:    dataPlaneEndpoint,
	}
	keysmith.Unlock()
}

func (c *Client) ConfigurationStoreIDFromEndpoint(ctx context.Context, subscriptionId commonids.SubscriptionId, configurationStoreEndpoint, domainSuffix string) (*string, error) {
	configurationStoreName, err := c.parseNameFromEndpoint(configurationStoreEndpoint, domainSuffix)
	if err != nil {
		return nil, err
	}

	cacheKey := c.cacheKeyForConfigurationStore(*configurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	// first check the cache
	if v, ok := configurationStoreCache[cacheKey]; ok {
		return &v.configurationStoreId, nil
	}

	// If it's not present, populate the entire cache
	configurationStores, err := c.ConfigurationStoresClient.ListComplete(ctx, subscriptionId)
	if err != nil {
		return nil, fmt.Errorf("retrieving the list of Configuration Stores in %s: %+v", subscriptionId, err)
	}
	for _, item := range configurationStores.Items {
		if item.Id == nil || item.Properties == nil || item.Properties.Endpoint == nil {
			continue
		}

		itemId := *item.Id
		endpointUri := *item.Properties.Endpoint
		configurationStoreId, err := configurationstores.ParseConfigurationStoreIDInsensitively(itemId)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", itemId, err)
		}

		c.AddToCache(*configurationStoreId, endpointUri)
	}

	// finally try and pull this from the cache
	if v, ok := configurationStoreCache[cacheKey]; ok {
		return &v.configurationStoreId, nil
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) EndpointForConfigurationStore(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId) (*string, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := configurationStoreCache[cacheKey]; ok {
		return &v.dataPlaneEndpoint, nil
	}

	resp, err := c.ConfigurationStoresClient.Get(ctx, configurationStoreId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s:%+v", configurationStoreId, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
		return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", configurationStoreId)
	}

	c.AddToCache(configurationStoreId, *resp.Model.Properties.Endpoint)

	return resp.Model.Properties.Endpoint, nil
}

func (c *Client) Exists(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId) (bool, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if _, ok := configurationStoreCache[cacheKey]; ok {
		return true, nil
	}

	resp, err := c.ConfigurationStoresClient.Get(ctx, configurationStoreId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving %s: %+v", configurationStoreId, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
		return false, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", configurationStoreId)
	}

	c.AddToCache(configurationStoreId, *resp.Model.Properties.Endpoint)

	return true, nil
}

func (c *Client) RemoveFromCache(configurationStoreId configurationstores.ConfigurationStoreId) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	delete(configurationStoreCache, cacheKey)
	lock[cacheKey].Unlock()
}

func (c *Client) cacheKeyForConfigurationStore(name string) string {
	return strings.ToLower(name)
}

func (c *Client) parseNameFromEndpoint(input, domainSuffix string) (*string, error) {
	uri, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, err
	}

	// https://the-appconfiguration.azconfig.io
	// https://the-appconfiguration.azconfig.azure.cn
	// https://the-appconfiguration.azconfig.azure.us
	if !strings.HasSuffix(uri.Host, domainSuffix) {
		return nil, fmt.Errorf("expected a URI in the format `https://somename.%s` but got %q", domainSuffix, uri.Host)
	}

	segments := strings.Split(uri.Host, ".")
	return &segments[0], nil
}
