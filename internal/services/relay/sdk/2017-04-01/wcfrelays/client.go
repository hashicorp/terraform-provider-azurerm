package wcfrelays

import "github.com/Azure/go-autorest/autorest"

type WCFRelaysClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWCFRelaysClientWithBaseURI(endpoint string) WCFRelaysClient {
	return WCFRelaysClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
