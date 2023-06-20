package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-01-31/communicationsgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-01-31/testlines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CommunicationsGatewaysClient *communicationsgateways.CommunicationsGatewaysClient
	TestLinesClient              *testlines.TestLinesClient
}

func NewClient(o *common.ClientOptions) *Client {

	communicationsGatewaysClient := communicationsgateways.NewCommunicationsGatewaysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&communicationsGatewaysClient.Client, o.ResourceManagerAuthorizer)

	testLinesClient := testlines.NewTestLinesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&testLinesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CommunicationsGatewaysClient: &communicationsGatewaysClient,
		TestLinesClient:              &testLinesClient,
	}
}
