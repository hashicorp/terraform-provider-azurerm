package roleassignmentscheduleinstances

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleInstancesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRoleAssignmentScheduleInstancesClientWithBaseURI(endpoint string) RoleAssignmentScheduleInstancesClient {
	return RoleAssignmentScheduleInstancesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
