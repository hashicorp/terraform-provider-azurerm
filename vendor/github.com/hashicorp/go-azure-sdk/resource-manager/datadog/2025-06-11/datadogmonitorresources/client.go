package datadogmonitorresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogMonitorResourcesClient struct {
	Client *resourcemanager.Client
}

func NewDatadogMonitorResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*DatadogMonitorResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "datadogmonitorresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatadogMonitorResourcesClient: %+v", err)
	}

	return &DatadogMonitorResourcesClient{
		Client: client,
	}, nil
}
