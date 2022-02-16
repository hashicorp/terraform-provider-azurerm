package fhirservices

import "github.com/Azure/go-autorest/autorest"

type FhirServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFhirServicesClientWithBaseURI(endpoint string) FhirServicesClient {
	return FhirServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
