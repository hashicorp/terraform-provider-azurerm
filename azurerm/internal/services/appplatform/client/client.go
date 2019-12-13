package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServicesClient    *appplatform.ServicesClient
	AppsClient        *appplatform.AppsClient
	DeploymentsClient *appplatform.DeploymentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ServicesClient := appplatform.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServicesClient.Client, o.ResourceManagerAuthorizer)

	AppsClient := appplatform.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppsClient.Client, o.ResourceManagerAuthorizer)

	DeploymentsClient := appplatform.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DeploymentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient:    &ServicesClient,
		AppsClient:        &AppsClient,
		DeploymentsClient: &DeploymentsClient,
	}
}
