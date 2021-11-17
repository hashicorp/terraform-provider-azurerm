package systemtopics

import "github.com/Azure/go-autorest/autorest"

type SystemTopicsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSystemTopicsClientWithBaseURI(endpoint string) SystemTopicsClient {
	return SystemTopicsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
