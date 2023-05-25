package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApplicationsClient   *applications.ApplicationsClient
	ClustersClient       *clusters.ClustersClient
	ConfigurationsClient *configurations.ConfigurationsClient
	ExtensionsClient     *extensions.ExtensionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// due to a bug in the HDInsight API we can't reuse client with the same x-ms-correlation-request-id for multiple updates
	opts := *o
	opts.DisableCorrelationRequestID = true

	applicationsClient, err := applications.NewApplicationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ApplicationsClient Client: %+v", err)
	}
	opts.Configure(applicationsClient.Client, opts.Authorizers.ResourceManager)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ClustersClient Client: %+v", err)
	}
	opts.Configure(clustersClient.Client, opts.Authorizers.ResourceManager)

	configurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationsClient Client: %+v", err)
	}
	opts.Configure(configurationsClient.Client, opts.Authorizers.ResourceManager)

	extensionsClient, err := extensions.NewExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExtensionsClient Client: %+v", err)
	}
	opts.Configure(extensionsClient.Client, opts.Authorizers.ResourceManager)

	return &Client{
		ApplicationsClient:   applicationsClient,
		ClustersClient:       clustersClient,
		ConfigurationsClient: configurationsClient,
		ExtensionsClient:     extensionsClient,
	}, nil
}
