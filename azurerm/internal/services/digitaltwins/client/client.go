package client

import (
	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	InstanceClient *digitaltwins.Client
}

func NewClient(o *common.ClientOptions) *Client {
	InstanceClient := digitaltwins.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&InstanceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		InstanceClient: &InstanceClient,
	}
}
