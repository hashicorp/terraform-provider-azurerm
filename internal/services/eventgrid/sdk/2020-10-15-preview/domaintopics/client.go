package domaintopics

import "github.com/Azure/go-autorest/autorest"

type DomainTopicsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDomainTopicsClientWithBaseURI(endpoint string) DomainTopicsClient {
	return DomainTopicsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
