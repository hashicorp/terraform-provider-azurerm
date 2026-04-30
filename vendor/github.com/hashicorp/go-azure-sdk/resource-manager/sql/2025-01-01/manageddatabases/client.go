package manageddatabases

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedDatabasesClient struct {
	Client *resourcemanager.Client
}

func NewManagedDatabasesClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedDatabasesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "manageddatabases", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedDatabasesClient: %+v", err)
	}

	return &ManagedDatabasesClient{
		Client: client,
	}, nil
}
