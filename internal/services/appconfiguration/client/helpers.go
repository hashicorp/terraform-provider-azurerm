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
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	ConfigurationStoreCache = map[string]ConfigurationStoreDetails{}
	keysmith                = &sync.RWMutex{}
	lock                    = map[string]*sync.RWMutex{}
)

type ConfigurationStoreDetails struct {
	configurationStoreId string
	dataPlaneEndpoint    string
}

func (c Client) AddToCache(configurationStoreId configurationstores.ConfigurationStoreId, dataPlaneEndpoint string) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	ConfigurationStoreCache[cacheKey] = ConfigurationStoreDetails{
		configurationStoreId: configurationStoreId.ID(),
		dataPlaneEndpoint:    dataPlaneEndpoint,
	}
	keysmith.Unlock()
}

func (c Client) ConfigurationStoreIDFromEndpoint(ctx context.Context, resourcesClient *resourcesClient.Client, configurationStoreEndpoint string) (*string, error) {
	configurationStoreName, err := c.parseNameFromEndpoint(configurationStoreEndpoint)
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

	if v, ok := ConfigurationStoreCache[cacheKey]; ok {
		return &v.configurationStoreId, nil
	}

	filter := fmt.Sprintf("resourceType eq 'Microsoft.AppConfiguration/configurationStores' and name eq '%s'", *configurationStoreName)
	result, err := resourcesClient.ResourcesClient.List(ctx, filter, "", utils.Int32(5))
	if err != nil {
		return nil, fmt.Errorf("listing resources matching %q: %+v", filter, err)
	}

	for result.NotDone() {
		for _, v := range result.Values() {
			if v.ID == nil {
				continue
			}

			id, err := configurationstores.ParseConfigurationStoreIDInsensitively(*v.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
			}
			if !strings.EqualFold(id.ConfigurationStoreName, *configurationStoreName) {
				continue
			}

			resp, err := c.ConfigurationStoresClient.Get(ctx, *id)
			if err != nil {
				return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
				return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", *id)
			}

			c.AddToCache(*id, *resp.Model.Properties.Endpoint)

			return utils.String(id.ID()), nil
		}

		if err := result.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("iterating over results: %+v", err)
		}
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c Client) EndpointForConfigurationStore(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId) (*string, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := ConfigurationStoreCache[cacheKey]; ok {
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

func (c Client) Exists(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId) (bool, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if _, ok := ConfigurationStoreCache[cacheKey]; ok {
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

func (c Client) RemoveFromCache(configurationStoreId configurationstores.ConfigurationStoreId) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	delete(ConfigurationStoreCache, cacheKey)
	lock[cacheKey].Unlock()
}

func (c Client) cacheKeyForConfigurationStore(name string) string {
	return strings.ToLower(name)
}

func (c Client) parseNameFromEndpoint(input string) (*string, error) {
	uri, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, err
	}

	// https://the-appconfiguration.azconfig.io

	segments := strings.Split(uri.Host, ".")
	if len(segments) < 3 || segments[1] != "azconfig" || segments[2] != "io" {
		return nil, fmt.Errorf("expected a URI in the format `https://the-appconfiguration.azconfig.io` but got %q", uri.Host)
	}
	return &segments[0], nil
}
