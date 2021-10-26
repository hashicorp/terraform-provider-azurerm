package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/sdk/2020-10-01/attestationproviders"
)

type Client struct {
	ProviderClient *attestationproviders.AttestationProvidersClient
}

func NewClient(o *common.ClientOptions) *Client {
	providerClient := attestationproviders.NewAttestationProvidersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&providerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ProviderClient: &providerClient,
	}
}
