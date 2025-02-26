package aad

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AADClient struct {
	Client *resourcemanager.Client
}

func NewAADClientWithBaseURI(sdkApi sdkEnv.Api) (*AADClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "aad", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AADClient: %+v", err)
	}

	return &AADClient{
		Client: client,
	}, nil
}
