package dedicatedhosts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostsClient struct {
	Client *resourcemanager.Client
}

func NewDedicatedHostsClientWithBaseURI(sdkApi sdkEnv.Api) (*DedicatedHostsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dedicatedhosts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DedicatedHostsClient: %+v", err)
	}

	return &DedicatedHostsClient{
		Client: client,
	}, nil
}
