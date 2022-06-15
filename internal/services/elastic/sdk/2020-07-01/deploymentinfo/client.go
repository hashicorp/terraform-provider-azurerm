package deploymentinfo

import "github.com/Azure/go-autorest/autorest"

type DeploymentInfoClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDeploymentInfoClientWithBaseURI(endpoint string) DeploymentInfoClient {
	return DeploymentInfoClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
