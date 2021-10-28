package keys

import "github.com/Azure/go-autorest/autorest"

type KeysClient struct {
	Client  autorest.Client
	baseUri string
}

func NewKeysClientWithBaseURI(endpoint string) KeysClient {
	return KeysClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
