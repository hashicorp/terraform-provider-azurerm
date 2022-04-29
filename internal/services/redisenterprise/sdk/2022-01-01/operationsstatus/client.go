package operationsstatus

import "github.com/Azure/go-autorest/autorest"

type OperationsStatusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOperationsStatusClientWithBaseURI(endpoint string) OperationsStatusClient {
	return OperationsStatusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
