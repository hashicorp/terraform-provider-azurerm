package azuresdkhacks

import "github.com/Azure/go-autorest/autorest"

// TODO 4.0: check if this is could be removed
// workaround for https://github.com/Azure/azure-rest-api-specs/issues/22947

type ReplicationProtectedItemsClient struct {
	Client  autorest.Client
	BaseUri string
}

func NewReplicationProtectedItemsClientWithBaseURI(endpoint string) ReplicationProtectedItemsClient {
	return ReplicationProtectedItemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		BaseUri: endpoint,
	}
}
