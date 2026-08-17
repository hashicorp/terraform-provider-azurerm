package systemassignedidentities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemAssignedIdentitiesClient struct {
	Client *resourcemanager.Client
}

func NewSystemAssignedIdentitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*SystemAssignedIdentitiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "systemassignedidentities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SystemAssignedIdentitiesClient: %+v", err)
	}

	return &SystemAssignedIdentitiesClient{
		Client: client,
	}, nil
}
