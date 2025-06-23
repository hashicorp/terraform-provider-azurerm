package standbycontainergrouppools

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyContainerGroupPoolsClient struct {
	Client *resourcemanager.Client
}

func NewStandbyContainerGroupPoolsClientWithBaseURI(sdkApi sdkEnv.Api) (*StandbyContainerGroupPoolsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "standbycontainergrouppools", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StandbyContainerGroupPoolsClient: %+v", err)
	}

	return &StandbyContainerGroupPoolsClient{
		Client: client,
	}, nil
}
