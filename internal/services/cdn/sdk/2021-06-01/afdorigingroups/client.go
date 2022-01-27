package afdorigingroups

import "github.com/Azure/go-autorest/autorest"

type AFDOriginGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAFDOriginGroupsClientWithBaseURI(endpoint string) AFDOriginGroupsClient {
	return AFDOriginGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
