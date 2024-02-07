package locationbasedperformancetier

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationBasedPerformanceTierClient struct {
	Client *resourcemanager.Client
}

func NewLocationBasedPerformanceTierClientWithBaseURI(sdkApi sdkEnv.Api) (*LocationBasedPerformanceTierClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "locationbasedperformancetier", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocationBasedPerformanceTierClient: %+v", err)
	}

	return &LocationBasedPerformanceTierClient{
		Client: client,
	}, nil
}
