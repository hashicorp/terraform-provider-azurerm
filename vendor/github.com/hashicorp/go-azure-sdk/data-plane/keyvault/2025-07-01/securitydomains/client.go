package securitydomains

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecuritydomainsClient struct {
	Client *dataplane.Client
}

func NewSecuritydomainsClientUnconfigured() (*SecuritydomainsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "securitydomains", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecuritydomainsClient: %+v", err)
	}

	return &SecuritydomainsClient{
		Client: client,
	}, nil
}

func (c *SecuritydomainsClient) SecuritydomainsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSecuritydomainsClientWithBaseURI(endpoint string) (*SecuritydomainsClient, error) {
	client, err := dataplane.NewClient(endpoint, "securitydomains", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecuritydomainsClient: %+v", err)
	}

	return &SecuritydomainsClient{
		Client: client,
	}, nil
}
