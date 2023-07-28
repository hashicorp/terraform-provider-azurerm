package workspaceprivateendpointconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePrivateEndpointConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewWorkspacePrivateEndpointConnectionsClientWithBaseURI(api environments.Api) (*WorkspacePrivateEndpointConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "workspaceprivateendpointconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkspacePrivateEndpointConnectionsClient: %+v", err)
	}

	return &WorkspacePrivateEndpointConnectionsClient{
		Client: client,
	}, nil
}
