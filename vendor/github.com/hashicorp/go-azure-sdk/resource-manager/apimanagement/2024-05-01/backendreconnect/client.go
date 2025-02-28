package backendreconnect

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendReconnectClient struct {
	Client *resourcemanager.Client
}

func NewBackendReconnectClientWithBaseURI(sdkApi sdkEnv.Api) (*BackendReconnectClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backendreconnect", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackendReconnectClient: %+v", err)
	}

	return &BackendReconnectClient{
		Client: client,
	}, nil
}
