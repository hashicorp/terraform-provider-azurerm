package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AuthorizationClient *authorizations.AuthorizationsClient
	ClusterClient       *clusters.ClustersClient
	PrivateCloudClient  *privateclouds.PrivateCloudsClient
}

func NewClient(o *common.ClientOptions) *Client {
	authorizationClient := authorizations.NewAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&authorizationClient.Client, o.ResourceManagerAuthorizer)

	clusterClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clusterClient.Client, o.ResourceManagerAuthorizer)

	privateCloudClient := privateclouds.NewPrivateCloudsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&privateCloudClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AuthorizationClient: &authorizationClient,
		ClusterClient:       &clusterClient,
		PrivateCloudClient:  &privateCloudClient,
	}
}
