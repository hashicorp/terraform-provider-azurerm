package managedhsmkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmKeysClient struct {
	Client *resourcemanager.Client
}

func NewManagedHsmKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedHsmKeysClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "managedhsmkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedHsmKeysClient: %+v", err)
	}

	return &ManagedHsmKeysClient{
		Client: client,
	}, nil
}
