package hybridkubernetes

import "github.com/Azure/go-autorest/autorest"

type HybridKubernetesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHybridKubernetesClientWithBaseURI(endpoint string) HybridKubernetesClient {
	return HybridKubernetesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
