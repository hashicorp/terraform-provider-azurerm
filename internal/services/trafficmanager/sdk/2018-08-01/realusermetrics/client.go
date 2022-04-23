package realusermetrics

import "github.com/Azure/go-autorest/autorest"

type RealUserMetricsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRealUserMetricsClientWithBaseURI(endpoint string) RealUserMetricsClient {
	return RealUserMetricsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
