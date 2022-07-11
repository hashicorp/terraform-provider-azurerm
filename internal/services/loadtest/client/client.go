package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LoadTestsClient *loadtests.LoadTestsClient
}

func NewClient(o *common.ClientOptions) *Client {
	loadTestsClient := loadtests.NewLoadTestsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&loadTestsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LoadTestsClient: &loadTestsClient,
	}
}
