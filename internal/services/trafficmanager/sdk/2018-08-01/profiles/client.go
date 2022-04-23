package profiles

import "github.com/Azure/go-autorest/autorest"

type ProfilesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProfilesClientWithBaseURI(endpoint string) ProfilesClient {
	return ProfilesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
