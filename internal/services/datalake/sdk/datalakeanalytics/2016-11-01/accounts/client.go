package accounts

import "github.com/Azure/go-autorest/autorest"

type AccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAccountsClientWithBaseURI(endpoint string) AccountsClient {
	return AccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
