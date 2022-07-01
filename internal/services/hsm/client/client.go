package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2018-10-31-preview/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DedicatedHsmClient *dedicatedhsms.DedicatedHsmsClient
}

func NewClient(o *common.ClientOptions) *Client {
	dedicatedHsmClient := dedicatedhsms.NewDedicatedHsmsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dedicatedHsmClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DedicatedHsmClient: &dedicatedHsmClient,
	}
}
