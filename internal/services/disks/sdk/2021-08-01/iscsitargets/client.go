package iscsitargets

import "github.com/Azure/go-autorest/autorest"

type IscsiTargetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIscsiTargetsClientWithBaseURI(endpoint string) IscsiTargetsClient {
	return IscsiTargetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
