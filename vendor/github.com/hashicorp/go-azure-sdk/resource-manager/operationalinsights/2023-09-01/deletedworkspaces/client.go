package deletedworkspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedWorkspacesClient struct {
	Client *resourcemanager.Client
}

func NewDeletedWorkspacesClientWithBaseURI(sdkApi sdkEnv.Api) (*DeletedWorkspacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deletedworkspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedWorkspacesClient: %+v", err)
	}

	return &DeletedWorkspacesClient{
		Client: client,
	}, nil
}
