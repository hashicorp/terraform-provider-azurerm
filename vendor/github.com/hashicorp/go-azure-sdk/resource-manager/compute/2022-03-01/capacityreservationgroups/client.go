package capacityreservationgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationGroupsClient struct {
	Client *resourcemanager.Client
}

func NewCapacityReservationGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*CapacityReservationGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "capacityreservationgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CapacityReservationGroupsClient: %+v", err)
	}

	return &CapacityReservationGroupsClient{
		Client: client,
	}, nil
}
