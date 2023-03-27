package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SAPVirtualInstancesClient *sapvirtualinstances.SAPVirtualInstancesClient
}

func NewClient(o *common.ClientOptions) *Client {
	sapVirtualInstancesClient := sapvirtualinstances.NewSAPVirtualInstancesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(sapVirtualInstancesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SAPVirtualInstancesClient: &sapVirtualInstancesClient,
	}
}
