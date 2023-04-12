package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ProviderClient *attestationproviders.AttestationProvidersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	providerClient, err := attestationproviders.NewAttestationProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(providerClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ProviderClient: providerClient,
	}, nil
}
