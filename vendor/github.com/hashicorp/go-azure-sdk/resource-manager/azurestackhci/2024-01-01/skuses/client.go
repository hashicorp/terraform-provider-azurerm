package skuses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusesClient struct {
	Client *resourcemanager.Client
}

func NewSkusesClientWithBaseURI(sdkApi sdkEnv.Api) (*SkusesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "skuses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SkusesClient: %+v", err)
	}

	return &SkusesClient{
		Client: client,
	}, nil
}
