package deviceupdates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeviceupdatesClient struct {
	Client *resourcemanager.Client
}

func NewDeviceupdatesClientWithBaseURI(api environments.Api) (*DeviceupdatesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "deviceupdates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeviceupdatesClient: %+v", err)
	}

	return &DeviceupdatesClient{
		Client: client,
	}, nil
}
