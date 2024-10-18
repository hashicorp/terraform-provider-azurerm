package roleeligibilityscheduleinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleEligibilityScheduleInstancesClient struct {
	Client *resourcemanager.Client
}

func NewRoleEligibilityScheduleInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*RoleEligibilityScheduleInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "roleeligibilityscheduleinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleEligibilityScheduleInstancesClient: %+v", err)
	}

	return &RoleEligibilityScheduleInstancesClient{
		Client: client,
	}, nil
}
