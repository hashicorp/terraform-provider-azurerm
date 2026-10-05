package autonomousdatabaseversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseVersionsClient struct {
	Client *resourcemanager.Client
}

func NewAutonomousDatabaseVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*AutonomousDatabaseVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autonomousdatabaseversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutonomousDatabaseVersionsClient: %+v", err)
	}

	return &AutonomousDatabaseVersionsClient{
		Client: client,
	}, nil
}
