package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AppServicePlansClient web.AppServicePlansClient
	AppServicesClient     web.AppsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AppServicePlansClient = web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AppServicePlansClient.Client, o.ResourceManagerAuthorizer)

	c.AppServicesClient = web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AppServicesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
