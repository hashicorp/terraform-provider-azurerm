package capacityreservations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationsClient struct {
	Client *resourcemanager.Client
}

func NewCapacityReservationsClientWithBaseURI(sdkApi sdkEnv.Api) (*CapacityReservationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "capacityreservations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CapacityReservationsClient: %+v", err)
	}

	return &CapacityReservationsClient{
		Client: client,
	}, nil
}
