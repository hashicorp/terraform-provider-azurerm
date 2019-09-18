package notificationhub

import (
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HubsClient       *notificationhubs.Client
	NamespacesClient *notificationhubs.NamespacesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	NamespacesClient := notificationhubs.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	HubsClient := notificationhubs.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HubsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HubsClient:       &HubsClient,
		NamespacesClient: &NamespacesClient,
	}
}
