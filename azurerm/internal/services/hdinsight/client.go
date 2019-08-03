package hdinsight

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationsClient   hdinsight.ApplicationsClient
	ClustersClient       hdinsight.ClustersClient
	ConfigurationsClient hdinsight.ConfigurationsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ApplicationsClient = hdinsight.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApplicationsClient.Client, o.ResourceManagerAuthorizer)

	c.ClustersClient = hdinsight.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ClustersClient.Client, o.ResourceManagerAuthorizer)

	c.ConfigurationsClient = hdinsight.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
