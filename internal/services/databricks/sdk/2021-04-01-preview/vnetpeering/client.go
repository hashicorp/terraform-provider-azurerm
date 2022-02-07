package vnetpeering

import "github.com/Azure/go-autorest/autorest"

type VNetPeeringClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVNetPeeringClientWithBaseURI(endpoint string) VNetPeeringClient {
	return VNetPeeringClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
