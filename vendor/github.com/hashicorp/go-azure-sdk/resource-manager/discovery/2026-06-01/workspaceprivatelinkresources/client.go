package workspaceprivatelinkresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePrivateLinkResourcesClient struct {
	Client *resourcemanager.Client
}

func NewWorkspacePrivateLinkResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkspacePrivateLinkResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "workspaceprivatelinkresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WorkspacePrivateLinkResourcesClient: %+v", err)
	}

	return &WorkspacePrivateLinkResourcesClient{
		Client: client,
	}, nil
}
