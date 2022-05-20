package client

import (
	"github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2021-11-01-preview/iotcentral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppsClient *iotcentral.AppsClient
}

func NewClient(o *common.ClientOptions) *Client {
	AppsClient := iotcentral.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AppsClient: &AppsClient,
	}
}
