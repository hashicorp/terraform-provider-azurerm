package dataconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewDataConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*DataConnectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dataconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataConnectionsClient: %+v", err)
	}

	return &DataConnectionsClient{
		Client: client,
	}, nil
}
