// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	searchDataPlane "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	dataplaneClient "github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AdminKeysClient                       *adminkeys.AdminKeysClient
	SearchDataPlaneClient                 *searchDataPlane.Client
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

	searchDataPlaneClient, err := searchDataPlane.NewClient(func(c *dataplaneClient.Client) {
		o.Configure(c.Client, o.Authorizers.Search)
	})
	if err != nil {
		return nil, fmt.Errorf("building data-plane Search client: %+v", err)
	}

	return &Client{
		AdminKeysClient:                       adminKeysClient,
		SearchDataPlaneClient:                 searchDataPlaneClient,
		QueryKeysClient:                       queryKeysClient,
		ServicesClient:                        servicesClient,
		SearchSharedPrivateLinkResourceClient: searchSharedPrivateLinkResourceClient,
	}, nil
}
