package newrelicmonitorresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NewRelicMonitorResourcesClient struct {
	Client *resourcemanager.Client
}

func NewNewRelicMonitorResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*NewRelicMonitorResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "newrelicmonitorresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NewRelicMonitorResourcesClient: %+v", err)
	}

	return &NewRelicMonitorResourcesClient{
		Client: client,
	}, nil
}
