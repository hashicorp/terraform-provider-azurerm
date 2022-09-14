package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
