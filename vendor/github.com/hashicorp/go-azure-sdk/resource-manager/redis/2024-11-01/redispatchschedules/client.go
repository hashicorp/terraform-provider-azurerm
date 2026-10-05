package redispatchschedules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisPatchSchedulesClient struct {
	Client *resourcemanager.Client
}

func NewRedisPatchSchedulesClientWithBaseURI(sdkApi sdkEnv.Api) (*RedisPatchSchedulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "redispatchschedules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RedisPatchSchedulesClient: %+v", err)
	}

	return &RedisPatchSchedulesClient{
		Client: client,
	}, nil
}
