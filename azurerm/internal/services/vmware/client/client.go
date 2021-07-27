package client

import (
	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/sdk/clusters"
)

type Client struct {
	AuthorizationClient *avs.AuthorizationsClient
	ClusterClient       *clusters.ClustersClient
	PrivateCloudClient  *avs.PrivateCloudsClient
}

func NewClient(o *common.ClientOptions) *Client {
	authorizationClient := avs.NewAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&authorizationClient.Client, o.ResourceManagerAuthorizer)

	clusterClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clusterClient.Client, o.ResourceManagerAuthorizer)

	privateCloudClient := avs.NewPrivateCloudsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&privateCloudClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AuthorizationClient: &authorizationClient,
		ClusterClient:       &clusterClient,
		PrivateCloudClient:  &privateCloudClient,
	}
}
