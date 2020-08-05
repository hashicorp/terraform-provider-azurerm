package client

import (
    "github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)
type Client struct {
    PrivateCloudClient    *avs.PrivateCloudsClient
    ClusterClient    *avs.ClustersClient
    HcxEnterpriseSiteClient    *avs.HcxEnterpriseSitesClient
    AuthorizationClient    *avs.AuthorizationsClient
}

func NewClient(o *common.ClientOptions) *Client {
    privateCloudClient := avs.NewPrivateCloudsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
    o.ConfigureClient(&privateCloudClient.Client, o.ResourceManagerAuthorizer)

    clusterClient := avs.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
    o.ConfigureClient(&clusterClient.Client, o.ResourceManagerAuthorizer)

    hcxEnterpriseSiteClient := avs.NewHcxEnterpriseSitesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
    o.ConfigureClient(&hcxEnterpriseSiteClient.Client, o.ResourceManagerAuthorizer)

    authorizationClient := avs.NewAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
    o.ConfigureClient(&authorizationClient.Client, o.ResourceManagerAuthorizer)

    return &Client{
        PrivateCloudClient:    &privateCloudClient,
        ClusterClient:    &clusterClient,
        HcxEnterpriseSiteClient:    &hcxEnterpriseSiteClient,
        AuthorizationClient:    &authorizationClient,
    }
}
