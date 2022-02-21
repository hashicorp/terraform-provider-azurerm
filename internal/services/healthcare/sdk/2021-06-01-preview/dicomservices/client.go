package dicomservices

import "github.com/Azure/go-autorest/autorest"

type DicomServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDicomServicesClientWithBaseURI(endpoint string) DicomServicesClient {
	return DicomServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
