package client

import (
	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationsClient   *hdinsight.ApplicationsClient
	ClustersClient       *hdinsight.ClustersClient
	ConfigurationsClient *hdinsight.ConfigurationsClient
	ExtensionsClient     *hdinsight.ExtensionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ApplicationsClient := hdinsight.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationsClient.Client, o.ResourceManagerAuthorizer)

	ClustersClient := hdinsight.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	ConfigurationsClient := hdinsight.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ExtensionsClient := hdinsight.NewExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExtensionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationsClient:   &ApplicationsClient,
		ClustersClient:       &ClustersClient,
		ConfigurationsClient: &ConfigurationsClient,
		ExtensionsClient:     &ExtensionsClient,
	}
}
