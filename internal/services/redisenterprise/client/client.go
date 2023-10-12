// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2023-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2023-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Client         *redisenterprise.RedisEnterpriseClient
	DatabaseClient *databases.DatabasesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := redisenterprise.NewRedisEnterpriseClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RedisEnterprise client: %+v", err)
	}
	o.Configure(client.Client, o.Authorizers.ResourceManager)

	databaseClient, err := databases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases client: %+v", err)
	}
	o.Configure(databaseClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		Client:         client,
		DatabaseClient: databaseClient,
	}, nil
}
