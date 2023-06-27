package guestconfigurationassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestConfigurationAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewGuestConfigurationAssignmentsClientWithBaseURI(api environments.Api) (*GuestConfigurationAssignmentsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "guestconfigurationassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GuestConfigurationAssignmentsClient: %+v", err)
	}

	return &GuestConfigurationAssignmentsClient{
		Client: client,
	}, nil
}
