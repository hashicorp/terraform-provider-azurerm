package devcenters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCentersClient struct {
	Client *resourcemanager.Client
}

func NewDevCentersClientWithBaseURI(sdkApi sdkEnv.Api) (*DevCentersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "devcenters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DevCentersClient: %+v", err)
	}

	return &DevCentersClient{
		Client: client,
	}, nil
}
