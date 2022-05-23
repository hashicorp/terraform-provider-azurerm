package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/sdk/2022-05-01/links"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/sdk/2022-05-01/servicelinker"
)

type Client struct {
	ServiceLinkerClient *servicelinker.ServiceLinkerClient
	LinksClient         *links.LinksClient
}

func NewClient(o *common.ClientOptions) *Client {
	serviceLinkerClient := servicelinker.NewServiceLinkerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serviceLinkerClient.Client, o.ResourceManagerAuthorizer)

	linksClient := links.NewLinksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&linksClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServiceLinkerClient: &serviceLinkerClient,
		LinksClient:         &linksClient,
	}
}
