package proxy

import "github.com/Azure/go-autorest/autorest"

type ProxyClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProxyClientWithBaseURI(endpoint string) ProxyClient {
	return ProxyClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
