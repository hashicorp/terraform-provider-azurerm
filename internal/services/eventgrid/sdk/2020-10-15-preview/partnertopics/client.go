package partnertopics

import "github.com/Azure/go-autorest/autorest"

type PartnerTopicsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPartnerTopicsClientWithBaseURI(endpoint string) PartnerTopicsClient {
	return PartnerTopicsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
