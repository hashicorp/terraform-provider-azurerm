package startstopmanagedinstanceschedules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartStopManagedInstanceSchedulesClient struct {
	Client *resourcemanager.Client
}

func NewStartStopManagedInstanceSchedulesClientWithBaseURI(sdkApi sdkEnv.Api) (*StartStopManagedInstanceSchedulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "startstopmanagedinstanceschedules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StartStopManagedInstanceSchedulesClient: %+v", err)
	}

	return &StartStopManagedInstanceSchedulesClient{
		Client: client,
	}, nil
}
