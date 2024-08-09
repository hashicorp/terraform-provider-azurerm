package backend

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendClient struct {
	Client *resourcemanager.Client
}

func NewBackendClientWithBaseURI(sdkApi sdkEnv.Api) (*BackendClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "backend", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackendClient: %+v", err)
	}

	return &BackendClient{
		Client: client,
	}, nil
}
