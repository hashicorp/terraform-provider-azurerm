package capacitypools

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityPoolsClient struct {
	Client *resourcemanager.Client
}

func NewCapacityPoolsClientWithBaseURI(sdkApi sdkEnv.Api) (*CapacityPoolsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "capacitypools", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CapacityPoolsClient: %+v", err)
	}

	return &CapacityPoolsClient{
		Client: client,
	}, nil
}
