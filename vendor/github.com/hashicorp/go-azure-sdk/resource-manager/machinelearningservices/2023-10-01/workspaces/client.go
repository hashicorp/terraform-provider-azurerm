package workspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacesClient struct {
	Client *resourcemanager.Client
}

func NewWorkspacesClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkspacesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "workspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkspacesClient: %+v", err)
	}

	return &WorkspacesClient{
		Client: client,
	}, nil
}
