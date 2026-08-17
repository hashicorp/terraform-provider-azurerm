package autonomousdatabasebackups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseBackupsClient struct {
	Client *resourcemanager.Client
}

func NewAutonomousDatabaseBackupsClientWithBaseURI(sdkApi sdkEnv.Api) (*AutonomousDatabaseBackupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autonomousdatabasebackups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutonomousDatabaseBackupsClient: %+v", err)
	}

	return &AutonomousDatabaseBackupsClient{
		Client: client,
	}, nil
}
