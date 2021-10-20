package videoanalyzer

import "github.com/Azure/go-autorest/autorest"

type VideoAnalyzerClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVideoAnalyzerClientWithBaseURI(endpoint string) VideoAnalyzerClient {
	return VideoAnalyzerClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
