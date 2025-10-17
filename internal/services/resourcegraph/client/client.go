// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/graphqueries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/graphquery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ResourceGraphQueryClient   *graphquery.GraphQueryClient
	ResourceGraphQueriesClient *graphqueries.GraphqueriesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	graphQueryClient, err := graphquery.NewGraphQueryClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(graphQueryClient.Client, o.Authorizers.ResourceManager)

	graphQueriesClient, err := graphqueries.NewGraphqueriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(graphQueriesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ResourceGraphQueryClient:   graphQueryClient,
		ResourceGraphQueriesClient: graphQueriesClient,
	}, nil
}
