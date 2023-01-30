package roleeligibilityscheduleinstances

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleEligibilityScheduleInstancesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRoleEligibilityScheduleInstancesClientWithBaseURI(endpoint string) RoleEligibilityScheduleInstancesClient {
	return RoleEligibilityScheduleInstancesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
