package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/attestation/mgmt/2018-09-01-preview/attestation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProviderClient *attestation.ProvidersClient
}

func NewClient(o *common.ClientOptions) *Client {
	providerClient := attestation.NewProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&providerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ProviderClient: &providerClient,
	}
}
