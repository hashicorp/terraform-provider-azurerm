package edgemodules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesClient struct {
	Client *resourcemanager.Client
}

func NewEdgeModulesClientWithBaseURI(sdkApi sdkEnv.Api) (*EdgeModulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "edgemodules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EdgeModulesClient: %+v", err)
	}

	return &EdgeModulesClient{
		Client: client,
	}, nil
}
