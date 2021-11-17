package topictypes

import "github.com/Azure/go-autorest/autorest"

type TopicTypesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTopicTypesClientWithBaseURI(endpoint string) TopicTypesClient {
	return TopicTypesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
