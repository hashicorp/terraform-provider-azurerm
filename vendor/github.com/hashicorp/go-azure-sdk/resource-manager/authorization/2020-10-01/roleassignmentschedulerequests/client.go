package roleassignmentschedulerequests

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleRequestsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRoleAssignmentScheduleRequestsClientWithBaseURI(endpoint string) RoleAssignmentScheduleRequestsClient {
	return RoleAssignmentScheduleRequestsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
