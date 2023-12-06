package schedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleClient struct {
	Client *resourcemanager.Client
}

func NewScheduleClientWithBaseURI(sdkApi sdkEnv.Api) (*ScheduleClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "schedule", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScheduleClient: %+v", err)
	}

	return &ScheduleClient{
		Client: client,
	}, nil
}
