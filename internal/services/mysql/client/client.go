// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	servers_v2017_12_01 "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01"
	servers_v2020_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01"
	flexibleServers_v2021_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/azureadadministrators"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FlexibleServers  *flexibleServers_v2021_05_01.Client
	MySqlClient      *servers_v2017_12_01.Client
	ServerKeysClient *servers_v2020_01_01.Client

	// TODO: port over to using the Meta Client (which involves bumping the API Version)
	AzureADAdministratorsClient *azureadadministrators.AzureADAdministratorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	flexibleServersMetaClient, err := flexibleServers_v2021_05_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Flexible Servers client: %+v", err)
	}

	mySqlMetaClient, err := servers_v2017_12_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building MySql client: %+v", err)
	}

	serverKeysClient, err := servers_v2020_01_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building MySql client: %+v", err)
	}

	azureADAdministratorsClient, err := azureadadministrators.NewAzureADAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure AD Administrators client: %+v", err)
	}
	o.Configure(azureADAdministratorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FlexibleServers:  flexibleServersMetaClient,
		MySqlClient:      mySqlMetaClient,
		ServerKeysClient: serverKeysClient,

		// TODO: switch to using the Meta Clients
		AzureADAdministratorsClient: azureADAdministratorsClient,
	}, nil
}
