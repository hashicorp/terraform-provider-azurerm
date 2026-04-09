package schedules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulesClient struct {
	Client *resourcemanager.Client
}

func NewSchedulesClientWithBaseURI(sdkApi sdkEnv.Api) (*SchedulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "schedules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SchedulesClient: %+v", err)
	}

	return &SchedulesClient{
		Client: client,
	}, nil
}
