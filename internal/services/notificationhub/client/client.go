package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HubsClient       *notificationhubs.NotificationHubsClient
	NamespacesClient *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hubsClient := notificationhubs.NewNotificationHubsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&hubsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HubsClient:       &hubsClient,
		NamespacesClient: &namespacesClient,
	}
}
