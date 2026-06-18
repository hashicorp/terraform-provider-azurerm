package linkedservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServicesClient struct {
	Client *dataplane.Client
}

func NewLinkedServicesClientUnconfigured() (*LinkedServicesClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "linkedservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LinkedServicesClient: %+v", err)
	}

	return &LinkedServicesClient{
		Client: client,
	}, nil
}

func (c *LinkedServicesClient) LinkedServicesClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewLinkedServicesClientWithBaseURI(endpoint string) (*LinkedServicesClient, error) {
	client, err := dataplane.NewClient(endpoint, "linkedservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LinkedServicesClient: %+v", err)
	}

	return &LinkedServicesClient{
		Client: client,
	}, nil
}

func (c *LinkedServicesClient) Clone(endpoint string) *LinkedServicesClient {
	return &LinkedServicesClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
