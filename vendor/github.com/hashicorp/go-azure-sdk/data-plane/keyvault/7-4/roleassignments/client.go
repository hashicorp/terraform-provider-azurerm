package roleassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentsClient struct {
	Client *dataplane.Client
}

func NewRoleAssignmentsClientUnconfigured() (*RoleAssignmentsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "roleassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleAssignmentsClient: %+v", err)
	}

	return &RoleAssignmentsClient{
		Client: client,
	}, nil
}

func (c *RoleAssignmentsClient) RoleAssignmentsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewRoleAssignmentsClientWithBaseURI(endpoint string) (*RoleAssignmentsClient, error) {
	client, err := dataplane.NewClient(endpoint, "roleassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleAssignmentsClient: %+v", err)
	}

	return &RoleAssignmentsClient{
		Client: client,
	}, nil
}

func (c *RoleAssignmentsClient) Clone(endpoint string) *RoleAssignmentsClient {
	return &RoleAssignmentsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
