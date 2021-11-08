package datalakestoreaccounts

import "github.com/Azure/go-autorest/autorest"

type DataLakeStoreAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDataLakeStoreAccountsClientWithBaseURI(endpoint string) DataLakeStoreAccountsClient {
	return DataLakeStoreAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
