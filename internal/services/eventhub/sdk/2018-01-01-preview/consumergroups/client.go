package consumergroups

import "github.com/Azure/go-autorest/autorest"

type ConsumerGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConsumerGroupsClientWithBaseURI(endpoint string) ConsumerGroupsClient {
	return ConsumerGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
