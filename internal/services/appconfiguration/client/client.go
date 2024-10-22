// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/deletedconfigurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/operations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/azuresdkhacks"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

type Client struct {
	ConfigurationStoresClient        *configurationstores.ConfigurationStoresClient
	DeletedConfigurationStoresClient *deletedconfigurationstores.DeletedConfigurationStoresClient
	OperationsClient                 *operations.OperationsClient
	ReplicasClient                   *replicas.ReplicasClient
	authorizerFunc                   common.ApiAuthorizerFunc
	configureClientFunc              func(c *autorest.Client, authorizer autorest.Authorizer)
}

func (c *Client) DataPlaneClientWithEndpoint(configurationStoreEndpoint string) (*appconfiguration.BaseClient, error) {
	api := environments.NewApiEndpoint("AppConfiguration", configurationStoreEndpoint, nil)
	appConfigAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", configurationStoreEndpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", configurationStoreEndpoint)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(appConfigAuth))

	return &client, nil
}

func (c *Client) LinkWorkaroundDataPlaneClientWithEndpoint(configurationStoreEndpoint string) (*azuresdkhacks.DataPlaneClient, error) {
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

	replicasClient, err := replicas.NewReplicasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeletedConfigurationStores client: %+v", err)
	}
	o.Configure(replicasClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConfigurationStoresClient:        configurationStores,
		DeletedConfigurationStoresClient: deletedConfigurationStores,
		OperationsClient:                 operationsClient,
		ReplicasClient:                   replicasClient,
		authorizerFunc:                   o.Authorizers.AuthorizerFunc,
		configureClientFunc:              o.ConfigureClient,
	}, nil
}
