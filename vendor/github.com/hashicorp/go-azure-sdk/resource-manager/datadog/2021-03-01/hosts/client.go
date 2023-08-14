package hosts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostsClient struct {
	Client *resourcemanager.Client
}

func NewHostsClientWithBaseURI(sdkApi sdkEnv.Api) (*HostsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "hosts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HostsClient: %+v", err)
	}

	return &HostsClient{
		Client: client,
	}, nil
}
