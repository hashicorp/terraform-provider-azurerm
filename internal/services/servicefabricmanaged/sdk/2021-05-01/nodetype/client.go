package nodetype

import "github.com/Azure/go-autorest/autorest"

type NodeTypeClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNodeTypeClientWithBaseURI(endpoint string) NodeTypeClient {
	return NodeTypeClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
