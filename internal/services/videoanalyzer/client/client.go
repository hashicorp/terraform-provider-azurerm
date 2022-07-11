package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	VideoAnalyzersClient *videoanalyzer.VideoAnalyzerClient
}

func NewClient(o *common.ClientOptions) *Client {
	VideoAnalyzersClient := videoanalyzer.NewVideoAnalyzerClientWithBaseURI(o.ResourceManagerEndpoint)

	o.ConfigureClient(&VideoAnalyzersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		VideoAnalyzersClient: &VideoAnalyzersClient,
	}
}
