package jobs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsClient struct {
	Client *dataplane.Client
}

func NewJobsClientUnconfigured() (*JobsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "", "jobs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobsClient: %+v", err)
	}

	return &JobsClient{
		Client: client,
	}, nil
}

func (c *JobsClient) JobsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func (c *JobsClient) JobsClientSetAdditionalEndpoint(endpoint string) {
	c.Client.AdditionalEndpoint = endpoint
}

func NewJobsClientWithBaseURI(endpoint string, additionalEndpoint string) (*JobsClient, error) {
	client, err := dataplane.NewClient(endpoint, additionalEndpoint, "jobs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobsClient: %+v", err)
	}

	return &JobsClient{
		Client: client,
	}, nil
}
