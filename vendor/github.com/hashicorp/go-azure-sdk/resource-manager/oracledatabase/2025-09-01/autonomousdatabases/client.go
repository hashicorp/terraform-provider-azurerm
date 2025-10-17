package autonomousdatabases

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabasesClient struct {
	Client *resourcemanager.Client
}

func NewAutonomousDatabasesClientWithBaseURI(sdkApi sdkEnv.Api) (*AutonomousDatabasesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autonomousdatabases", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutonomousDatabasesClient: %+v", err)
	}

	return &AutonomousDatabasesClient{
		Client: client,
	}, nil
}
