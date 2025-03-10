package networkstatus

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkStatusClient struct {
	Client *resourcemanager.Client
}

func NewNetworkStatusClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkStatusClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkstatus", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkStatusClient: %+v", err)
	}

	return &NetworkStatusClient{
		Client: client,
	}, nil
}
