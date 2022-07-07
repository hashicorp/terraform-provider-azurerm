package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	UserAssignedIdentitiesClient *managedidentity.ManagedIdentityClient
}

func NewClient(o *common.ClientOptions) *Client {
	userAssignedIdentitiesClient := managedidentity.NewManagedIdentityClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&userAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &userAssignedIdentitiesClient,
	}
}
