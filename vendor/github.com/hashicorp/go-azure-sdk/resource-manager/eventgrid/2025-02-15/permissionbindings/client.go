package permissionbindings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PermissionBindingsClient struct {
	Client *resourcemanager.Client
}

func NewPermissionBindingsClientWithBaseURI(sdkApi sdkEnv.Api) (*PermissionBindingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "permissionbindings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PermissionBindingsClient: %+v", err)
	}

	return &PermissionBindingsClient{
		Client: client,
	}, nil
}
