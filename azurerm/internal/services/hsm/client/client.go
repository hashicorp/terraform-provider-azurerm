package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/hardwaresecuritymodules/mgmt/2018-10-31-preview/hardwaresecuritymodules"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DedicatedHsmClient *hardwaresecuritymodules.DedicatedHsmClient
}

func NewClient(o *common.ClientOptions) *Client {
	dedicatedHsmClient := hardwaresecuritymodules.NewDedicatedHsmClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dedicatedHsmClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DedicatedHsmClient: &dedicatedHsmClient,
	}
}
