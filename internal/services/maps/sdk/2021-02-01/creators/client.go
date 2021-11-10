package creators

import "github.com/Azure/go-autorest/autorest"

type CreatorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCreatorsClientWithBaseURI(endpoint string) CreatorsClient {
	return CreatorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
