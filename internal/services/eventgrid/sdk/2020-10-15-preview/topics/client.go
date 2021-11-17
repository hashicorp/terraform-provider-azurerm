package topics

import "github.com/Azure/go-autorest/autorest"

type TopicsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTopicsClientWithBaseURI(endpoint string) TopicsClient {
	return TopicsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
