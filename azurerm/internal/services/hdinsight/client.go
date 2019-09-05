package hdinsight

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationsClient   *hdinsight.ApplicationsClient
	ClustersClient       *hdinsight.ClustersClient
	ConfigurationsClient *hdinsight.ConfigurationsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ApplicationsClient := hdinsight.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationsClient.Client, o.ResourceManagerAuthorizer)

	ClustersClient := hdinsight.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	ConfigurationsClient := hdinsight.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationsClient:   &ApplicationsClient,
		ClustersClient:       &ClustersClient,
		ConfigurationsClient: &ConfigurationsClient,
	}
}
