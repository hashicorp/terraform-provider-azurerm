package iotconnectors

import "github.com/Azure/go-autorest/autorest"

type IotConnectorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIotConnectorsClientWithBaseURI(endpoint string) IotConnectorsClient {
	return IotConnectorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
