package portal

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DashboardsClient *portal.DashboardsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	DashboardsClient := portal.NewDashboardsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DashboardsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DashboardsClient: &DashboardsClient,
	}
}
