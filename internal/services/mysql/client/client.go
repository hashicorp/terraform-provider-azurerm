// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	flexibleServers_v2023_12_30 "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FlexibleServers      *flexibleServers_v2023_12_30.Client
	FlexibleServerRename *servers.ServersClient

	// TODO: port over to using the Meta Client (which involves bumping the API Version)
	AzureADAdministratorsClient *azureadadministrators.AzureADAdministratorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	flexibleServersMetaClient, err := flexibleServers_v2023_12_30.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Flexible Servers client: %+v", err)
	}

	azureADAdministratorsClient, err := azureadadministrators.NewAzureADAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure AD Administrators client: %+v", err)
	}
	o.Configure(azureADAdministratorsClient.Client, o.Authorizers.ResourceManager)

	// Use the ARM-registered non-preview API for rename; keep the existing version for ordinary server operations.
	// The request uses only the unchanged location, createMode and sourceServerResourceId fields.
	renameClient, err := resourcemanager.NewClient(o.Environment.ResourceManager, "mysql-server-rename", "2025-12-30")
	if err != nil {
		return nil, fmt.Errorf("building Flexible Server Rename client: %+v", err)
	}
	o.Configure(renameClient, o.Authorizers.ResourceManager)

	return &Client{
		FlexibleServers:      flexibleServersMetaClient,
		FlexibleServerRename: &servers.ServersClient{Client: renameClient},

		// TODO: switch to using the Meta Clients
		AzureADAdministratorsClient: azureADAdministratorsClient,
	}, nil
}
