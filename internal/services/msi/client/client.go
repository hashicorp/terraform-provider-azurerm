package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/sdk/2018-11-30/managedidentity"
)

type Client struct {
	UserAssignedIdentitiesClient *managedidentity.ManagedIdentityClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserAssignedIdentitiesClient := managedidentity.NewManagedIdentityClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&UserAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &UserAssignedIdentitiesClient,
	}
}
