package managedidentities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIdentitiesClient struct {
	Client *resourcemanager.Client
}

func NewManagedIdentitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedIdentitiesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "managedidentities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedIdentitiesClient: %+v", err)
	}

	return &ManagedIdentitiesClient{
		Client: client,
	}, nil
}
