package vaults

import "github.com/Azure/go-autorest/autorest"

type VaultsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVaultsClientWithBaseURI(endpoint string) VaultsClient {
	return VaultsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
