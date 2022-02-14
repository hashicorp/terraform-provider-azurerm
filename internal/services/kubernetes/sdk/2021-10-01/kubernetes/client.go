package kubernetes

import "github.com/Azure/go-autorest/autorest"

type KubernetesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewKubernetesClientWithBaseURI(endpoint string) KubernetesClient {
	return KubernetesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
