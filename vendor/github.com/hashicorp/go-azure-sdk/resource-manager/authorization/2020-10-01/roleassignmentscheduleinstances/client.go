package roleassignmentscheduleinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleInstancesClient struct {
	Client *resourcemanager.Client
}

func NewRoleAssignmentScheduleInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*RoleAssignmentScheduleInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "roleassignmentscheduleinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleAssignmentScheduleInstancesClient: %+v", err)
	}

	return &RoleAssignmentScheduleInstancesClient{
		Client: client,
	}, nil
}
