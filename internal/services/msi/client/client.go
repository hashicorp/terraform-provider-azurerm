package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	UserAssignedIdentitiesClient *managedidentities.ManagedIdentitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	userAssignedIdentitiesClient := managedidentities.NewManagedIdentitiesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&userAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &userAssignedIdentitiesClient,
	}
}
