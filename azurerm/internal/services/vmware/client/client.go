package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/sdk/authorizations"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/sdk/clusters"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/sdk/privateclouds"
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
