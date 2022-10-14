package v2018_11_30

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentities"
)

type Client struct {
	ManagedIdentities *managedidentities.ManagedIdentitiesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	managedIdentitiesClient := managedidentities.NewManagedIdentitiesClientWithBaseURI(endpoint)
	configureAuthFunc(&managedIdentitiesClient.Client)

	return Client{
		ManagedIdentities: &managedIdentitiesClient,
	}
}
