package communicationservice

import "github.com/Azure/go-autorest/autorest"

type CommunicationServiceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCommunicationServiceClientWithBaseURI(endpoint string) CommunicationServiceClient {
	return CommunicationServiceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
