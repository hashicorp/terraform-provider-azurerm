package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2020-08-20/communicationservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServiceClient *communicationservice.CommunicationServiceClient
}

func NewClient(o *common.ClientOptions) *Client {
	serviceClient := communicationservice.NewCommunicationServiceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serviceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServiceClient: &serviceClient,
	}
}
