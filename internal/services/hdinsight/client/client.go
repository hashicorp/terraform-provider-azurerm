package client

import (
	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApplicationsClient   *hdinsight.ApplicationsClient
	ClustersClient       *hdinsight.ClustersClient
	ConfigurationsClient *hdinsight.ConfigurationsClient
	ExtensionsClient     *hdinsight.ExtensionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	// due to a bug in the HDInsight API we can't reuse client with the same x-ms-correlation-request-id for multiple updates
	opts := *o
	opts.DisableCorrelationRequestID = true

	ApplicationsClient := hdinsight.NewApplicationsClientWithBaseURI(opts.ResourceManagerEndpoint, opts.SubscriptionId)
	opts.ConfigureClient(&ApplicationsClient.Client, opts.ResourceManagerAuthorizer)

	ClustersClient := hdinsight.NewClustersClientWithBaseURI(opts.ResourceManagerEndpoint, opts.SubscriptionId)
	opts.ConfigureClient(&ClustersClient.Client, opts.ResourceManagerAuthorizer)

	ConfigurationsClient := hdinsight.NewConfigurationsClientWithBaseURI(opts.ResourceManagerEndpoint, opts.SubscriptionId)
	opts.ConfigureClient(&ConfigurationsClient.Client, opts.ResourceManagerAuthorizer)

	ExtensionsClient := hdinsight.NewExtensionsClientWithBaseURI(opts.ResourceManagerEndpoint, opts.SubscriptionId)
	opts.ConfigureClient(&ExtensionsClient.Client, opts.ResourceManagerAuthorizer)

	c := &Client{
		ApplicationsClient:   &ApplicationsClient,
		ClustersClient:       &ClustersClient,
		ConfigurationsClient: &ConfigurationsClient,
		ExtensionsClient:     &ExtensionsClient,
	}

	return c
}
