package graphquery

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphQueryClient struct {
	Client *resourcemanager.Client
}

func NewGraphQueryClientWithBaseURI(sdkApi sdkEnv.Api) (*GraphQueryClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "graphquery", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GraphQueryClient: %+v", err)
	}

	return &GraphQueryClient{
		Client: client,
	}, nil
}
