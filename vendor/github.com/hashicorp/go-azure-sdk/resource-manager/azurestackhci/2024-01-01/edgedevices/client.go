package edgedevices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeDevicesClient struct {
	Client *resourcemanager.Client
}

func NewEdgeDevicesClientWithBaseURI(sdkApi sdkEnv.Api) (*EdgeDevicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "edgedevices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EdgeDevicesClient: %+v", err)
	}

	return &EdgeDevicesClient{
		Client: client,
	}, nil
}
