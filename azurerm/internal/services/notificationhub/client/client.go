package client

import (
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HubsClient       *notificationhubs.Client
	NamespacesClient *notificationhubs.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hubsClient := notificationhubs.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&hubsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := notificationhubs.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HubsClient:       &hubsClient,
		NamespacesClient: &namespacesClient,
	}
}
