// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AuthorizationClient *authorizations.AuthorizationsClient
	ClusterClient       *clusters.ClustersClient
	PrivateCloudClient  *privateclouds.PrivateCloudsClient
	DataStoreClient     *datastores.DataStoresClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	authorizationClient, err := authorizations.NewAuthorizationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Authorization Client: %+v", err)
	}
	o.Configure(authorizationClient.Client, o.Authorizers.ResourceManager)

	clusterClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Cluster Client: %+v", err)
	}
	o.Configure(clusterClient.Client, o.Authorizers.ResourceManager)

	privateCloudClient, err := privateclouds.NewPrivateCloudsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Private Cloud Client: %+v", err)
	}
	o.Configure(privateCloudClient.Client, o.Authorizers.ResourceManager)

	dataStoresClient, err := datastores.NewDataStoresClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Stores Client: %+v", err)
	}
	o.Configure(dataStoresClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AuthorizationClient: authorizationClient,
		ClusterClient:       clusterClient,
		PrivateCloudClient:  privateCloudClient,
		DataStoreClient:     dataStoresClient,
	}, nil
}
