package certificates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesClient struct {
	Client *dataplane.Client
}

func NewCertificatesClientUnconfigured() (*CertificatesClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "certificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificatesClient: %+v", err)
	}

	return &CertificatesClient{
		Client: client,
	}, nil
}

func (c *CertificatesClient) CertificatesClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewCertificatesClientWithBaseURI(endpoint string) (*CertificatesClient, error) {
	client, err := dataplane.NewClient(endpoint, "certificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificatesClient: %+v", err)
	}

	return &CertificatesClient{
		Client: client,
	}, nil
}

func (c *CertificatesClient) Clone(endpoint string) *CertificatesClient {
	return &CertificatesClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
