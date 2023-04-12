package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2020-08-20/communicationservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServiceClient *communicationservice.CommunicationServiceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	serviceClient, err := communicationservice.NewCommunicationServiceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Service client: %+v", err)
	}
	o.Configure(serviceClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ServiceClient: serviceClient,
	}, nil
}
