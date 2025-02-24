package saplandscapemonitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorClient struct {
	Client *resourcemanager.Client
}

func NewSapLandscapeMonitorClientWithBaseURI(sdkApi sdkEnv.Api) (*SapLandscapeMonitorClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "saplandscapemonitor", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SapLandscapeMonitorClient: %+v", err)
	}

	return &SapLandscapeMonitorClient{
		Client: client,
	}, nil
}
