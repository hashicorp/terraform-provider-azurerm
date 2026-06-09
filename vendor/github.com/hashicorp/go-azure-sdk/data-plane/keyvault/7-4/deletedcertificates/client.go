package deletedcertificates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedCertificatesClient struct {
	Client *dataplane.Client
}

func NewDeletedCertificatesClientUnconfigured() (*DeletedCertificatesClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "deletedcertificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedCertificatesClient: %+v", err)
	}

	return &DeletedCertificatesClient{
		Client: client,
	}, nil
}

func (c *DeletedCertificatesClient) DeletedCertificatesClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewDeletedCertificatesClientWithBaseURI(endpoint string) (*DeletedCertificatesClient, error) {
	client, err := dataplane.NewClient(endpoint, "deletedcertificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedCertificatesClient: %+v", err)
	}

	return &DeletedCertificatesClient{
		Client: client,
	}, nil
}

func (c *DeletedCertificatesClient) Clone(endpoint string) *DeletedCertificatesClient {
	return &DeletedCertificatesClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
