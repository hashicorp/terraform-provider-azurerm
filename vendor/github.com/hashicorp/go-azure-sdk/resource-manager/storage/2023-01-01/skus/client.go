package skus

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusClient struct {
	Client *resourcemanager.Client
}

func NewSkusClientWithBaseURI(sdkApi sdkEnv.Api) (*SkusClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "skus", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SkusClient: %+v", err)
	}

	return &SkusClient{
		Client: client,
	}, nil
}
