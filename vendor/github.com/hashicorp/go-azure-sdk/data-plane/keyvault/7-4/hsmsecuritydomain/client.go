package hsmsecuritydomain

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HSMSecurityDomainClient struct {
	Client *dataplane.Client
}

func NewHSMSecurityDomainClientUnconfigured() (*HSMSecurityDomainClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "hsmsecuritydomain", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HSMSecurityDomainClient: %+v", err)
	}

	return &HSMSecurityDomainClient{
		Client: client,
	}, nil
}

func (c *HSMSecurityDomainClient) HSMSecurityDomainClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewHSMSecurityDomainClientWithBaseURI(endpoint string) (*HSMSecurityDomainClient, error) {
	client, err := dataplane.NewClient(endpoint, "hsmsecuritydomain", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HSMSecurityDomainClient: %+v", err)
	}

	return &HSMSecurityDomainClient{
		Client: client,
	}, nil
}

func (c *HSMSecurityDomainClient) Clone(endpoint string) *HSMSecurityDomainClient {
	return &HSMSecurityDomainClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
