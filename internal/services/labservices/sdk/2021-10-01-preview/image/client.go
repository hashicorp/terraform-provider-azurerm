package image

import "github.com/Azure/go-autorest/autorest"

type ImageClient struct {
	Client  autorest.Client
	baseUri string
}

func NewImageClientWithBaseURI(endpoint string) ImageClient {
	return ImageClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
