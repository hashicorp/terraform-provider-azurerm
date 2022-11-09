package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	UserClient *user.UserClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserClient := user.NewUserClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&UserClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserClient: &UserClient,
	}
}
