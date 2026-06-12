package autonomousdatabasenationalcharactersets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseNationalCharacterSetsClient struct {
	Client *resourcemanager.Client
}

func NewAutonomousDatabaseNationalCharacterSetsClientWithBaseURI(sdkApi sdkEnv.Api) (*AutonomousDatabaseNationalCharacterSetsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "autonomousdatabasenationalcharactersets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutonomousDatabaseNationalCharacterSetsClient: %+v", err)
	}

	return &AutonomousDatabaseNationalCharacterSetsClient{
		Client: client,
	}, nil
}
