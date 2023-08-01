package managedapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedAPIsClient struct {
	Client *resourcemanager.Client
}

func NewManagedAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedAPIsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "managedapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedAPIsClient: %+v", err)
	}

	return &ManagedAPIsClient{
		Client: client,
	}, nil
}
