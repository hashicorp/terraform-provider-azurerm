package graph

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationsClient      graphrbac.ApplicationsClient
	ServicePrincipalsClient graphrbac.ServicePrincipalsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ApplicationsClient = graphrbac.NewApplicationsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&c.ApplicationsClient.Client, o.GraphAuthorizer)

	c.ServicePrincipalsClient = graphrbac.NewServicePrincipalsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&c.ServicePrincipalsClient.Client, o.GraphAuthorizer)

	return &c
}
