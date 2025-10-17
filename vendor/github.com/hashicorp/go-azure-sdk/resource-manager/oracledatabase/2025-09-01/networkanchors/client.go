package networkanchors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAnchorsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkAnchorsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkAnchorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkanchors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkAnchorsClient: %+v", err)
	}

	return &NetworkAnchorsClient{
		Client: client,
	}, nil
}
