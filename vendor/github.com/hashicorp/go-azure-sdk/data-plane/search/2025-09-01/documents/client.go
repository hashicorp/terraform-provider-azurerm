package documents

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DocumentsClient struct {
	Client *dataplane.Client
}

func NewDocumentsClientUnconfigured() (*DocumentsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "documents", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DocumentsClient: %+v", err)
	}

	return &DocumentsClient{
		Client: client,
	}, nil
}

func (c *DocumentsClient) DocumentsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewDocumentsClientWithBaseURI(endpoint string) (*DocumentsClient, error) {
	client, err := dataplane.NewClient(endpoint, "documents", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DocumentsClient: %+v", err)
	}

	return &DocumentsClient{
		Client: client,
	}, nil
}

func (c *DocumentsClient) Clone(endpoint string) *DocumentsClient {
	return &DocumentsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
