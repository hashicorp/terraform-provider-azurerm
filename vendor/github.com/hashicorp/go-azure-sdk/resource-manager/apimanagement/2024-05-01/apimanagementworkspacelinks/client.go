package apimanagementworkspacelinks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementWorkspaceLinksClient struct {
	Client *resourcemanager.Client
}

func NewApiManagementWorkspaceLinksClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiManagementWorkspaceLinksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apimanagementworkspacelinks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiManagementWorkspaceLinksClient: %+v", err)
	}

	return &ApiManagementWorkspaceLinksClient{
		Client: client,
	}, nil
}
