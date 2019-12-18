package client

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationsClient      *graphrbac.ApplicationsClient
	ServicePrincipalsClient *graphrbac.ServicePrincipalsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ApplicationsClient := graphrbac.NewApplicationsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&ApplicationsClient.Client, o.GraphAuthorizer)

	ServicePrincipalsClient := graphrbac.NewServicePrincipalsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&ServicePrincipalsClient.Client, o.GraphAuthorizer)

	return &Client{
		ApplicationsClient:      &ApplicationsClient,
		ServicePrincipalsClient: &ServicePrincipalsClient,
	}
}
