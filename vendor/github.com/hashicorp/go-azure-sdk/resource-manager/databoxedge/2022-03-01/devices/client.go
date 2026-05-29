package devices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevicesClient struct {
	Client *resourcemanager.Client
}

func NewDevicesClientWithBaseURI(sdkApi sdkEnv.Api) (*DevicesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "devices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DevicesClient: %+v", err)
	}

	return &DevicesClient{
		Client: client,
	}, nil
}
