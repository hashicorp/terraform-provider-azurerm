package jobschedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScheduleClient struct {
	Client *resourcemanager.Client
}

func NewJobScheduleClientWithBaseURI(api environments.Api) (*JobScheduleClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "jobschedule", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating JobScheduleClient: %+v", err)
	}

	return &JobScheduleClient{
		Client: client,
	}, nil
}
