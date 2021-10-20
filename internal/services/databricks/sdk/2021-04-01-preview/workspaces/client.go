package workspaces

import "github.com/Azure/go-autorest/autorest"

type WorkspacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWorkspacesClientWithBaseURI(endpoint string) WorkspacesClient {
	return WorkspacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
