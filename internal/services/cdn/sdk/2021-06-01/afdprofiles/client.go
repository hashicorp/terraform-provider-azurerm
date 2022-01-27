package afdprofiles

import "github.com/Azure/go-autorest/autorest"

type AFDProfilesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAFDProfilesClientWithBaseURI(endpoint string) AFDProfilesClient {
	return AFDProfilesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
