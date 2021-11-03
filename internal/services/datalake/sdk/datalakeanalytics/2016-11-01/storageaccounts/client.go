package storageaccounts

import "github.com/Azure/go-autorest/autorest"

type StorageAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStorageAccountsClientWithBaseURI(endpoint string) StorageAccountsClient {
	return StorageAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
