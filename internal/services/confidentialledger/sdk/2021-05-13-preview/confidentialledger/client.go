package confidentialledger

import "github.com/Azure/go-autorest/autorest"

type ConfidentialLedgerClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfidentialLedgerClientWithBaseURI(endpoint string) ConfidentialLedgerClient {
	return ConfidentialLedgerClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
