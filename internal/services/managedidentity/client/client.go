package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview/managedidentities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

// @stephybun: temporary and manual definition of ManagedIdentity Client to make upstream changes in Pandora's generator-terraform

type Client struct {
	ManagedIdentityClient *managedidentities.ManagedIdentitiesClient
}

func NewManagedIdentityClient(o *common.ClientOptions) (*Client, error) {
	client := managedidentities.NewManagedIdentitiesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedIdentityClient: &client,
	}, nil
}
