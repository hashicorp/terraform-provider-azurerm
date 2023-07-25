// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/operations"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/azuresdkhacks"
	"github.com/tombuildsstuff/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

type Client struct {
	ConfigurationStoresClient        *configurationstores.ConfigurationStoresClient
	DeletedConfigurationStoresClient *deletedconfigurationstores.DeletedConfigurationStoresClient
	OperationsClient                 *operations.OperationsClient
	authorizerFunc                   common.ApiAuthorizerFunc
	configureClientFunc              func(c *autorest.Client, authorizer autorest.Authorizer)
}

func (c Client) DataPlaneClientWithEndpoint(configurationStoreEndpoint string) (*appconfiguration.BaseClient, error) {
	api := environments.NewApiEndpoint("AppConfiguration", configurationStoreEndpoint, nil)
	appConfigAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", configurationStoreEndpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", configurationStoreEndpoint)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(appConfigAuth))

	return &client, nil
}

func (c Client) LinkWorkaroundDataPlaneClientWithEndpoint(configurationStoreEndpoint string) (*azuresdkhacks.DataPlaneClient, error) {
	api := environments.NewApiEndpoint("AppConfiguration", configurationStoreEndpoint, nil)
	appConfigAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", configurationStoreEndpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", configurationStoreEndpoint)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(appConfigAuth))
	workaroundClient := azuresdkhacks.NewDataPlaneClient(client)

	return &workaroundClient, nil
}

func (c Client) DataPlaneClient(ctx context.Context, configurationStoreId string) (*appconfiguration.BaseClient, error) {
	appConfigId, err := configurationstores.ParseConfigurationStoreID(configurationStoreId)
	if err != nil {
		return nil, err
	}

	// TODO: caching all of this
	appConfig, err := c.ConfigurationStoresClient.Get(ctx, *appConfigId)
	if err != nil {
		if response.WasNotFound(appConfig.HttpResponse) {
			return nil, nil
		}

		return nil, err
	}

	if appConfig.Model == nil || appConfig.Model.Properties == nil || appConfig.Model.Properties.Endpoint == nil {
		return nil, fmt.Errorf("endpoint was nil")
	}

	endpoint := *appConfig.Model.Properties.Endpoint

	api := environments.NewApiEndpoint("AppConfiguration", endpoint, nil)
	appConfigAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", endpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", endpoint)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(appConfigAuth))

	return &client, nil
}

func (c Client) LinkWorkaroundDataPlaneClient(ctx context.Context, configurationStoreId string) (*azuresdkhacks.DataPlaneClient, error) {
	appConfigId, err := configurationstores.ParseConfigurationStoreID(configurationStoreId)
	if err != nil {
		return nil, err
	}

	// TODO: caching all of this
	appConfig, err := c.ConfigurationStoresClient.Get(ctx, *appConfigId)
	if err != nil {
		if response.WasNotFound(appConfig.HttpResponse) {
			return nil, nil
		}

		return nil, err
	}

	if appConfig.Model == nil || appConfig.Model.Properties == nil || appConfig.Model.Properties.Endpoint == nil {
		return nil, fmt.Errorf("endpoint was nil")
	}

	api := environments.NewApiEndpoint("AppConfiguration", *appConfig.Model.Properties.Endpoint, nil)
	appConfigAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", *appConfig.Model.Properties.Endpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", *appConfig.Model.Properties.Endpoint)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(appConfigAuth))
	workaroundClient := azuresdkhacks.NewDataPlaneClient(client)

	return &workaroundClient, nil
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	configurationStores, err := configurationstores.NewConfigurationStoresClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationStores client: %+v", err)
	}
	o.Configure(configurationStores.Client, o.Authorizers.ResourceManager)

	deletedConfigurationStores, err := deletedconfigurationstores.NewDeletedConfigurationStoresClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeletedConfigurationStores client: %+v", err)
	}
	o.Configure(deletedConfigurationStores.Client, o.Authorizers.ResourceManager)

	operationsClient, err := operations.NewOperationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Operations client: %+v", err)
	}
	o.Configure(operationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConfigurationStoresClient:        configurationStores,
		DeletedConfigurationStoresClient: deletedConfigurationStores,
		OperationsClient:                 operationsClient,
		authorizerFunc:                   o.Authorizers.AuthorizerFunc,
		configureClientFunc:              o.ConfigureClient,
	}, nil
}
