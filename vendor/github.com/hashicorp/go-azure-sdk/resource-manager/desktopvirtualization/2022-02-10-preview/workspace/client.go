package workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceClient struct {
	Client *resourcemanager.Client
}

func NewWorkspaceClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkspaceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "workspace", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkspaceClient: %+v", err)
	}

	return &WorkspaceClient{
		Client: client,
	}, nil
}
