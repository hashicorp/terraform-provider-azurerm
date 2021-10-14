package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/videoanalyzer/mgmt/2021-05-01-preview/videoanalyzer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	VideoAnalyzersClient *videoanalyzer.VideoAnalyzersClient
	EdgeModulesClient    *videoanalyzer.EdgeModulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	VideoAnalyzersClient := videoanalyzer.NewVideoAnalyzersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	EdgeModulesClient := videoanalyzer.NewEdgeModulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)

	o.ConfigureClient(&VideoAnalyzersClient.Client, o.ResourceManagerAuthorizer)
	o.ConfigureClient(&EdgeModulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		VideoAnalyzersClient: &VideoAnalyzersClient,
		EdgeModulesClient:    &EdgeModulesClient,
	}
}
