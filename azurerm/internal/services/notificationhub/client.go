package notificationhub

import (
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HubsClient       notificationhubs.Client
	NamespacesClient notificationhubs.NamespacesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.NamespacesClient = notificationhubs.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.NamespacesClient.Client, o.ResourceManagerAuthorizer)

	c.HubsClient = notificationhubs.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.HubsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
