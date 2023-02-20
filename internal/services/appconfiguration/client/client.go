package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/deletedconfigurationstores"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
)

type Client struct {
	ConfigurationStoresClient        *configurationstores.ConfigurationStoresClient
	DeletedConfigurationStoresClient *deletedconfigurationstores.DeletedConfigurationStoresClient
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

func NewClient(o *common.ClientOptions) *Client {
	configurationStores := configurationstores.NewConfigurationStoresClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationStores.Client, o.ResourceManagerAuthorizer)

	deletedConfigurationStores := deletedconfigurationstores.NewDeletedConfigurationStoresClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&deletedConfigurationStores.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationStoresClient:        &configurationStores,
		DeletedConfigurationStoresClient: &deletedConfigurationStores,
		authorizerFunc:                   o.Authorizers.AuthorizerFunc,
		configureClientFunc:              o.ConfigureClient,
	}
}
