package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServiceLinkerClient *servicelinker.ServicelinkerClient
	LinksClient         *links.LinksClient
}

func NewClient(o *common.ClientOptions) *Client {
	serviceLinkerClient := servicelinker.NewServicelinkerClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serviceLinkerClient.Client, o.ResourceManagerAuthorizer)

	linksClient := links.NewLinksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&linksClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServiceLinkerClient: &serviceLinkerClient,
		LinksClient:         &linksClient,
	}
}
