// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AdminKeysClient                       *adminkeys.AdminKeysClient
	QueryKeysClient                       *querykeys.QueryKeysClient
	ServicesClient                        *services.ServicesClient
	SearchSharedPrivateLinkResourceClient *sharedprivatelinkresources.SharedPrivateLinkResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	adminKeysClient, err := adminkeys.NewAdminKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AdminKeys Client: %+v", err)
	}
	o.Configure(adminKeysClient.Client, o.Authorizers.ResourceManager)

	queryKeysClient, err := querykeys.NewQueryKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building QueryKeys Client: %+v", err)
	}
	o.Configure(queryKeysClient.Client, o.Authorizers.ResourceManager)

	servicesClient, err := services.NewServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Services Client: %+v", err)
	}
	o.Configure(servicesClient.Client, o.Authorizers.ResourceManager)

	searchSharedPrivateLinkResourceClient, err := sharedprivatelinkresources.NewSharedPrivateLinkResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SharedPrivateLinkResource Client: %+v", err)
	}
	o.Configure(searchSharedPrivateLinkResourceClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AdminKeysClient:                       adminKeysClient,
		QueryKeysClient:                       queryKeysClient,
		ServicesClient:                        servicesClient,
		SearchSharedPrivateLinkResourceClient: searchSharedPrivateLinkResourceClient,
	}, nil
}
