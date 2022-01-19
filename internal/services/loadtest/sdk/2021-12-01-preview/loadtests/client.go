package loadtests

import "github.com/Azure/go-autorest/autorest"

type LoadTestsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLoadTestsClientWithBaseURI(endpoint string) LoadTestsClient {
	return LoadTestsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
