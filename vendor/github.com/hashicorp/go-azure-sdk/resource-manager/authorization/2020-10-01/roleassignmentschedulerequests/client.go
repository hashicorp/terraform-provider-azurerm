package roleassignmentschedulerequests

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleRequestsClient struct {
	Client *resourcemanager.Client
}

func NewRoleAssignmentScheduleRequestsClientWithBaseURI(sdkApi sdkEnv.Api) (*RoleAssignmentScheduleRequestsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "roleassignmentschedulerequests", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleAssignmentScheduleRequestsClient: %+v", err)
	}

	return &RoleAssignmentScheduleRequestsClient{
		Client: client,
	}, nil
}
