package edgeactions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionsClient struct {
	Client *resourcemanager.Client
}

func NewEdgeActionsClientWithBaseURI(sdkApi sdkEnv.Api) (*EdgeActionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "edgeactions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EdgeActionsClient: %+v", err)
	}

	return &EdgeActionsClient{
		Client: client,
	}, nil
}
